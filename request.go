package stdutil

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
	"github.com/gorilla/mux"
	"github.com/narsilworks/livenote"
)

const REQUEST_VERSION string = "1.0.0.0"
const REQUEST_MODIFIED string = "15112023"

var (
	reqTimeOut int
	ct         *http.Transport
)

// CustomPayload - payload for JWT
type CustomPayload struct {
	jwt.Payload
	UserName      string `json:"usr,omitempty"` // Username payload for JWT
	Domain        string `json:"dom,omitempty"` // Domain payload for JWT
	ApplicationID string `json:"app,omitempty"` // Application payload for JWT
	DeviceID      string `json:"dev,omitempty"` // Device id payload for JWT
	TenantID      string `json:"tnt,omitempty"` // Tenant id payload for JWT
}

// ResultData - a result structure and a JSON raw message
type ResultData struct {
	Result
	Data json.RawMessage `json:"data"`
}

// ResultAny - a result structure with an empty interface
type ResultAny[T any] struct {
	Result
	Data T `json:"data"`
}

func init() {
	reqTimeOut = 30
	ct = http.DefaultTransport.(*http.Transport).Clone()
	ct.MaxIdleConns = 100
	ct.MaxConnsPerHost = 100
	ct.MaxIdleConnsPerHost = 100
}

// SetRequestTimeOut sets the new timeout value
func SetRequestTimeout(timeout int) {
	reqTimeOut = timeout
}

// ExecuteJSONAPI wraps http operation that change or read data and returns a custom result
func ExecuteJSONAPI(method string, endpoint string, payload []byte, compressed bool, headers map[string]string, timeout int) (rd ResultData) {

	rd = ResultData{
		Result: InitResult(),
	}
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json"
	data, err := ExecuteAPI(method, endpoint, payload, compressed, headers, timeout)
	if err != nil {
		rd.Result.AddErr(err)
		return
	}
	if len(data) == 0 {
		return
	}

	// Create a temporary result data for unmarshalling purposes
	// The internal LiveNote field is not populated when unmarshalling
	trd := ResultData{}
	if err = json.Unmarshal(data, &trd); err != nil {
		rd.Result.AddErr(err)
		rd.Data = data // This is not marshable to resultdata, we'll try to send the real result
		return
	}

	// Assign temp to result
	rd.Data = trd.Data
	for _, m := range trd.Messages {
		if m == "" {
			continue
		}
		msgType := m[0:3]
		msg := m[3:]
		if strings.HasPrefix(msg, "[") {
			if endBr := strings.Index(msg, "]"); endBr != -1 {
				rd.MessagePrefix = msg[1:endBr]
				msg = msg[endBr+3:]
			}
		}

		switch msgType {
		case string(livenote.Warn):
			rd.Result.AddWarning(msg)
		case string(livenote.Error):
			rd.Result.AddError(msg)
		case string(livenote.Fatal):
			rd.Result.AddError(msg)
		case string(livenote.App):
			rd.Result.ln.AddAppMsg(msg)
		}
	}
	if !rd.Result.ln.HasErrors() {
		rd.Return(OK)
	}

	return
}

// ExecuteAPI wraps http operation that change or read data and returns a byte array
func ExecuteAPI(method string, endpoint string, payload []byte, compressed bool, headers map[string]string, timeout int) ([]byte, error) {

	nr, err := http.NewRequest(method, endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	nr.Close = true
	nr.Header.Set("User-Agent", fmt.Sprintf("stdutil.request/%s-%s", REQUEST_VERSION, REQUEST_MODIFIED))
	nr.Header.Set("Connection", "keep-alive")
	nr.Header.Set("Accept", "*/*")
	if compressed {
		nr.Header.Set("Accept-Encoding", "gzip, deflate, br")
		switch strings.ToUpper(nr.Method) {
		case "POST", "PUT", "PATCH":
			nr.Header.Add("Content-Encoding", "gzip")
		}
	}
	for k, v := range headers {
		k = strings.ToLower(k)
		if k != "cookie" {
			nr.Header.Set(k, v)
			continue
		}
		for _, nvs := range strings.Split(v, `;`) {
			if nv := strings.Split(nvs, `=`); len(nv) > 1 {
				nr.AddCookie(&http.Cookie{
					Name:  strings.TrimSpace(nv[0]),
					Value: strings.TrimSpace(nv[1]),
				})
			}
		}
	}
	if timeout == 0 {
		timeout = 30
	}
	cli := http.Client{
		Timeout:   time.Second * time.Duration(timeout),
		Transport: ct,
	}
	resp, err := cli.Do(nr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	var data []byte

	if !resp.Uncompressed {
		ce := resp.Header.Get("Content-Encoding")
		if ce == "gzip" {
			raw, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			gzr, err := gzip.NewReader(bytes.NewBuffer(raw))
			if err != nil {
				return nil, err
			}
			defer gzr.Close()
			for {
				uz := make([]byte, 1024)
				cnt, err := gzr.Read(uz)
				if err != nil {
					if !errors.Is(err, io.ErrUnexpectedEOF) {
						return nil, err
					}
					break
				}
				if cnt == 0 {
					break
				}
				data = append(data, uz[0:cnt]...)
			}
			return data, nil
		}
	}

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			return nil, err
		}
	}
	return data, nil
}

// PostJSON wraps http.Post with custom result
func PostJSON(endpoint string, payload []byte, gzipped bool, headers map[string]string) ResultData {
	return ExecuteJSONAPI("POST", endpoint, payload, gzipped, headers, reqTimeOut)
}

// PutJSON wraps http.Put with custom result
func PutJSON(endpoint string, payload []byte, gzipped bool, headers map[string]string) ResultData {
	return ExecuteJSONAPI("PUT", endpoint, payload, gzipped, headers, reqTimeOut)
}

// PatchJSON wraps http.Patch with custom result
func PatchJSON(endpoint string, payload []byte, gzipped bool, headers map[string]string) ResultData {
	return ExecuteJSONAPI("PATCH", endpoint, payload, gzipped, headers, reqTimeOut)
}

// GetJSON wraps http.Get with custom result
func GetJSON(endpoint string, headers map[string]string) ResultData {
	return ExecuteJSONAPI("GET", endpoint, nil, false, headers, reqTimeOut)
}

// DeleteJSON wraps http.Delete with custom result
func DeleteJSON(endpoint string, headers map[string]string) ResultData {
	return ExecuteJSONAPI("DELETE", endpoint, nil, false, headers, reqTimeOut)
}

// ParseQueryString parses the query string into a column value
func ParseQueryString(qs *string) NameValues {
	ret := NameValues{
		Pair: make(map[string]any),
	}
	rv, _ := url.ParseQuery(*qs)
	for k, v := range rv {
		ret.Pair[k] = strings.Join(v[:], ",")
	}
	return ret
}

// ParseRouteVars parses custom routes from a mux handler
func ParseRouteVars(r *http.Request) (Command []string, Key string) {
	cmd := make([]string, 0, 10)
	key := ""

	m := mux.CurrentRoute(r)
	pt, _ := m.GetPathTemplate()
	ptn := strings.Replace(r.URL.Path, pt, "", -1) // Trim the url by URL path. The remaining text will be the path to evaluate

	hasTrailingSlash := false
	if ptn != "" {
		hasTrailingSlash = ptn[len(ptn)-1:] == `/`
	}

	path := strings.FieldsFunc(ptn, func(c rune) bool {
		return c == '/'
	})

	pathlen := len(path)

	// If path length is 1, we might have a key.
	// But if the path is not a number, it might be a command
	if pathlen == 1 {
		if pth := path[0]; len(pth) > 0 {
			if hasTrailingSlash {
				cmd = append(cmd, strings.ToLower(pth))
			} else {
				key = pth
			}
		}
	}

	// If path length is greater than 1, we transfer all paths
	// to the cmd array except the last one. The last one will
	// be checked if it has a trailing slash
	if pathlen > 1 {
		for i, ck := range path {
			if i < pathlen-1 && len(ck) > 0 {
				cmd = append(cmd, strings.ToLower(ck))
			}
		}

		if pth := path[pathlen-1]; len(pth) > 0 {
			if hasTrailingSlash {
				cmd = append(cmd, strings.ToLower(pth))
			} else {
				key = pth
			}
		}
	}

	return cmd, key
}

// BuildAccessToken builds a JWT token
func BuildAccessToken(header *map[string]interface{}, claims *map[string]interface{}, secretkey string) string {
	clm := *claims
	var (
		usr, dom, app, dev string
		iss, sub, jti, tnt string
		exp, nbf, iat      int64
	)

	aud := jwt.Audience{}
	var ifc interface{}

	if ifc = clm["iss"]; ifc != nil {
		iss = ifc.(string)
	}

	if ifc = clm["sub"]; ifc != nil {
		sub = ifc.(string)
	}

	if ifc = clm["aud"]; ifc != nil {
		t := reflect.TypeOf(ifc)

		// check if this is a slice
		if t.Kind() == reflect.Slice {
			// check if what type of slice are the elements
			if t.Elem().Kind() == reflect.String {
				aud = ifc.([]string)
			}
		}

		// check if this is a string
		if t.Kind() == reflect.String {
			aud = jwt.Audience([]string{ifc.(string)})
		}
	}
	if ifc = clm["exp"]; ifc != nil {
		exp = ifc.(int64)
	}
	if ifc = clm["nbf"]; ifc != nil {
		nbf = ifc.(int64)
	}
	if ifc = clm["iat"]; ifc != nil {
		iat = ifc.(int64)
	}
	if ifc = clm["usr"]; ifc != nil {
		usr = ifc.(string)
	}
	if ifc = clm["dom"]; ifc != nil {
		dom = ifc.(string)
	}
	if ifc = clm["app"]; ifc != nil {
		app = ifc.(string)
	}
	if ifc = clm["dev"]; ifc != nil {
		dev = ifc.(string)
	}
	if ifc = clm["jti"]; ifc != nil {
		jti = ifc.(string)
	}
	if ifc = clm["tnt"]; ifc != nil {
		tnt = ifc.(string)
	}

	unixt := func(unixts int64) *jwt.Time {
		epoch := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
		tt := time.Unix(unixts, 0)
		if tt.Before(epoch) {
			tt = epoch
		}
		return &jwt.Time{Time: tt}
	}

	pl := CustomPayload{
		Payload: jwt.Payload{
			Issuer:         iss,
			Subject:        sub,
			Audience:       aud,
			ExpirationTime: unixt(exp),
			NotBefore:      unixt(nbf),
			IssuedAt:       unixt(iat),
			JWTID:          jti,
		},
		UserName:      usr,
		Domain:        dom,
		ApplicationID: app,
		DeviceID:      dev,
		TenantID:      tnt,
	}

	HMAC := jwt.NewHS256([]byte(secretkey))
	token, err := jwt.Sign(pl, HMAC)
	if err != nil {
		return ""
	}

	return string(token)
}

// GetRequestVarsOnly get request variables
func GetRequestVarsOnly(r *http.Request) RequestVars {

	rv := &RequestVars{}
	const (
		mulpart string = "multipart/form-data"
		furlenc string = "application/x-www-form-urlencoded"
	)

	rv.Method = strings.ToUpper(r.Method)
	ctype := strings.Split(r.Header.Get("Content-Type"), ";")
	c1 := strings.TrimSpace(ctype[0])

	useBody := (c1 != furlenc && c1 != mulpart) && (rv.IsPostOrPut() || rv.IsDelete())

	if useBody {
		// We are receiving body as bytes to Unmarshall later depending on the type
		b := func() []byte {
			if r.Body != nil {
				b, _ := io.ReadAll(r.Body)
				defer r.Body.Close()
				return b
			}
			return []byte{}
		}
		if rv.Body = b(); rv.Body != nil {
			rv.HasBody = len(rv.Body) > 0
		}
	}

	// Query Strings
	rv.Variables.QueryString = ParseQueryString(&r.URL.RawQuery)
	rv.Variables.HasQueryString = len(rv.Variables.QueryString.Pair) > 0
	rv.Variables.IsMultipart = (c1 == mulpart)

	if rv.Variables.IsMultipart {
		r.ParseMultipartForm(30 << 20)
	} else {
		r.ParseForm()
	}

	// Get Form data
	rv.Variables.FormData = NameValues{
		Pair: make(map[string]any),
	}

	for k, v := range r.PostForm {
		rv.Variables.FormData.Pair[k] = strings.Join(v[:], ",")
	}
	rv.Variables.HasFormData = len(rv.Variables.FormData.Pair) > 0

	// Get route commands
	rv.Variables.Command, rv.Variables.Key = ParseRouteVars(r)

	return *rv
}

// ValidateJWT validates JWT and returns information
func ValidateJWT(r *http.Request, secretKey string, validateTimes bool) (*JWTInfo, error) {

	var (
		jwtfromck,
		jwth string
		jwtp []string
	)

	// Get Authorization header
	if jwth = r.Header.Get("Authorization"); len(jwth) == 0 {
		return nil, fmt.Errorf(`authorization header not set`)
	}
	if jwtp = strings.Split(jwth, " "); len(jwtp) < 2 {
		return nil, fmt.Errorf(`invalid authorization header`)
	}
	if !strings.EqualFold(strings.TrimSpace(jwtp[0]), "bearer") {
		return nil, fmt.Errorf(`invalid authorization bearer`)
	}
	if jwtfromck = strings.TrimSpace(jwtp[1]); len(jwtfromck) == 0 {
		return nil, fmt.Errorf(`invalid authorization token`)
	}
	return ParseJWT(jwtfromck, secretKey, validateTimes)
}

// ParseJWT validates, parses JWT and returns information
func ParseJWT(token, secretKey string, validateTimes bool) (*JWTInfo, error) {

	ji := &JWTInfo{}

	if len(secretKey) == 0 {
		return ji, fmt.Errorf(`secret key not set`)
	}

	var (
		pl  CustomPayload
		err error
	)

	// Parse JWT
	HMAC := jwt.NewHS256([]byte(secretKey))
	now := time.Now()

	// Validate claims "iat", "exp" and "aud".
	if validateTimes {
		iatValidator := jwt.IssuedAtValidator(now)
		expValidator := jwt.ExpirationTimeValidator(now)
		nbfValidator := jwt.NotBeforeValidator(now)

		// Use jwt.ValidatePayload to build a jwt.VerifyOption.
		// Validators are run in the order informed.
		validator := jwt.ValidatePayload(&pl.Payload, iatValidator, expValidator, nbfValidator)
		_, err = jwt.Verify([]byte(token), HMAC, &pl, validator)
	} else {
		_, err = jwt.Verify([]byte(token), HMAC, &pl)
	}

	if err != nil {
		return ji, err
	}

	ji.Audience = pl.Audience
	ji.UserName = pl.UserName
	ji.Domain = pl.Domain
	ji.DeviceID = pl.DeviceID
	ji.ApplicationID = pl.ApplicationID
	ji.TenantID = pl.TenantID
	ji.Raw = token

	ji.Valid = true

	return ji, nil
}

// GetRequestVars requests variables and return JWT validation result
func GetRequestVars(r *http.Request, secretKey string, validateTimes bool) (RequestVars, error) {

	rv := GetRequestVarsOnly(r)
	rv.Token = nil

	// silently ignore OPTION methid
	if strings.EqualFold(r.Method, "OPTION") {
		return rv, nil
	}

	ji, err := ValidateJWT(r, secretKey, validateTimes)
	if err != nil {
		return rv, err
	}
	rv.Token = ji

	return rv, nil
}
