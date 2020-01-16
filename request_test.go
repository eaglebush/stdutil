package stdutil

import (
	"fmt"
	"testing"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
)

func TestExecuteAPIPOST(t *testing.T) {
	payload := `{
		"username": "zaldy.baguinon",
		"sessionid": "eyJhbGciOiJIUzI1NiIsImRldmljZWlkIjoiIiwiZG9tYWluIjoiIiwidHlwIjoiSldUIiwidXNlciI6InphbGR5LmJhZ3Vpbm9uIn0.eyJuYmYiOjE1NzgyNzU3MDV9.8NbRqiIIQ6Kx03Zo_aOyf_5rFnhYQtM8O990TEv0_aM"
	}`

	exapi := PostJSON("http://hulk.vdimdci.com.ph/api/appshub/auth/svalid/", []byte(payload), nil)
	fmt.Println(exapi)
}

func TestExecuteAPIGET(t *testing.T) {
	hdr := make(map[string]string, 1)
	hdr["Cookie"] = "APPSHUB-WF-login=zaldy.baguinon; APPSHUB-WF-appdomain=MDCI"

	exapi := GetJSON("http://hulk.vdimdci.com.ph/api/appshub/user/19", hdr)
	fmt.Printf("%v", exapi)
}

func TestJWTParse(t *testing.T) {
	jwtfromck := []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJBUFBTSFVCLUFVVEgiLCJuYmYiOjE1NzkxNjE4OTgsInVzciI6InphbGR5LmJhZ3Vpbm9uIiwiZG9tIjoiTURDSSIsImRldiI6IjFXU3hhQ0h1V2VYSDREN0tXd0dZSkZlTTRwRiIsImFwcCI6IkFQUFNIVUItQVVUSCJ9.SKU6lfcVO5JAk81zvDYxvcOl6IUY7Kg_QJW4IFS3xso")

	type CustomPayload struct {
		jwt.Payload
		UserName      string `json:"usr,omitempty"`
		Domain        string `json:"dom,omitempty"`
		ApplicationID string `json:"app,omitempty"`
		DeviceID      string `json:"dev,omitempty"`
	}

	HMAC := jwt.NewHS256([]byte("thisisanhmacsecretkey"))

	var pl CustomPayload

	now := time.Now()

	// Validate claims "iat", "exp" and "aud".
	iatValidator := jwt.IssuedAtValidator(now)
	expValidator := jwt.ExpirationTimeValidator(now)

	// Use jwt.ValidatePayload to build a jwt.VerifyOption.
	// Validators are run in the order informed.
	validatePayload := jwt.ValidatePayload(&pl.Payload, iatValidator, expValidator)

	_, err := jwt.Verify(jwtfromck, HMAC, &pl, validatePayload)
	if err != nil {
		fmt.Println(err)
	}

}
