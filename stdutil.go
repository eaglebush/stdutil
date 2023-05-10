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

type NumericConstraint interface {
	constraints.Integer | constraints.Float
}

type StringValidationOptions struct {
	Empty    bool // Allow empty string. Default: false, will raise an error if the string is empty
	Null     bool // Allow null. Default: false, will raise an error if the string is null
	Min      int  // Minimum length. Default: 0
	Max      int  // Maximum length. Default: 0
	NoSpaces bool // Do not allow spaces in the string. Default: false. Setting to true will raise an error if the string has spaces
	Extended []func(value *string) error
}

type TimeValidationOptions struct {
	Null     bool       // Allow null. Default: false, will raise an error if the time is null
	Empty    bool       // Allow zero time Default: false, will raise an error if the time is zero
	Min      *time.Time // Minimum time. Default: nil
	Max      *time.Time // Maximum time. Default: nil
	DateOnly bool       // Compare dates only. Default: false
	Extended []func(value *time.Time) error
}

type NumericValidationOptions[T NumericConstraint] struct {
	Null     bool // Allow null. Default: false, will raise an error if the time is null
	Empty    bool // Allow zero time Default: false, will raise an error if the time is zero
	Min      T    // Minimum time. Default: nil
	Max      T    // Maximum time. Default: nil
	Extended []func(value *T) error
}

type DecimalValidationOptions struct {
	Null     bool         // Allow null. Default: false, will raise an error if the decimal is null
	Empty    bool         // Allow zero decimal. Default: false, will raise an error if the decimal is zero
	Min      *ssd.Decimal // Minimum decimal value. Default: nil
	Max      *ssd.Decimal // Maximum decimal value. Default: nil
	Extended []func(value *ssd.Decimal) error
}

type SeriesOptions struct {
	Prefix string // Prefix of series
	Suffix string // Suffix of series
	Length int    // Fixed length of the series
}

// AnyToString converts any variable to string
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

// ToInterfaceArray converts a value to interface array
//
// This function requires version 1.18+
func ToInterfaceArray[T FieldTypeConstraint](values T) []interface{} {
	return []interface{}{values}
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

// Interpolate interpolates string with the name value pairs
func Interpolate(base string, keyValues NameValues) (string, []interface{}) {

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
func ValidateEmail(email *string) error {
	if email == nil || *email == "" {
		return fmt.Errorf("is an invalid email address")
	}

	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !re.MatchString(*email) {
		return fmt.Errorf("is an invalid email address")
	}
	return nil
}

// IsStringNumeric checks if a string is numeric
func IsStringNumeric(value *string) error {
	if value == nil {
		return fmt.Errorf("is empty")
	}

	if _, err := strconv.ParseFloat(*value, 64); err != nil {
		return fmt.Errorf(`is not a number (%s)`, err)
	}

	return nil
}

// ValidateString validates an input string against the string validation options
func ValidateString(value *string, opts *StringValidationOptions) error {

	// If options were not set, this string is valid
	// If value is nil and the Null option is false, we raise an error
	// If value is empty and the Empty option is false, we raise an error
	if opts == nil {
		return nil
	}

	if value == nil {
		if !opts.Null {
			return fmt.Errorf("must be provided (nil)")
		}
		return nil
	}

	ln := len([]rune(*value))

	if ln == 0 {
		if !opts.Empty {
			return fmt.Errorf("must be provided (empty)")
		}
		return nil
	}

	if opts.Min > 0 && ln < opts.Min {
		return fmt.Errorf("is shorter than %d characters", opts.Min)
	}

	if opts.Max > 0 && ln > opts.Max {
		return fmt.Errorf("is longer than %d characters", opts.Max)
	}

	if opts.NoSpaces && strings.Contains(*value, " ") {
		return fmt.Errorf("contains spaces")
	}

	for _, f := range opts.Extended {
		if err := f(value); err != nil {
			return err
		}
	}

	return nil
}

// ValidateTime validates an input time against the time validation options
func ValidateTime(value *time.Time, opts *TimeValidationOptions) error {

	// If options were not set, this time is valid
	// If value is nil and the Null option is false, we raise an error
	if opts == nil {
		return nil
	}

	if value == nil {
		if !opts.Null {
			return fmt.Errorf("must be provided (nil)")
		}
		return nil
	}

	if value.IsZero() {
		if !opts.Empty {
			return fmt.Errorf("must be provided (empty)")
		}
		return nil
	}

	if opts.DateOnly {
		dv := *value
		*value = time.Date(dv.Year(), dv.Month(), dv.Day(), 0, 0, 0, 0, dv.Location())

		if opts.Min != nil {
			dc := opts.Min
			*opts.Min = time.Date(dc.Year(), dc.Month(), dc.Day(), 0, 0, 0, 0, dc.Location())
		}

		if opts.Max != nil {
			dc := opts.Max
			*opts.Max = time.Date(dc.Year(), dc.Month(), dc.Day(), 0, 0, 0, 0, dc.Location())
		}
	}

	if opts.Min != nil && value.Before(*opts.Min) {
		return fmt.Errorf("is earlier than %s minimum time", opts.Min)
	}

	if opts.Max != nil && value.After(*opts.Max) {
		return fmt.Errorf("is later than %s maximum time", opts.Max)
	}

	for _, f := range opts.Extended {
		if err := f(value); err != nil {
			return err
		}
	}

	return nil
}

// ValidateNumeric validates a numeric input against numeric validation options
func ValidateNumeric[T NumericConstraint](value *T, opts *NumericValidationOptions[T]) error {

	// If options were not set, this time is valid
	// If value is nil and the Null option is false, we raise an error
	if opts == nil {
		return nil
	}

	if value == nil {
		if !opts.Null {
			return fmt.Errorf("must be provided (nil)")
		}
		return nil
	}

	if *value == 0 {
		if !opts.Empty {
			return fmt.Errorf("must be provided (empty)")
		}
	}

	if opts.Min > 0 && *value < opts.Min {
		return fmt.Errorf("is lesser than %v minimum value", opts.Min)
	}

	if opts.Max > 0 && *value > opts.Max {
		return fmt.Errorf("is greater than %v maximum value", opts.Max)
	}

	for _, f := range opts.Extended {
		if err := f(value); err != nil {
			return err
		}
	}

	return nil
}

// ValidateDecimal validates a decimal input against decimal validation options
func ValidateDecimal(value *ssd.Decimal, opts *DecimalValidationOptions) error {

	// If options were not set, this decimal is valid
	// If value is nil and the Null option is false, we raise an error
	if opts == nil {
		return nil
	}

	if value == nil {
		if !opts.Null {
			return fmt.Errorf("must be provided (nil)")
		}
		return nil
	}

	if value.IsZero() {
		if !opts.Empty {
			return fmt.Errorf("must be provided (empty)")
		}
	}

	zero := ssd.NewFromInt(0)

	if opts.Min != nil && opts.Min.GreaterThan(zero) && value.LessThan(*opts.Min) {
		return fmt.Errorf("is lesser than %v minimum value", *opts.Min)
	}

	if opts.Max != nil && opts.Max.GreaterThan(zero) && value.GreaterThan(*opts.Max) {
		return fmt.Errorf("is greater than %v maximum value", *opts.Max)
	}

	for _, f := range opts.Extended {
		if err := f(value); err != nil {
			return err
		}
	}

	return nil
}

// BuildSeries builds series based on options
func BuildSeries(series int, opt SeriesOptions) string {

	// If length is specified, we get the difference between suffix and prefix
	if opt.Length > 0 {
		diff := opt.Length - (len(opt.Prefix) + len(opt.Suffix))
		ds := `%0` + strconv.Itoa(diff) + `d`
		return fmt.Sprintf(`%s`+ds+`%s`, opt.Prefix, series, opt.Suffix)
	}

	return fmt.Sprintf(`%s%d%s`, opt.Prefix, series, opt.Suffix)
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

// SortByKey reorders keys and values based on a keyOrder array sequence
func SortByKey(values *NameValues, keyOrder *[]string) NameValues {

	if keyOrder == nil {
		return *values
	}

	ko := *keyOrder
	if len(ko) == 0 {
		return *values
	}

	ret := NameValues{}
	ret.Pair = make(map[string]any)

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

// StripEndingForwardSlash removes the ending forward slash of a string
func StripEndingForwardSlash(value string) string {
	v := strings.TrimSpace(value)
	v = strings.ReplaceAll(v, `\`, `/`)
	ix := strings.LastIndex(v, `/`)

	if ix == (len(v) - 1) {
		return v[0:ix]
	}

	return v
}

// StripTrailing strips string of trailing characters after the length
func StripTrailing(value string, length int) string {

	str := []rune(value)
	if len(str) > length {
		return string(str[0:length])
	}

	return value
}

// StripLeading strips string of leading characters by an offset
func StripLeading(value string, offset int) string {

	str := []rune(value)
	if len(str) > offset {
		return string(str[offset:])
	}

	return value
}

// Elem returns the element of an array as specified by the index
//
// If the index exceeds the length of an array, it will return a non-nil value of the type.
// To monitor if the element exists, define a boolean value in the exists parameter
//
// This function requires version 1.18+
func Elem[T any](array *[]T, index int, exists *bool) T {

	var result T

	if exists != nil {
		*exists = false
	}

	if array == nil {
		return result
	}

	arrl := len(*array)
	if arrl == 0 {
		return result
	}
	arrl--

	if arrl >= index {
		if exists != nil {
			*exists = true
		}
		return (*array)[index]
	}

	return result
}

// ElemPtr returns a pointer to the element of an array as specified by the index
//
// If the index exceeds the length of an array, it will return a non-nil value of the type.
// To monitor if the element exists, define a boolean value in the exists parameter
//
// This function requires version 1.18+
func ElemPtr[T any](array *[]T, index int, exists *bool) *T {
	r := Elem(array, index, exists)
	return &r
}

// GetZero gets the zero value of the types defined as
// constraints.Ordered, time.Time and shopspring/decimal
//
// This function requires version 1.18+
func GetZero[T FieldTypeConstraint]() T {
	var result T
	return result
}
