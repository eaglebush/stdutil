package stdutil

import (
	"fmt"
	"testing"
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
