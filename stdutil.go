package stdutil

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// AnyToString - convert any variable to string
func AnyToString(value interface{}) string {
	var b string

	switch value.(type) {
	case string:
		b = value.(string)
	case int:
		b = strconv.FormatInt(int64(value.(int)), 10)
	case int8:
		b = strconv.FormatInt(int64(value.(int8)), 10)
	case int16:
		b = strconv.FormatInt(int64(value.(int16)), 10)
	case int32:
		b = strconv.FormatInt(int64(value.(int32)), 10)
	case int64:
		b = strconv.FormatInt(value.(int64), 10)
	case uint:
		b = strconv.FormatUint(uint64(value.(uint)), 10)
	case uint8:
		b = strconv.FormatUint(uint64(value.(uint8)), 10)
	case uint16:
		b = strconv.FormatUint(uint64(value.(uint16)), 10)
	case uint32:
		b = strconv.FormatUint(uint64(value.(uint32)), 10)
	case uint64:
		b = strconv.FormatUint(uint64(value.(uint64)), 10)
	case float32:
		b = fmt.Sprintf("%f", value.(float32))
	case float64:
		b = fmt.Sprintf("%f", value.(float64))
	case bool:
		b = "false"
		s := strings.ToLower(value.(string))
		if len(s) > 0 {
			if s == "true" || s == "on" || s == "yes" || s == "1" || s == "-1" {
				b = "true"
			}
		}
	case time.Time:
		b = "'" + value.(time.Time).Format(time.RFC3339) + "'"
	}

	return b
}

// IntToInterfaceArray - converts a name value array to interface array
func IntToInterfaceArray(values int) []interface{} {
	args := make([]interface{}, 1)
	args[0] = values
	return args
}

//IsNumeric - checks if a string is numeric
func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// NameValuesToInterfaceArray - converts name values to interface array
func NameValuesToInterfaceArray(values NameValues) []interface{} {

	args := make([]interface{}, len(values.Pair))
	i := 0
	for _, v := range values.Pair {
		args[i] = v.Value
		i++
	}

	return args
}

// NameValuesToValidationExpressionArray - converts name values to ValidationExpression array
func NameValuesToValidationExpressionArray(values NameValues) []ValidationExpression {

	args := make([]ValidationExpression, len(values.Pair))
	i := 0
	for _, v := range values.Pair {
		args[i].Name = AnyToString(v.Name)
		args[i].Value = AnyToString(v.Value)
		args[i].Operator = `=`
		i++
	}

	return args
}

// InterpolateString - Interpolate string with the name value pairs
func InterpolateString(base string, keyValues NameValues) (string, []interface{}) {
	retstr := base

	hasmatch := false

	pattern := `\$\{(\w*)\}` //search for ${*}
	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(base, -1)

	retif := make([]interface{}, len(matches))

	for i, match := range matches {
		hasmatch = false
		for _, vs := range keyValues.Pair {
			n := strings.ToLower(vs.Name)
			v := vs.Value
			if match == `${`+n+`}` {
				retstr = strings.Replace(retstr, match, AnyToString(v), -1)
				retif[i] = v
				hasmatch = true
				break
			}
		}

		/* The matches needs to have a default value */
		if !hasmatch {
			retstr = strings.Replace(retstr, match, "0", -1)
			retif[i] = "0" //a string 0 would cater to both string and number columns
		}
	}

	return retstr, retif
}

// ValidateEmail - validate an e-mail address
func ValidateEmail(email string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return re.MatchString(email)
}

// SortByKeyArray - reorder keys and values based on a keyOrder array sequence
func SortByKeyArray(values *NameValues, keyOrder *[]string) NameValues {
	ret := NameValues{}
	ret.Pair = make([]NameValue, 0)

	//If keyorder was specified, the order of keys will be sorted according to the specifications
	ko := *keyOrder
	if len(ko) > 0 {
		for i := 0; i < len(ko); i++ {
			for _, v := range values.Pair {
				kv := v.Name
				if strings.ToLower(ko[i]) == strings.ToLower(kv) {
					ret.Pair = append(ret.Pair, v)
					break
				}
			}
		}
	}
	return ret
}

// StripEndingForwardSlash - remove the ending forward slash of a string
func StripEndingForwardSlash(value string) string {
	v := strings.TrimSpace(value)
	v = strings.ReplaceAll(v, `\`, `/`)
	ix := strings.LastIndex(v, `/`)
	if ix == (len(v) - 1) {
		v = v[0:ix]
	}
	return v
}

// StripTrailing - strip string of trailing characters after the length
func StripTrailing(value string, length int) string {

	if len(value) > length {
		return value[0:length]
	}

	return value
}

// StripLeading - strip string of leading characters by an offset
func StripLeading(value string, offset int) string {

	if len(value) > offset {
		return value[offset:]
	}

	return value
}

// NewString initializes a string pointer with an initial value
func NewString(initial string) (init *string) {
	init = new(string)
	*init = initial
	return
}

// NewByte initializes a byte pointer with an initial value
func NewByte(initial byte) (init *byte) {
	init = new(byte)
	*init = initial
	return
}

// NewInt initializes an int pointer with an initial value
func NewInt(initial int) (init *int) {
	init = new(int)
	*init = initial
	return
}

// NewInt32 initializes an int32 pointer with an initial value
func NewInt32(initial int32) (init *int32) {
	init = new(int32)
	*init = initial
	return
}

// NewInt64 initializes an int64 pointer with an initial value
func NewInt64(initial int64) (init *int64) {
	init = new(int64)
	*init = initial
	return
}

// NewBool initializes a bool pointer with an initial value
func NewBool(initial bool) (init *bool) {
	init = new(bool)
	*init = initial
	return
}

// NewFloat32 initializes a float32 pointer with an initial value
func NewFloat32(initial float32) (init *float32) {
	init = new(float32)
	*init = initial
	return
}

// NewFloat64 initializes a float64 pointer with an initial value
func NewFloat64(initial float64) (init *float64) {
	init = new(float64)
	*init = initial
	return
}

// NewTime initializes a time.Time pointer with an initial value
func NewTime(initial *time.Time) (init *time.Time) {
	init = new(time.Time)

	if initial != nil {
		init = initial
	}

	return
}
