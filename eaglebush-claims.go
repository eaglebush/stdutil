package stdutil

import (
	"crypto/subtle"
	"fmt"
	"time"
)

// // For a type to be a Claims object, it must just have a Valid method that determines
// // if the token is invalid for any supported reason
// type Claims interface {
// 	Valid() error
// }

// EaglebushClaims - eaglebush claims structure
type EaglebushClaims struct {
	Audience      string `json:"aud,omitempty"`
	ExpiresAt     int64  `json:"exp,omitempty"`
	ID            string `json:"jti,omitempty"`
	IssuedAt      int64  `json:"iat,omitempty"`
	Issuer        string `json:"iss,omitempty"`
	NotBefore     int64  `json:"nbf,omitempty"`
	Subject       string `json:"sub,omitempty"`
	UserName      string `json:"usr,omitempty"`
	Domain        string `json:"dom,omitempty"`
	DeviceID      string `json:"dev,omitempty"`
	ApplicationID string `json:"app,omitempty"`
}

// TimeFunc provides the current time when parsing token to validate "exp" claim (expiration time).
// You can override it to use another time value.  This is useful for testing or if your
// server uses a different time zone than your tokens.
var TimeFunc = time.Now

// Valid - Validates time based claims "exp, iat, nbf".
// There is no accounting for clock skew.
// As well, if any of the above claims are not in the token, it will still
// be considered a valid claim.
func (c EaglebushClaims) Valid() error {

	vErr := new(ValidationError)
	now := TimeFunc().Unix()

	// Required claims
	if c.VerifyUserName("", true) == false {
		vErr.Inner = fmt.Errorf("Token has no user name")
		vErr.Errors |= ValidationErrorUserName
	}

	if c.VerifyDomain("", true) == false {
		vErr.Inner = fmt.Errorf("Token has no domain")
		vErr.Errors |= ValidationErrorDomain
	}

	if c.VerifyDeviceID("", true) == false {
		vErr.Inner = fmt.Errorf("Token has no device id")
		vErr.Errors |= ValidationErrorDeviceID
	}

	if c.VerifyApplicationID("", true) == false {
		vErr.Inner = fmt.Errorf("Token has no application id")
		vErr.Errors |= ValidationErrorDeviceID
	}

	// The claims below are optional, by default, so if they are set to the
	// default value in Go, let's not fail the verification for them.
	if c.VerifyExpiresAt(now, false) == false {
		delta := time.Unix(now, 0).Sub(time.Unix(c.ExpiresAt, 0))
		vErr.Inner = fmt.Errorf("token is expired by %v", delta)
		vErr.Errors |= ValidationErrorExpired
	}

	if c.VerifyIssuedAt(now, false) == false {
		vErr.Inner = fmt.Errorf("Token used before issued")
		vErr.Errors |= ValidationErrorIssuedAt
	}

	if c.VerifyNotBefore(now, false) == false {
		vErr.Inner = fmt.Errorf("token is not valid yet")
		vErr.Errors |= ValidationErrorNotValidYet
	}

	if vErr.valid() {
		return nil
	}

	return vErr
}

// VerifyAudience - Compares the aud claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *EaglebushClaims) VerifyAudience(cmp string, req bool) bool {
	return verifyAud(c.Audience, cmp, req)
}

// VerifyExpiresAt - Compares the exp claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *EaglebushClaims) VerifyExpiresAt(cmp int64, req bool) bool {
	return verifyExp(c.ExpiresAt, cmp, req)
}

// VerifyIssuedAt - Compares the iat claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *EaglebushClaims) VerifyIssuedAt(cmp int64, req bool) bool {
	return verifyIat(c.IssuedAt, cmp, req)
}

// VerifyIssuer - Compares the iss claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *EaglebushClaims) VerifyIssuer(cmp string, req bool) bool {
	return verifyIss(c.Issuer, cmp, req)
}

// VerifyNotBefore - Compares the nbf claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *EaglebushClaims) VerifyNotBefore(cmp int64, req bool) bool {
	return verifyNbf(c.NotBefore, cmp, req)
}

// VerifyUserName - Compares the usr claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *EaglebushClaims) VerifyUserName(cmp string, req bool) bool {
	return verifyUsr(c.UserName, cmp, req)
}

// VerifyDomain - Compares the dom claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *EaglebushClaims) VerifyDomain(cmp string, req bool) bool {
	return verifyDom(c.Domain, cmp, req)
}

// VerifyDeviceID - Compares the dom claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *EaglebushClaims) VerifyDeviceID(cmp string, req bool) bool {
	return verifyDev(c.DeviceID, cmp, req)
}

// VerifyApplicationID - Compares the app claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *EaglebushClaims) VerifyApplicationID(cmp string, req bool) bool {
	return verifyDev(c.ApplicationID, cmp, req)
}

// ----- helpers

func verifyAud(aud string, cmp string, required bool) bool {
	if aud == "" {
		return !required
	}
	if subtle.ConstantTimeCompare([]byte(aud), []byte(cmp)) != 0 {
		return true
	}

	return false
}

func verifyExp(exp int64, now int64, required bool) bool {
	if exp == 0 {
		return !required
	}
	return now <= exp
}

func verifyIat(iat int64, now int64, required bool) bool {
	if iat == 0 {
		return !required
	}
	return now >= iat
}

func verifyIss(iss string, cmp string, required bool) bool {
	if iss == "" {
		return !required
	}
	if subtle.ConstantTimeCompare([]byte(iss), []byte(cmp)) != 0 {
		return true
	}

	return false
}

func verifyNbf(nbf int64, now int64, required bool) bool {
	if nbf == 0 {
		return !required
	}
	return now >= nbf
}

func verifyUsr(usr string, cmp string, required bool) bool {
	if usr == "" {
		return !required
	}
	if subtle.ConstantTimeCompare([]byte(usr), []byte(cmp)) != 0 {
		return true
	}

	return false
}

func verifyDom(dom string, cmp string, required bool) bool {
	if dom == "" {
		return !required
	}
	if subtle.ConstantTimeCompare([]byte(dom), []byte(cmp)) != 0 {
		return true
	}

	return false
}

func verifyDev(dev string, cmp string, required bool) bool {
	if dev == "" {
		return !required
	}
	if subtle.ConstantTimeCompare([]byte(dev), []byte(cmp)) != 0 {
		return true
	}

	return false
}

func verifyApp(app string, cmp string, required bool) bool {
	if app == "" {
		return !required
	}
	if subtle.ConstantTimeCompare([]byte(app), []byte(cmp)) != 0 {
		return true
	}

	return false
}
