package stdutil

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
	"github.com/gorilla/mux"
)

//CustomVars - command struct
type CustomVars struct {
	Command        []string
	Key            string
	QueryString    NameValues
	HasQueryString bool
	FormData       NameValues
	HasFormData    bool
	IsMultipart    bool
	DecodedCommand NameValues
}

// CustomPayload - payload for JWT
type CustomPayload struct {
	jwt.Payload
	UserName      string `json:"usr,omitempty"`
	Domain        string `json:"dom,omitempty"`
	ApplicationID string `json:"app,omitempty"`
	DeviceID      string `json:"dev,omitempty"`
}

//RequestVars - contains necessary request variables
type RequestVars struct {
	Method             string
	Variables          CustomVars
	Body               []byte
	Cookies            map[string]string
	HasBody            bool
	ValidAuthToken     bool
	TokenRaw           string
	TokenApplicationID string
	TokenAudience      []string
	TokenDeviceID      string
	TokenUserName      string
	TokenDomain        string
}

// ResultData - a result structure and a generic data
type ResultData struct {
	Result
	Data json.RawMessage
}

//IsGet - a shortcut method to check if the request is a GET
func (rv *RequestVars) IsGet() bool {
	return rv.Method == "GET"
}

//IsPost - a shortcut method to check if the request is a POST
func (rv *RequestVars) IsPost() bool {
	return rv.Method == "POST"
}

//IsPut - a shortcut method to check if the request is a PUT
func (rv *RequestVars) IsPut() bool {
	return rv.Method == "PUT"
}

//IsDelete - a shortcut method to check if the request is a DELETE
func (rv *RequestVars) IsDelete() bool {
	return rv.Method == "DELETE"
}

// IsHead - a shortcut method to check if the request is a HEAD
func (rv *RequestVars) IsHead() bool {
	return rv.Method == "HEAD"
}

// IsOptions - a shortcut method to check if the request is OPTIONS
func (rv *RequestVars) IsOptions() bool {
	return rv.Method == "OPTIONS"
}

// ExecuteJSONAPI - a wrapper for http operation that can change or read data that returns a custom result
func ExecuteJSONAPI(method string, endpoint string, payload []byte, gzipped bool, headers map[string]string, timeout int) (rd ResultData) {

	rd = ResultData{}
	rd.Result = InitResult()

	nr, err := http.NewRequest(method, endpoint, bytes.NewBuffer(payload))
	if err != nil {
		rd.Result.Messages = append(rd.Result.Messages, err.Error())
		return
	}

	if headers != nil {
		for k, v := range headers {

			k = strings.ToLower(k)

			switch k {
			case "cookie":

				// split values with semi-colons
				cnvs := strings.Split(v, `;`)

				for _, nvs := range cnvs {
					nv := strings.Split(nvs, `=`)

					if len(nv) > 1 {
						nv[0] = strings.TrimSpace(nv[0])
						nv[1] = strings.TrimSpace(nv[1])

						nr.AddCookie(&http.Cookie{
							Name:  nv[0],
							Value: nv[1],
						})
					}
				}

			default:
				nr.Header.Add(k, v)
			}
		}
	}

	// Standard header
	nr.Header.Add("Content-Type", "application/json")
	nr.Header.Add("Accept-Encoding", "gzip")
	if gzipped {
		nr.Header.Add("Content-Encoding", "gzip")
	}

	if timeout == 0 {
		timeout = 30
	}

	cli := http.Client{Timeout: time.Second * time.Duration(timeout)}
	resp, err := cli.Do(nr)
	if err != nil {
		rd.Result.Messages = append(rd.Result.Messages, err.Error())
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		rd.Result.Messages = append(rd.Result.Messages, err.Error())
		return
	}

	if len(data) != 0 {
		if err = json.Unmarshal(data, &rd); err != nil {
			rd.Result.Messages = append(rd.Result.Messages, err.Error())
			return
		}
	}

	return
}

// PostJSON - a wrapper for http.Post with custom result
func PostJSON(endpoint string, payload []byte, gzipped bool, headers map[string]string) ResultData {
	return ExecuteJSONAPI("POST", endpoint, payload, gzipped, headers, 30)
}

// PutJSON - a wrapper for http.Put with custom result
func PutJSON(endpoint string, payload []byte, gzipped bool, headers map[string]string) ResultData {
	return ExecuteJSONAPI("PUT", endpoint, payload, gzipped, headers, 30)
}

// GetJSON - a wrapper for http.Get with returns with a custom result
func GetJSON(endpoint string, headers map[string]string) ResultData {
	return ExecuteJSONAPI("GET", endpoint, nil, false, headers, 30)
}

// DeleteJSON - a wrapper for http.Delete with custom result
func DeleteJSON(endpoint string, headers map[string]string) ResultData {
	return ExecuteJSONAPI("DELETE", endpoint, nil, false, headers, 30)
}

//ParseQueryString - parse the query string into a column value
func ParseQueryString(qs *string) NameValues {
	rv, _ := url.ParseQuery(*qs)

	ret := NameValues{}
	ret.Pair = make([]NameValue, 0)
	for k, v := range rv {
		ret.Pair = append(ret.Pair, NameValue{k, strings.Join(v[:], ",")})
	}

	return ret
}

// ParseRouteVars - parse custom routes from a mux handler
func ParseRouteVars(r *http.Request) (Command []string, Key string) {
	cmd := make([]string, 0)
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

	/* If path length is 1, we might have a key. But if the path is not a number, it might be a command  */
	if pathlen == 1 {
		pth := path[0]
		if len(pth) > 0 {
			if hasTrailingSlash {
				cmd = append(cmd, strings.ToLower(pth))
			} else {
				key = pth
			}
		}
	}

	/* If path length is greater than 1, we transfer all paths to the cmd array except the last one. The last one will be checked if it has a trailing slash */
	if pathlen > 1 {
		for i, ck := range path {
			if i < pathlen-1 && len(ck) > 0 {
				cmd = append(cmd, strings.ToLower(ck))
			}
		}

		pth := path[pathlen-1]
		if len(pth) > 0 {
			if hasTrailingSlash {
				cmd = append(cmd, strings.ToLower(pth))
			} else {
				key = pth //key will not be set to lower case
			}
		}
	}

	return cmd, key
}

// BuildAccessToken - build a JWT token
func BuildAccessToken(header *map[string]interface{}, claims *map[string]interface{}, secretkey string) string {
	clm := *claims

	iss := ""
	sub := ""
	aud := jwt.Audience{}
	exp := int64(0)
	nbf := int64(0)
	iat := int64(0)
	usr := ""
	dom := ""
	app := ""
	dev := ""

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
		},
		UserName:      usr,
		Domain:        dom,
		ApplicationID: app,
		DeviceID:      dev,
	}

	HMAC := jwt.NewHS256([]byte(secretkey))

	token, err := jwt.Sign(pl, HMAC)
	if err != nil {
		return ""
	}

	return string(token)
}

// GetRequestVars - get request variables
func GetRequestVars(r *http.Request, secretkey string) RequestVars {

	rv := &RequestVars{}

	const mp string = "multipart/form-data"

	rv.Method = strings.ToUpper(r.Method)
	ctype := strings.Split(r.Header.Get("Content-Type"), ";")
	c1 := strings.TrimSpace(ctype[0])
	useBody := (c1 != "application/x-www-form-urlencoded" && c1 != mp) && (rv.IsPost() || rv.IsPut() || rv.IsDelete())

	if useBody {
		// We are receiving body as bytes to Unmarshall later depending on the type
		b := func() []byte {
			if r.Body != nil {
				b, _ := ioutil.ReadAll(r.Body)
				defer r.Body.Close()
				return b
			}
			return []byte{}
		}
		rv.Body = b()

		if rv.Body != nil {
			rv.HasBody = len(rv.Body) > 0
		}
	}

	// Query Strings
	rv.Variables.QueryString = ParseQueryString(&r.URL.RawQuery)
	rv.Variables.HasQueryString = len(rv.Variables.QueryString.Pair) > 0
	rv.Variables.IsMultipart = (c1 == mp)

	if rv.Variables.IsMultipart {
		r.ParseMultipartForm(30 << 20)
	} else {
		r.ParseForm()
	}

	// Get Form data
	rv.Variables.FormData = NameValues{}
	rv.Variables.FormData.Pair = make([]NameValue, 0)
	for k, v := range r.PostForm {
		rv.Variables.FormData.Pair = append(rv.Variables.FormData.Pair, NameValue{
			k, strings.Join(v[:], ","),
		})
	}
	rv.Variables.HasFormData = len(rv.Variables.FormData.Pair) > 0

	// Get route commands
	rv.Variables.Command, rv.Variables.Key = ParseRouteVars(r)

	jwtfromck := ""

	// Get Authorization header
	if jwth := r.Header.Get("Authorization"); len(jwth) > 0 {
		if jwtp := strings.Split(jwth, " "); len(jwtp) > 1 {
			if strings.ToLower(strings.TrimSpace(jwtp[0])) == "bearer" {
				jwtfromck = strings.TrimSpace(jwtp[1])
			}
		}
	}

	// Parse JWT
	if len(jwtfromck) > 0 {

		var pl CustomPayload

		HMAC := jwt.NewHS256([]byte(secretkey))

		// Commented as this is not yet implemented
		//now := time.Now()

		// // Validate claims "iat", "exp" and "aud".
		// iatValidator := jwt.IssuedAtValidator(now)
		// expValidator := jwt.ExpirationTimeValidator(now)
		// nbfValidator := jwt.NotBeforeValidator(now)

		// // Use jwt.ValidatePayload to build a jwt.VerifyOption.
		// // Validators are run in the order informed.
		// validatePayload := jwt.ValidatePayload(&pl.Payload, iatValidator, expValidator, nbfValidator)
		// if _, err := jwt.Verify([]byte(jwtfromck), HMAC, &pl, validatePayload); err == nil {

		if _, err := jwt.Verify([]byte(jwtfromck), HMAC, &pl); err == nil {
			rv.TokenAudience = pl.Audience

			rv.TokenUserName = pl.UserName
			rv.TokenDomain = pl.Domain
			rv.TokenDeviceID = pl.DeviceID
			rv.TokenApplicationID = pl.ApplicationID

			rv.TokenRaw = jwtfromck
			rv.ValidAuthToken = true
		}
	}

	return *rv
}

// FirstCommand - get first command from route
func (cv *CustomVars) FirstCommand() string {
	_, ret := cv.GetCommand(0)
	return ret
}

// SecondCommand - get second command from route
func (cv *CustomVars) SecondCommand() string {
	_, ret := cv.GetCommand(1)
	return ret
}

// ThirdCommand - get third command from route
func (cv *CustomVars) ThirdCommand() string {
	_, ret := cv.GetCommand(2)
	return ret
}

// LastCommand - get third command from route
func (cv *CustomVars) LastCommand() string {
	_, ret := cv.GetCommand(uint(len(cv.Command) - 1))
	return ret
}

// GetCommand - get command by index
func (cv *CustomVars) GetCommand(index uint) (exists bool, value string) {
	lenc := uint(len(cv.Command))

	// if there's no command, return at once
	if lenc == 0 {
		return false, ""
	}

	// if the index is greater than the length of the array
	if index > lenc-1 {
		return false, ""
	}

	return true, cv.Command[index]
}
