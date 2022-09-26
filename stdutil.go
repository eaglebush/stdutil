package stdutil

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	ssd "github.com/shopspring/decimal"
	"golang.org/x/exp/constraints"
)

type FieldTypeConstraint interface {
	constraints.Ordered | time.Time | ssd.Decimal | bool | byte
}

// AnyToString - convert any variable to string
func AnyToString(value interface{}) string {
	var b string

	if value == nil {
		return ""
	}

	switch t := value.(type) {
	case string:
		b = t
	case int:
		b = strconv.FormatInt(int64(t), 10)
	case int8:
		b = strconv.FormatInt(int64(t), 10)
	case int16:
		b = strconv.FormatInt(int64(t), 10)
	case int32:
		b = strconv.FormatInt(int64(t), 10)
	case int64:
		b = strconv.FormatInt(t, 10)
	case uint:
		b = strconv.FormatUint(uint64(t), 10)
	case uint8:
		b = strconv.FormatUint(uint64(t), 10)
	case uint16:
		b = strconv.FormatUint(uint64(t), 10)
	case uint32:
		b = strconv.FormatUint(uint64(t), 10)
	case uint64:
		b = strconv.FormatUint(uint64(t), 10)
	case float32:
		b = fmt.Sprintf("%f", t)
	case float64:
		b = fmt.Sprintf("%f", t)
	case bool:
		if t {
			return "true"
		} else {
			return "false"
		}
	case time.Time:
		b = "'" + t.Format(time.RFC3339) + "'"
	case *string:
		if t == nil {
			return ""
		}

		b = *t
	case *int:
		if t == nil {
			return "0"
		}

		b = strconv.FormatInt(int64(*t), 10)
	case *int8:
		if t == nil {
			return "0"
		}

		b = strconv.FormatInt(int64(*t), 10)
	case *int16:
		if t == nil {
			return "0"
		}

		b = strconv.FormatInt(int64(*t), 10)
	case *int32:
		if t == nil {
			return "0"
		}

		b = strconv.FormatInt(int64(*t), 10)
	case *int64:
		if t == nil {
			return "0"
		}

		b = strconv.FormatInt(*t, 10)
	case *uint:
		if t == nil {
			return "0"
		}
		b = strconv.FormatUint(uint64(*t), 10)
	case *uint8:
		if t == nil {
			return "0"
		}

		b = strconv.FormatUint(uint64(*t), 10)
	case *uint16:
		if t == nil {
			return "0"
		}

		b = strconv.FormatUint(uint64(*t), 10)
	case *uint32:
		if t == nil {
			return "0"
		}

		b = strconv.FormatUint(uint64(*t), 10)
	case *uint64:
		if t == nil {
			return "0"
		}

		b = strconv.FormatUint(uint64(*t), 10)
	case *float32:
		if t == nil {
			return "0"
		}

		b = fmt.Sprintf("%f", *t)
	case *float64:
		if t == nil {
			return "0"
		}

		b = fmt.Sprintf("%f", *t)
	case *bool:
		if t == nil {
			return "false"
		}

		if *t {
			return "true"
		} else {
			return "false"
		}
	case *time.Time:
		if t == nil {
			return "'" + time.Time{}.Format(time.RFC3339) + "'"
		}

		tm := *t
		b = "'" + tm.Format(time.RFC3339) + "'"
	}

	return b
}

// Itos is a shortcut to AnyToString. I means Interface
func Itos(value interface{}) string {
	return AnyToString(value)
}

// IntToInterfaceArray - converts a name value array to interface array
func IntToInterfaceArray(values int) []interface{} {
	args := make([]interface{}, 1)
	args[0] = values
	return args
}

// IsNumeric - checks if a string is numeric
func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// IsNullOrEmpty checks for nullity and emptiness of a pointer variable
// Currently supported data types are the ones in the constraints.Ordered,
// time.Time, bool and shopspring/decimal
//
// This function requires version 1.18+
func IsNullOrEmpty[T FieldTypeConstraint](value *T) bool {
	return value == nil || *value == GetZero[T]()
}

// IsNullOrEmpty checks for emptiness of a pointer variable ignoring nullity
// Currently supported data types are the ones in the constraints.Ordered,
// time.Time, bool and shopspring/decimal
//
// This function requires version 1.18+
func IsEmpty[T FieldTypeConstraint](value *T) bool {
	return value != nil && *value == GetZero[T]()
}

// Val gets the value of a pointer in order
// Currently supported data types are the ones in the constraints.Ordered,
// time.Time, bool and shopspring/decimal
//
// This function requires version 1.18+
func Val[T FieldTypeConstraint](value *T) T {
	if value == nil {
		return GetZero[T]()
	}

	return *value
}

// New initializes a variable and returns a pointer of its type
// Currently supported data types are the ones in the constraints.Ordered,
// time.Time, bool and shopspring/decimal
//
// This function requires version 1.18+
func New[T FieldTypeConstraint](value T) *T {
	n := new(T)
	*n = value
	return n
}

// NameValuesToInterfaceArray converts name values to interface array
func NameValuesToInterfaceArray(values NameValues) []interface{} {

	args := make([]interface{}, len(values.Pair))
	i := 0

	for _, v := range values.Pair {
		args[i] = v
		i++
	}

	return args
}

// InterpolateString interpolates string with the name value pairs
func InterpolateString(base string, keyValues NameValues) (string, []interface{}) {

	retstr := base
	hasmatch := false
	pattern := `\$\{(\w*)\}` //search for ${*}
	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(base, -1)

	retif := make([]interface{}, len(matches))

	for i, match := range matches {
		hasmatch = false
		for n, v := range keyValues.Pair {
			if strings.EqualFold(match, `${`+n+`}`) {
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

// In checks if the seek parameter is in the list parameter
func In[T comparable](seek T, list ...T) bool {
	for _, li := range list {
		if li == seek {
			return true
		}
	}
	return false
}

// SortByKeyArray - reorder keys and values based on a keyOrder array sequence
func SortByKeyArray(values *NameValues, keyOrder *[]string) NameValues {
	ret := NameValues{}
	ret.Pair = make(map[string]any)

	//If keyorder was specified, the order of keys will be sorted according to the specifications
	ko := *keyOrder
	if len(ko) == 0 {
		return ret
	}

	for i := 0; i < len(ko); i++ {
		for k, v := range values.Pair {
			if strings.EqualFold(ko[i], k) {
				ret.Pair[k] = v
				break
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

// GetZero gets the zero value of the types defined as
// constraints.Ordered, time.Time and shopspring/decimal
//
// This function requires version 1.18+
func GetZero[T FieldTypeConstraint]() T {
	var result T
	return result
}
