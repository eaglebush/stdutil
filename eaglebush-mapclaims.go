package stdutil

import (
	"encoding/json"
	"errors"
	"fmt"
	// "fmt"
)

// EaglebushMapClaims - Claims type that uses the map[string]interface{} for JSON decoding
// This is the default claims type if you don't supply one
type EaglebushMapClaims map[string]interface{}

// VerifyAudience - Compares the aud claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (m EaglebushMapClaims) VerifyAudience(cmp string, req bool) bool {
	aud, _ := m["aud"].(string)
	return verifyAud(aud, cmp, req)
}

// VerifyExpiresAt - Compares the exp claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (m EaglebushMapClaims) VerifyExpiresAt(cmp int64, req bool) bool {
	switch exp := m["exp"].(type) {
	case float64:
		return verifyExp(int64(exp), cmp, req)
	case json.Number:
		v, _ := exp.Int64()
		return verifyExp(v, cmp, req)
	}
	return req == false
}

// VerifyIssuedAt - the iat claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (m EaglebushMapClaims) VerifyIssuedAt(cmp int64, req bool) bool {
	switch iat := m["iat"].(type) {
	case float64:
		return verifyIat(int64(iat), cmp, req)
	case json.Number:
		v, _ := iat.Int64()
		return verifyIat(v, cmp, req)
	}
	return req == false
}

// VerifyIssuer - Compares the iss claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (m EaglebushMapClaims) VerifyIssuer(cmp string, req bool) bool {
	iss, _ := m["iss"].(string)
	return verifyIss(iss, cmp, req)
}

// VerifyNotBefore - Compares the nbf claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (m EaglebushMapClaims) VerifyNotBefore(cmp int64, req bool) bool {
	switch nbf := m["nbf"].(type) {
	case float64:
		return verifyNbf(int64(nbf), cmp, req)
	case json.Number:
		v, _ := nbf.Int64()
		return verifyNbf(v, cmp, req)
	}
	return req == false
}

// VerifyUserName - Compares the usr claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (m EaglebushMapClaims) VerifyUserName(cmp string, req bool) bool {
	usr, _ := m["usr"].(string)
	return verifyUsr(usr, cmp, req)
}

// VerifyDomain - Compares the dom claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (m EaglebushMapClaims) VerifyDomain(cmp string, req bool) bool {
	dom, _ := m["dom"].(string)
	return verifyDom(dom, cmp, req)
}

// VerifyDeviceID - Compares the dom claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (m EaglebushMapClaims) VerifyDeviceID(cmp string, req bool) bool {
	dev, _ := m["dev"].(string)
	return verifyDev(dev, cmp, req)
}

// VerifyApplicationID - Compares app dom claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (m EaglebushMapClaims) VerifyApplicationID(cmp string, req bool) bool {
	app, _ := m["app"].(string)
	return verifyApp(app, cmp, req)
}

// Valid - Validates time based claims "exp, iat, nbf".
// There is no accounting for clock skew.
// As well, if any of the above claims are not in the token, it will still
// be considered a valid claim.
func (m EaglebushMapClaims) Valid() error {
	vErr := new(ValidationError)
	now := TimeFunc().Unix()

	// Required claims
	if m.VerifyUserName("", true) == false {
		vErr.Inner = fmt.Errorf("Token has no user name")
		vErr.Errors |= ValidationErrorUserName
	}

	if m.VerifyDomain("", true) == false {
		vErr.Inner = fmt.Errorf("Token has no domain")
		vErr.Errors |= ValidationErrorDomain
	}

	if m.VerifyDeviceID("", true) == false {
		vErr.Inner = fmt.Errorf("Token has no device id")
		vErr.Errors |= ValidationErrorDeviceID
	}

	if m.VerifyApplicationID("", true) == false {
		vErr.Inner = fmt.Errorf("Token has no application id")
		vErr.Errors |= ValidationErrorDeviceID
	}

	if m.VerifyExpiresAt(now, false) == false {
		vErr.Inner = errors.New("Token is expired")
		vErr.Errors |= ValidationErrorExpired
	}

	if m.VerifyIssuedAt(now, false) == false {
		vErr.Inner = errors.New("Token used before issued")
		vErr.Errors |= ValidationErrorIssuedAt
	}

	if m.VerifyNotBefore(now, false) == false {
		vErr.Inner = errors.New("Token is not valid yet")
		vErr.Errors |= ValidationErrorNotValidYet
	}

	if vErr.valid() {
		return nil
	}

	return vErr
}
