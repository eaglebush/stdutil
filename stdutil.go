// Package stdutil is a collection of often-needed functions for development use
//
// The package is still supports v1.19.
package stdutil

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	ssd "github.com/shopspring/decimal"
	"golang.org/x/exp/constraints"
)

type (
	FieldTypeConstraint interface {
		constraints.Ordered | time.Time | ssd.Decimal | bool | byte
	}
	NumericConstraint interface {
		constraints.Integer | constraints.Float
	}
)

const (
	INTERPOLATE_PATTERN string = `\$\{(\w*)\}` // search for ${*}
	EMAIL_PATTERN       string = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9]" +
		"(?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
)

type (
	StringValidationOptions struct {
		Empty    bool // Allow empty string. Default: false, will raise an error if the string is empty
		Null     bool // Allow null. Default: false, will raise an error if the string is null
		Min      int  // Minimum length. Default: 0
		Max      int  // Maximum length. Default: 0
		NoSpaces bool // Do not allow spaces in the string. Default: false. Setting to true will raise an error if the string has spaces
		Extended []func(value *string) error
	}
	TimeValidationOptions struct {
		Null     bool       // Allow null. Default: false, will raise an error if the time is null
		Empty    bool       // Allow zero time Default: false, will raise an error if the time is zero
		Min      *time.Time // Minimum time. Default: nil
		Max      *time.Time // Maximum time. Default: nil
		DateOnly bool       // Compare dates only. Default: false
		Extended []func(value *time.Time) error
	}
	NumericValidationOptions[T NumericConstraint] struct {
		Null     bool // Allow null. Default: false, will raise an error if the time is null
		Empty    bool // Allow zero time Default: false, will raise an error if the time is zero
		Min      T    // Minimum time. Default: nil
		Max      T    // Maximum time. Default: nil
		Extended []func(value *T) error
	}
	DecimalValidationOptions struct {
		Null     bool         // Allow null. Default: false, will raise an error if the decimal is null
		Empty    bool         // Allow zero decimal. Default: false, will raise an error if the decimal is zero
		Min      *ssd.Decimal // Minimum decimal value. Default: nil
		Max      *ssd.Decimal // Maximum decimal value. Default: nil
		Extended []func(value *ssd.Decimal) error
	}
	SeriesOptions struct {
		Prefix string // Prefix of series
		Suffix string // Suffix of series
		Length int    // Fixed length of the series
	}
)

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
		if t == nil || !*t {
			return "false"
		}
		return "true"
	case *time.Time:
		if t == nil {
			return "'" + time.Time{}.Format(time.RFC3339) + "'"
		}
		tm := *t
		b = "'" + tm.Format(time.RFC3339) + "'"
	}

	return b
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

// Elem returns the element of an array as specified by the index
//
// If the index exceeds the length of an array, it will return a non-nil value of the type.
// To monitor if the element exists, define a boolean value in the exists parameter
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
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
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func ElemPtr[T any](array *[]T, index int, exists *bool) *T {
	r := Elem(array, index, exists)
	return &r
}

// GetZero gets the zero value of the type.
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func GetZero[T FieldTypeConstraint]() T {
	var result T
	return result
}

// If is a basic ternary operator to return whatever is set in
// truthy and falsey parameter.
// If the subject is nil, empty string, 0, -0 or false, it will return the falsey parameter
//
// This function requires version 1.18+
func If[T constraints.Ordered](subject any, truthy T, falsey T) T {
	if subject == nil {
		return falsey
	}
	switch t := subject.(type) {
	case string:
		if t == "" {
			return falsey
		}
	case *string:
		if t == nil || *t == "" {
			return falsey
		}
	case
		int8, int16, int32, int64, int,
		uint8, uint16, uint32, uint64, uint,
		float32, float64, complex64, complex128:
		if t == 0 || t == -0 {
			return falsey
		}
	case
		*int8, *int16, *int32, *int64, *int,
		*uint8, *uint16, *uint32, *uint64, *uint,
		*float32, *float64, *complex64, *complex128:
		vo := reflect.ValueOf(t)
		tx := vo.Elem()
		if !tx.IsValid() || tx.IsZero() {
			return falsey
		}
	case bool:
		if !t {
			return falsey
		}
	case *bool:
		if !*t {
			return falsey
		}
	}
	return truthy
}

// In checks if the seek parameter is in the list parameter
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func In[T comparable](seek T, list ...T) bool {
	for _, li := range list {
		if li == seek {
			return true
		}
	}
	return false
}

// IsNullOrEmpty checks for emptiness of a pointer variable ignoring nullity
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func IsEmpty[T FieldTypeConstraint](value *T) bool {
	// if value == nil {
	// 	return false
	// }
	// if *value == GetZero[T]() {
	// 	return true
	// }
	// return false
	return value != nil && *value == GetZero[T]()
}

// IsNullOrEmpty checks for nullity and emptiness of a pointer variable
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func IsNullOrEmpty[T FieldTypeConstraint](value *T) bool {
	return value == nil || *value == GetZero[T]()
}

// IsNumeric checks if a string is numeric
func IsNumeric(value string) error {
	if value == "" {
		return fmt.Errorf("is empty")
	}
	if _, err := strconv.ParseFloat(value, 64); err != nil {
		return fmt.Errorf(`is not a number (%s)`, err)
	}
	return nil
}

// Interpolate interpolates string with the name value pairs
func Interpolate(base string, nv NameValues) (string, []any) {
	var (
		val  any
		sval string
	)

	nstr := base
	re := regexp.MustCompile(INTERPOLATE_PATTERN)
	matches := re.FindAllString(base, -1)
	vals := make([]any, len(matches))
	for i, match := range matches {
		val = "0"
		sval = "0"
		for n, v := range nv.Pair {
			if strings.EqualFold(match, `${`+n+`}`) {
				sval = AnyToString(v)
				val = v
				break
			}
		}
		nstr = strings.Replace(nstr, match, sval, -1)
		vals[i] = val //a string 0 would cater to both string and number columns
	}
	return nstr, vals
}

// MapVal retrieves a value from a map by a key and converts it to the type indicated by T.
// Returns a pointer to the value if found, or nil if not found.
//
// The third parameter, dateLayout can be set with many time layouts. The specified layouts
// will be the only one to try parsing.
//
// If not set, the built-in date layouts of ParseDate are used. See function for supported layouts
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func MapVal[T FieldTypeConstraint](kvmap *map[string]any, key string, dateLayout ...string) *T {
	var (
		ok bool
	)
	if kvmap == nil {
		return nil
	}
	miv, ok := (*kvmap)[key]
	if !ok {
		return nil
	}
	mv, ok := miv.(T)
	if ok {
		return &mv
	}
	t := new(T)  // initialize a variable to the return generic type
	a := any(*t) // initialize a variable to the value type of the generic type
	switch a.(type) {
	case time.Time:
		if s, ok := miv.(string); ok {
			var dlo *string
			if len(dateLayout) > 0 {
				dlo = &dateLayout[0]
			}
			v, _, err := ParseDate(s, dlo)
			if err != nil {
				return nil
			}
			*t = any(v).(T) // Convert the parsed value to any before asserting the type for the return
			return t
		}
	case ssd.Decimal:
		if s, ok := miv.(string); ok {
			v, err := ssd.NewFromString(s)
			if err != nil {
				return nil
			}
			*t = any(v).(T)
			return t
		}
	default:

	}
	return nil
}

// New initializes a variable and returns a pointer of its type
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func New[T FieldTypeConstraint](value T) *T {
	n := new(T)
	*n = value
	return n
}

// NonNullComp compares two parameters when both are not nil.
//
//   - When one or both of the parameters is nil, the function returns -1
//   - When the parameters are equal, the function returns 0.
//   - else it returns 1
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func NonNullComp[T FieldTypeConstraint](param1 *T, param2 *T) int {
	if param1 == nil || param2 == nil {
		return -1
	}
	if *param1 == *param2 {
		return 0
	}
	return 1
}

// Null accepts a value to test and the default value
// if it fails. It returns a non-pointer value of T.
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func Null[T any](testValue any, defaultValue any) T {
	var (
		def T
	)
	if defaultValue == nil {
		defaultValue = def
	}
	if testValue == nil {
		return defaultValue.(T)
	}
	vo := reflect.ValueOf(testValue)
	if k := vo.Kind(); k == reflect.Map ||
		k == reflect.Func ||
		k == reflect.Ptr ||
		k == reflect.Slice ||
		k == reflect.Interface {
		if vo.IsZero() && vo.IsNil() {
			return defaultValue.(T)
		}
		ifv := vo.Elem().Interface()
		return ifv.(T)
	}
	return def
}

// NullPtr accepts a value to test and the default value
// if it fails. It returns a pointer value of T.
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func NullPtr[T any](testValue any, defaultValue any) *T {
	val := Null[T](testValue, defaultValue)
	return &val
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

// ParseDate parses a string as date.
//
// If dateLayout is not provided, this function try all layout combinations.
// The following date layouts has been provided:
//   - "2006-01-02"
//   - "2006-02-01"
//   - "01-02-2006"
//   - "02-01-2006"
//   - "2006/01/02"
//   - "2006/02/01"
//   - "01/02/2006"
//   - "02/01/2006"
//   - "1/2/2006"
//   - "2/1/2006"
//   - "2006/1/2"
//   - "2006/2/1"
//   - "06-01-02"
//   - "06-02-01"
//   - "01-02-06"
//   - "02-01-06"
//   - "06/01/02"
//   - "06/02/01"
//   - "01/02/06"
//   - "02/01/06"
//
// The date layout partitions means:
//   - Anything with 1, with or without zero is the month
//   - Anything with 2, with or without zero is the day
//   - Anything with 06, with or without the prefix 20 is the year
func ParseDate(dtText string, dateLayout *string) (time.Time, string, error) {
	var (
		rtm time.Time
		rlo string
		err error
	)
	if dtText == "" {
		return rtm, rlo, fmt.Errorf("invalid date or time input")
	}
	dlo :=
		[]string{
			// Dashed full year
			"2006-01-02",
			"2006-02-01",
			"01-02-2006",
			"02-01-2006",
			// Forward-slashed full year
			"2006/01/02",
			"2006/02/01",
			"01/02/2006",
			"02/01/2006",
			// Forward-slashed single digit day and month
			"1/2/2006",
			"2/1/2006",
			"2006/1/2",
			"2006/2/1",
			// 2-digit dashed year with full digit day and month
			"06-01-02",
			"06-02-01",
			"01-02-06",
			"02-01-06",
			// 2-digit forward-slashed year with full digit day and month
			"06/01/02",
			"06/02/01",
			"01/02/06",
			"02/01/06",
		}
	// Try to parse using layout provided
	// The function will return upon failure or success
	if dateLayout != nil {
		// Check if layout is in the array
		if !In(*dateLayout, dlo...) {
			return rtm, *dateLayout, fmt.Errorf("layout provided not supported")
		}
		rtm, err = time.Parse(*dateLayout, dtText)
		return rtm, *dateLayout, err
	}
	// Try each layout until it succeeds
	for _, lo := range dlo {
		rtm, err = time.Parse(lo, dtText)
		if err != nil {
			continue
		}
		return rtm, lo, err
	}
	return rtm, rlo, fmt.Errorf("date parsing failed")
}

// SafeMapWrite allows writing to maps by locking, preventing the library from crashing
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func SafeMapWrite[T any](ptrMap *map[string]T, key string, value T, rw *sync.RWMutex) bool {
	defer func() {
		recover()
	}()
	// Prepare mutex
	// attempt writing to map
	if rw.TryLock() {
		defer rw.Unlock()
		(*ptrMap)[key] = value
	}
	return true
}

// SafeMapRead allows reading maps by locking it, preventing the library from crashing
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func SafeMapRead[T any](ptrMap *map[string]T, key string, rw *sync.RWMutex) T {
	var result T
	defer func() {
		recover()
	}()
	if rw.TryLock() {
		defer rw.Unlock()
		result = (*ptrMap)[key]
	}
	return result
}

// Seek checks if the seek parameter is in the list parameter and returns it.
// If the value is not found in the list, the function returns nil
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func Seek[T comparable](seek T, list ...T) *T {
	for _, li := range list {
		if li == seek {
			return &li
		}
	}
	return nil
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
	ret := NameValues{
		Pair: make(map[string]any),
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

// StringToByte converts string or strings to byte array with an option for a separator
func StringToByte(sep string, elems ...string) []byte {
	return []byte(strings.Join(elems, sep))
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

// StripLeading strips string of leading characters by an offset
func StripLeading(value string, offset int) string {
	str := []rune(value)
	if len(str) > offset {
		return string(str[offset:])
	}
	return value
}

// StripTrailing strips string of trailing characters after the length
func StripTrailing(value string, length int) string {
	str := []rune(value)
	if len(str) > length {
		return string(str[0:length])
	}
	return value
}

// ToInterfaceArray converts a value to interface array
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func ToInterfaceArray[T FieldTypeConstraint](values T) []any {
	var value [1]any
	value[0] = values
	return value[:]
}

// Val gets the value of a pointer in order
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func Val[T FieldTypeConstraint](value *T) T {
	if value == nil {
		return GetZero[T]()
	}
	return *value
}

// ValidateEmail validates an e-mail address
func ValidateEmail(email *string) error {
	if email == nil || *email == "" {
		return fmt.Errorf("is an invalid email address")
	}
	re := regexp.MustCompile(EMAIL_PATTERN)
	if !re.MatchString(*email) {
		return fmt.Errorf("is an invalid email address")
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
//
// Currently supported data types are:
//   - constraints.Integer (Signed | Unsigned)
//   - constraints.Float (~float32 | ~float64)
//
// This function requires version 1.18+
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
