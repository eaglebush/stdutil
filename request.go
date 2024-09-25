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
	"sync"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
	"github.com/gorilla/mux"
	"github.com/narsilworks/livenote"
)

const (
	REQUEST_VERSION  string = "1.0.0.0"
	REQUEST_MODIFIED string = "16052024"
)

var (
	reqTimeOut int
	ct         *http.Transport
)

type (
	// CustomPayload - payload for JWT
	CustomPayload struct {
		jwt.Payload
		UserName      string `json:"usr,omitempty"` // Username payload for JWT
		Domain        string `json:"dom,omitempty"` // Domain payload for JWT
		ApplicationID string `json:"app,omitempty"` // Application payload for JWT
		DeviceID      string `json:"dev,omitempty"` // Device id payload for JWT
		TenantID      string `json:"tnt,omitempty"` // Tenant id payload for JWT
	}
	// ResultData - a result structure and a JSON raw message
	ResultData struct {
		Result
		Data json.RawMessage `json:"data"`
	}
	// RequestParam for <REST verb>Api request functions
	RequestParam struct {
		TimeOut    int               // Request time out
		Compressed bool              // Compressed
		Headers    map[string]string // Headers for the request
		Mutex      *sync.RWMutex     // Mutex lock for header modification
	}
	// RequestOption for <REST verb>Api request functions
	RequestOption func(opt *RequestParam) error
)

func init() {
	reqTimeOut = 30
	ct = http.DefaultTransport.(*http.Transport).Clone()
	ct.MaxIdleConns = 100
	ct.MaxConnsPerHost = 100
	ct.MaxIdleConnsPerHost = 100
}

// SetRequestTimeOut sets the new timeout value
func SetRequestTimeout(timeOut int) {
	reqTimeOut = timeOut
}

// ExecuteJsonApi wraps http operation that change or read data and returns a custom result
func ExecuteJsonApi(method string, endPoint string, payload []byte, compressed bool, header map[string]string, timeOut int, rw *sync.RWMutex) (rd ResultData) {
	rd = ResultData{
		Result: InitResult(),
	}
	if header == nil {
		header = make(map[string]string)
	}
	if rw == nil {
		rw = &sync.RWMutex{}
	}
	SafeMapWrite(&header, "Content-Type", "application/json", rw)
	data, err := ExecuteApi(method, endPoint, payload, compressed, header, timeOut)
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
	rd.Return(Status(trd.Status))
	for _, m := range trd.Messages {
		if m == "" {
			continue
		}
		msgType := m[0:3]
		msg := m[3:]
		if strings.HasPrefix(msg, "[") {
			if endBr := strings.Index(msg, "]"); endBr != -1 {
				rd.ln.Prefix = msg[1:endBr]
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

	return
}

// ExecuteApi wraps http operation that change or read data and returns a byte array
//
// On headers:
//   - Content-Type: If this header is not set, it defaults to "application/json"//
//   - Content-Encoding: If compressed is true, it is set to "gzip"
func ExecuteApi(method string, endPoint string, payload []byte, compressed bool, header map[string]string, timeOut int) ([]byte, error) {
	nr, err := http.NewRequest(method, endPoint, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	nr.Close = true
	nr.Header.Set(
		"User-Agent",
		fmt.Sprintf("com.github.eaglebush.stdutil.request/%s-%s",
			REQUEST_VERSION, REQUEST_MODIFIED))
	nr.Header.Set("Connection", "keep-alive")
	nr.Header.Set("Accept", "*/*")
	if ct := nr.Header.Get("Content-Type"); ct == "" {
		nr.Header.Set("Content-Type", "application/json")
	}
	if compressed {
		nr.Header.Set("Accept-Encoding", "gzip, deflate, br")
		switch strings.ToUpper(nr.Method) {
		case "POST", "PUT", "PATCH":
			nr.Header.Add("Content-Encoding", "gzip")
		}
	}
	for k, v := range header {
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
	if timeOut == 0 {
		timeOut = 30
	}
	cli := http.Client{
		Timeout:   time.Second * time.Duration(timeOut),
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

// GetJson wraps http.Get and gets a raw json message data
func GetJson(endpoint string, headers map[string]string, rw *sync.RWMutex) ResultData {
	return ExecuteJsonApi("GET", endpoint, nil, false, headers, reqTimeOut, rw)
}

// DeleteJson wraps http.Delete and gets a raw json message data
func DeleteJson(endpoint string, headers map[string]string, rw *sync.RWMutex) ResultData {
	return ExecuteJsonApi("DELETE", endpoint, nil, false, headers, reqTimeOut, rw)
}

// PostJson wraps http.Post and gets a raw json message data
func PostJson(endpoint string, payload []byte, gzipped bool, headers map[string]string, rw *sync.RWMutex) ResultData {
	return ExecuteJsonApi("POST", endpoint, payload, gzipped, headers, reqTimeOut, rw)
}

// PutJson wraps http.Put and gets a raw json message data
func PutJson(endpoint string, payload []byte, gzipped bool, headers map[string]string, rw *sync.RWMutex) ResultData {
	return ExecuteJsonApi("PUT", endpoint, payload, gzipped, headers, reqTimeOut, rw)
}

// PatchJson wraps http.Patch and gets a raw json message data
func PatchJson(endpoint string, payload []byte, gzipped bool, headers map[string]string, rw *sync.RWMutex) ResultData {
	return ExecuteJsonApi("PATCH", endpoint, payload, gzipped, headers, reqTimeOut, rw)
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
	var (
		c1 string
	)
	const (
		mulpart string = "multipart/form-data"
		furlenc string = "application/x-www-form-urlencoded"
	)
	rv := &RequestVars{
		Method: strings.ToUpper(r.Method),
	}
	if ctype := strings.Split(r.Header.Get("Content-Type"), ";"); len(ctype) > 0 {
		c1 = strings.TrimSpace(ctype[0])
	}
	if useBody := (c1 != furlenc && c1 != mulpart) && (rv.IsPostOrPut() || rv.IsDelete()); useBody {
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
	if len(secretKey) == 0 {
		return nil, fmt.Errorf(`secret key not set`)
	}
	var (
		pl  CustomPayload
		err error
	)

	// Parse JWT
	HMAC := jwt.NewHS256([]byte(secretKey))

	// Validate claims "iat", "exp" and "aud".
	if validateTimes {
		now := time.Now()
		// Use jwt.ValidatePayload to build a jwt.VerifyOption.
		// Validators are run in the order informed.
		validator := jwt.ValidatePayload(
			&pl.Payload,
			jwt.IssuedAtValidator(now),
			jwt.ExpirationTimeValidator(now),
			jwt.NotBeforeValidator(now))
		_, err = jwt.Verify([]byte(token), HMAC, &pl, validator)
	} else {
		_, err = jwt.Verify([]byte(token), HMAC, &pl)
	}
	if err != nil {
		return nil, err
	}
	return &JWTInfo{
		Audience:      pl.Audience,
		UserName:      pl.UserName,
		Domain:        pl.Domain,
		DeviceID:      pl.DeviceID,
		ApplicationID: pl.ApplicationID,
		TenantID:      pl.TenantID,
		Raw:           token,
		Valid:         true,
	}, nil
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

func getJsonConverted[T any](rslt *ResultData) ResultAny[T] {
	var data T
	if !rslt.OK() {
		return ResultAny[T]{
			Result: rslt.Result,
			Data:   data,
		}
	}
	if len(rslt.Data) == 0 {
		return ResultAny[T]{
			Result: InitResult(
				NameValue[string]{
					Name: "status", Value: string(EXCEPTION),
				},
				NameValue[string]{
					Name: "message", Value: "No data retrieved",
				},
			),
			Data: data,
		}
	}
	if err := json.Unmarshal(rslt.Data, &data); err != nil {
		return ResultAny[T]{
			Result: InitResult(
				NameValue[string]{
					Name: "status", Value: string(EXCEPTION),
				},
				NameValue[string]{
					Name: "message", Value: err.Error(),
				},
			),
			Data: data,
		}
	}
	return ResultAny[T]{
		Result: InitResult(
			NameValue[string]{
				Name: "status", Value: rslt.Status,
			},
		),
		Data: data,
	}
}

// TimeOut sets the request timeout as an option
//
// This is used with <REST verb>Api functions
func TimeOut(timeOut int) RequestOption {
	return func(rp *RequestParam) error {
		rp.TimeOut = timeOut
		return nil
	}
}

// Compressed sets the request compression as an option
//
// This is used with <REST verb>Api functions
func Compressed(compressed bool) RequestOption {
	return func(rp *RequestParam) error {
		rp.Compressed = compressed
		return nil
	}
}

// Headers adds request headers as an option
//
// This is used with <REST verb>Api functions
func Headers(hdr map[string]string, mut *sync.RWMutex) RequestOption {
	return func(rp *RequestParam) error {
		rp.Headers = hdr
		rp.Mutex = mut
		return nil
	}
}

// CreateApi posts data on an API endpoint and converts the returned data into a resulting type
func CreateApi[T any, U any](url string, pl U, gzpd bool, hdrs map[string]string, rw *sync.RWMutex) ResultAny[T] {
	b, err := json.Marshal(pl)
	if err != nil {
		return ResultAny[T]{
			Result: InitResult(
				NameValue[string]{
					Name:  "message",
					Value: err.Error(),
				},
			),
		}
	}
	rd := ExecuteJsonApi("POST", url, b, gzpd, hdrs, reqTimeOut, rw)
	return getJsonConverted[T](&rd)
}

// ReadApi retrieves data on an API endpoint and converts the returned data into a resulting type
func ReadApi[T any](url string, opts ...RequestOption) ResultAny[T] {
	rp := RequestParam{}
	for _, o := range opts {
		if o == nil {
			continue
		}
		o(&rp)
	}
	rd := ExecuteJsonApi("GET", url, nil, rp.Compressed, rp.Headers, rp.TimeOut, rp.Mutex)
	return getJsonConverted[T](&rd)
}

// UpdateApi updates data on an API endpoint and converts the returned data into a resulting type
func UpdateApi[T any, U any](url string, pl U, opts ...RequestOption) ResultAny[T] {
	b, err := json.Marshal(pl)
	if err != nil {
		return ResultAny[T]{
			Result: InitResult(
				NameValue[string]{
					Name:  "message",
					Value: err.Error(),
				},
			),
		}
	}
	rp := RequestParam{}
	for _, o := range opts {
		if o == nil {
			continue
		}
		o(&rp)
	}
	rd := ExecuteJsonApi("PUT", url, b, rp.Compressed, rp.Headers, rp.TimeOut, rp.Mutex)
	return getJsonConverted[T](&rd)
}

// DeleteApi deletes data on an API endpoint and converts the returned data into a resulting type
func DeleteApi[T any](url string, opts ...RequestOption) ResultAny[T] {
	rp := RequestParam{}
	for _, o := range opts {
		if o == nil {
			continue
		}
		o(&rp)
	}
	rd := ExecuteJsonApi("DELETE", url, nil, rp.Compressed, rp.Headers, rp.TimeOut, rp.Mutex)
	return getJsonConverted[T](&rd)
}

// PatchApi patches data on an API endpoint and converts the returned data into a resulting type
func PatchApi[T any, U any](url string, pl U, opts ...RequestOption) ResultAny[T] {
	b, err := json.Marshal(pl)
	if err != nil {
		return ResultAny[T]{
			Result: InitResult(
				NameValue[string]{
					Name:  "message",
					Value: err.Error(),
				},
			),
		}
	}
	rp := RequestParam{}
	for _, o := range opts {
		if o == nil {
			continue
		}
		o(&rp)
	}
	rd := ExecuteJsonApi("PATCH", url, b, rp.Compressed, rp.Headers, rp.TimeOut, rp.Mutex)
	return getJsonConverted[T](&rd)
}
