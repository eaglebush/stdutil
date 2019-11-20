package stdutil

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

//CustomVars - command struct
type CustomVars struct {
	Command        string
	Key            string
	QueryString    NameValues
	HasQueryString bool
	FormData       NameValues
	HasFormData    bool
	IsMultipart    bool
	DecodedCommand NameValues
}

//RequestVars - contains necessary request variables
type RequestVars struct {
	Login          string
	Domain         string
	Method         string
	Variables      CustomVars
	Body           []byte
	Token          jwt.Token
	Cookies        map[string]string
	HasBody        bool
	ValidAuthToken bool
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

//IsHead - a shortcut method to check if the request is a HEAD
func (rv *RequestVars) IsHead() bool {
	return rv.Method == "HEAD"
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

//ParseRouteVars - parse custom routes from a mux handler
func ParseRouteVars(r *http.Request) (Command string, Key string) {
	cmd := ""
	key := ""
	/*
		1. Last part of the path should be the key, if this is not numeric, it will be the command
		2. If the total number of paths is 4, we check if the key is numeric, if it is, it will be the key.
	*/

	m := mux.CurrentRoute(r)
	pt, _ := m.GetPathTemplate()
	ptn := strings.Replace(r.URL.Path, pt, "", -1) // Trim the url by URL path. The remaining text will be the path to evaluate

	path := strings.FieldsFunc(ptn, func(c rune) bool {
		return c == '/'
	})

	/* If path length is 1, we might have a key. But if the path is not a number, it might be a command  */
	if len(path) == 1 {
		key = path[0]
	}

	/* If path length is 2, we might have a key and a command */
	if len(path) == 2 {
		cmd = path[0] /* the second to the last would be the command */
		key = path[1]
	}

	return cmd, key
}

//BuildAccessToken - build a JWT token
func BuildAccessToken(header *map[string]interface{}, claims *map[string]interface{}, HMAC string) string {
	/*
		token := jwt.NewWithClaims(jwt.SigningMethodHMAC, jwt.MapClaims{
			"foo": "bar",
			"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		})
	*/
	// Create the Claims
	stdclaims := jwt.StandardClaims{}

	for k, v := range *claims {
		switch k {
		case "aud":
			stdclaims.Audience = v.(string)
		case "exp":
			stdclaims.ExpiresAt = v.(int64)
		case "jti":
			stdclaims.Id = v.(string)
		case "iat":
			stdclaims.IssuedAt = v.(int64)
		case "iss":
			stdclaims.Issuer = v.(string)
		case "nbf":
			stdclaims.NotBefore = v.(int64)
		case "sub":
			stdclaims.Subject = v.(string)
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, stdclaims)
	token.Header = *header

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(HMAC))

	if err != nil {
		return ""
	}

	return tokenString
}

// GetRequestVars - get request variables
func GetRequestVars(r *http.Request, ApplicationID string, HMAC string) RequestVars {
	rv := &RequestVars{}

	const mp string = "multipart/form-data"

	rv.Method = strings.ToUpper(r.Method)
	ctype := strings.Split(r.Header.Get("Content-Type"), ";")
	c1 := strings.TrimSpace(ctype[0])
	useBody := (c1 != "application/x-www-form-urlencoded" && c1 != mp) && (rv.IsPost() || rv.IsPut())

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

	rv.Variables.FormData = NameValues{}
	rv.Variables.FormData.Pair = make([]NameValue, 0)
	for k, v := range r.PostForm {
		rv.Variables.FormData.Pair = append(rv.Variables.FormData.Pair, NameValue{k, strings.Join(v[:], ",")})
	}
	rv.Variables.HasFormData = len(rv.Variables.FormData.Pair) > 0

	rv.Variables.Command, rv.Variables.Key = ParseRouteVars(r)

	jwtfromck := ""
	rv.ValidAuthToken = false

	// Get cookies that matters
	rv.Cookies = make(map[string]string)
	for _, c := range r.Cookies() {
		rv.Cookies[c.Name] = c.Value
		// Set login value if met
		if c.Name == ApplicationID+"-login" {
			rv.Login = c.Value
		}

		if c.Name == ApplicationID+"-appdomain" {
			rv.Domain = c.Value
		}

		if c.Name == ApplicationID+"-sessionid" {
			jwtfromck = c.Value
		}
	}

	// Get JWT from request headers if the cookie has none
	if jwtfromck == "" {
		jwth := r.Header.Get("Authorization")
		if len(jwth) > 0 {
			jwtp := strings.Split(jwth, " ")
			if len(jwtp) > 1 {
				if strings.ToLower(strings.TrimSpace(jwtp[0])) == "bearer" {
					jwtfromck = strings.TrimSpace(jwtp[1])
				}
			}
		}
	}

	if len(jwtfromck) > 0 {
		token, _ := jwt.Parse(jwtfromck, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(HMAC), nil
		})

		if token.Valid {
			rv.Token = *token
		}
	}

	// Validate the user name from cookie against the retrieved token. Expiry Data could also be validated manually here
	if rv.Token.Valid {
		tokuser := rv.Token.Header["user"]
		if tokuser == nil {
			rv.ValidAuthToken = false
		}

		tokdom := rv.Token.Header["domain"]
		if tokdom == nil {
			rv.ValidAuthToken = false
		}

		rv.ValidAuthToken = tokuser == rv.Login && tokdom == rv.Domain
	}

	return *rv
}
