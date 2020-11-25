package stdutil

import (
	"strconv"
	"strings"
)

// ArgumentValue - argument value
type ArgumentValue struct {
	Value string
}

// Bool - Return boolean from value
func (av *ArgumentValue) Bool() bool {

	s := strings.ToLower(av.Value)

	if s == "true" || s == "on" || s == "yes" || s == "1" || s == "-1" {
		return true
	}

	return false
}

// Int64 - Return int64 from value
func (av *ArgumentValue) Int64() int64 {
	v, _ := strconv.ParseInt(av.Value, 10, 64)
	return v
}

// Float64 - Return float64 from value
func (av *ArgumentValue) Float64() float64 {
	v, _ := strconv.ParseFloat(av.Value, 64)
	return v
}

// String - return string from value
func (av *ArgumentValue) String() string {
	return av.Value
}
