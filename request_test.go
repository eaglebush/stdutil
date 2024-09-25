package stdutil

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
)

func TestExecuteAPIPOST(t *testing.T) {
	payload := `{
		"username": "zaldy.baguinon",
		"sessionid": "eyJhbGciOiJIUzI1NiIsImRldmljZWlkIjoiIiwiZG9tYWluIjoiIiwidHlwIjoiSldUIiwidXNlciI6InphbGR5LmJhZ3Vpbm9uIn0.eyJuYmYiOjE1NzgyNzU3MDV9.8NbRqiIIQ6Kx03Zo_aOyf_5rFnhYQtM8O990TEv0_aM"
	}`

	exapi := PostJson("https://appcore.vdimdci.com.ph/api/auth/svalid/", []byte(payload), false, nil, nil)
	fmt.Println(exapi)
}

func TestExecuteAPIGET(t *testing.T) {
	hdr := make(map[string]string)
	hdr["Cookie"] = "APPSHUB-WF-login=zaldy.baguinon; APPSHUB-WF-appdomain=MDCI"
	hdr["Authorization"] = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJLZWFudS1VcGxvYWRlciIsImV4cCI6MCwibmJmIjoxNjc1MzE0MTg0LCJpYXQiOjAsInVzciI6ImphbWVzLmx1bWliYW9AbWRjaS5jb20ucGgiLCJkb20iOiJNRENJIiwiYXBwIjoiS2VhbnUtVXBsb2FkZXIiLCJkZXYiOiIyS1JzS3Z4Y2NuOUp0RjNxbDIxMmN1MmhwS1MifQ.961xUrBObQfN6fkO_s7OYhFTqKC_aMrr1OKVwvPhkLU"

	exapi := GetJson("http://appcore.vdimdci.com.ph/api/user/88", hdr, nil)
	if !exapi.OK() {
		t.Fail()
	}
	fmt.Printf("%v", string(exapi.Data))
}

func TestExecuteAPIGET2(t *testing.T) {
	hdr := make(map[string]string)
	//hdr["Cookie"] = "APPSHUB-WF-login=zaldy.baguinon; APPSHUB-WF-appdomain=MDCI"
	//hdr["Authorization"] = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJLZWFudS1VcGxvYWRlciIsImV4cCI6MCwibmJmIjoxNjc1MzE0MTg0LCJpYXQiOjAsInVzciI6ImphbWVzLmx1bWliYW9AbWRjaS5jb20ucGgiLCJkb20iOiJNRENJIiwiYXBwIjoiS2VhbnUtVXBsb2FkZXIiLCJkZXYiOiIyS1JzS3Z4Y2NuOUp0RjNxbDIxMmN1MmhwS1MifQ.961xUrBObQfN6fkO_s7OYhFTqKC_aMrr1OKVwvPhkLU"

	exapi := GetJson("http://inform.vdimdci.com.ph/api/email/?num=50", hdr, nil)
	//exapi := GetJSON("http://localhost:15001/email/?num=50", hdr)
	if !exapi.OK() {
		t.Fail()
	}
	fmt.Printf("%v", string(exapi.Data))
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

func TestBuildAccessToken(t *testing.T) {
	jwth := map[string]interface{}{
		"alg": "HS256",
		"typ": "JWT",
	}

	jwtc := map[string]interface{}{
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		//"aud": []string{"APPSHUB-AUTH"},
		"aud": "APPSHUB-AUTH",
		"usr": "zaldy.baguinon",
		"dom": "MDCI",
		"dev": "j4h2j34h23jk4h3kj4hfdsfsdf",
		"app": "APPSHUB-AITH",
	}

	token := BuildAccessToken(&jwth, &jwtc, "thisisanhmacsecretkey")
	if token == "" {
		t.Fail()
	}
	fmt.Println(token)
}

func TestParseAccessToken(t *testing.T) {
	var pl CustomPayload

	HMAC := jwt.NewHS256([]byte("thisisanhmacsecretkey"))

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

	jwtfromck := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJBUFBTSFVCLUFVVEgiLCJleHAiOjAsIm5iZiI6MTU3OTIyMTM1MywiaWF0IjoxNTc5MjIxMzUzLCJ1c3IiOiJ6YWxkeS5iYWd1aW5vbiIsImRvbSI6Ik1EQ0kiLCJhcHAiOiJBUFBTSFVCLUFJVEgiLCJkZXYiOiJqNGgyajM0aDIzams0aDNrajRoZmRzZnNkZiJ9.MS77eSy7rg0a8-wTyaGmSbR8kOtZCv0092qVoucpG9k"

	if _, err := jwt.Verify([]byte(jwtfromck), HMAC, &pl); err == nil {

		fmt.Printf("%+v", pl)
	} else {
		t.Fail()
	}
}

func TestParseJWT(t *testing.T) {
	jwtfromck := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJBUFBTSFVCLUFVVEgiLCJleHAiOjAsIm5iZiI6MTU3OTIyMTM1MywiaWF0IjoxNTc5MjIxMzUzLCJ1c3IiOiJ6YWxkeS5iYWd1aW5vbiIsImRvbSI6Ik1EQ0kiLCJhcHAiOiJBUFBTSFVCLUFJVEgiLCJkZXYiOiJqNGgyajM0aDIzams0aDNrajRoZmRzZnNkZiJ9.MS77eSy7rg0a8-wTyaGmSbR8kOtZCv0092qVoucpG9k"

	vj, err := ParseJWT(jwtfromck, "thisisanhmacsecretkey", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(vj)
}

func BenchmarkParseJWT(b *testing.B) {
	jwtfromck := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJBUFBTSFVCLUFVVEgiLCJleHAiOjAsIm5iZiI6MTU3OTIyMTM1MywiaWF0IjoxNTc5MjIxMzUzLCJ1c3IiOiJ6YWxkeS5iYWd1aW5vbiIsImRvbSI6Ik1EQ0kiLCJhcHAiOiJBUFBTSFVCLUFJVEgiLCJkZXYiOiJqNGgyajM0aDIzams0aDNrajRoZmRzZnNkZiJ9.MS77eSy7rg0a8-wTyaGmSbR8kOtZCv0092qVoucpG9k"

	vj, err := ParseJWT(jwtfromck, "thisisanhmacsecretkey", false)
	if err != nil {
		b.Fatal(err)
	}
	_ = vj
	//b.Log(vj)
}

func TestResultAny(t *testing.T) {
	res := ResultAny[string]{}
	res.AddErr(errors.New("overriden result"))
	t.Log(res)
}

func TestGetAny(t *testing.T) {
	hdr := make(map[string]string)
	hdr["Cookie"] = "APPSHUB-WF-login=zaldy.baguinon; APPSHUB-WF-appdomain=MDCI"
	hdr["Authorization"] = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJLZWFudS1VcGxvYWRlciIsImV4cCI6MCwibmJmIjoxNjc1MzE0MTg0LCJpYXQiOjAsInVzciI6ImphbWVzLmx1bWliYW9AbWRjaS5jb20ucGgiLCJkb20iOiJNRENJIiwiYXBwIjoiS2VhbnUtVXBsb2FkZXIiLCJkZXYiOiIyS1JzS3Z4Y2NuOUp0RjNxbDIxMmN1MmhwS1MifQ.961xUrBObQfN6fkO_s7OYhFTqKC_aMrr1OKVwvPhkLU"
	type UserAccount struct {
		Key                 *int       `json:"user_key,omitempty"`
		UserName            *string    `json:"user_name,omitempty"` // UserName, required
		Password            *string    `json:"password,omitempty"`  // Password, required
		Active              *bool      `json:"active,omitempty"`
		AppsHubAdmin        *bool      `json:"apps_hub_admin,omitempty"`
		ApplicationID       *string    `json:"application_id,omitempty"` // ApplicationID, for specific uses only, not on the database
		DisplayName         *string    `json:"display_name,omitempty"`
		ProfileImageURI     *string    `json:"profile_image_uri,omitempty"`
		ProfileImageExt     *string    `json:"profile_image_ext,omitempty"`
		EmailAddress        *string    `json:"email_address,omitempty"` // EmailAddress, required
		MobileNo            *string    `json:"mobile_no,omitempty"`     // MobileNo, required if email address is not available
		LDAPLogin           *bool      `json:"ldap_login,omitempty"`
		TransactionID       *string    `json:"transaction_id,omitempty"`
		ActivationCode      *string    `json:"activation_code,omitempty"`
		ActivationStatus    *string    `json:"activation_status,omitempty"`
		RegistrationChannel *string    `json:"registration_channel,omitempty"`
		DateLastLoggedIn    *time.Time `json:"date_last_logged_in,omitempty"`
	}
	exapi := ReadApi[[]UserAccount](
		"https://appcore.vdimdci.com.ph/api/user/19",
		TimeOut(1000),
		Headers(hdr, &sync.RWMutex{}),
		Compressed(true))
	if !exapi.OK() {
		t.Fail()
	}
	fmt.Printf("%v", exapi.Data)

	//ua := UserAccount{}
	//CreateApi[[]UserAccount]("https://appcore.vdimdci.com.ph/api/user/19", ua, true, hdr, &sync.RWMutex{})

}
