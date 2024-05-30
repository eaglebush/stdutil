package stdutil

import (
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

type (
	// NameValue is a struct to contain a name-value pair
	NameValue[T any] struct {
		Name  string `json:"name,omitempty"`
		Value T      `json:"value,omitempty"`
	}
	// NameValues is a struct to manage value structs
	NameValues struct {
		Pair     map[string]any
		prepared bool
	}
)

func (nvp *NameValues) prepare() {
	np := make(map[string]any)
	for n := range nvp.Pair {
		ln := strings.ToLower(n)
		np[ln] = nvp.Pair[n]
		delete(nvp.Pair, n)
	}
	for n := range np {
		nvp.Pair[n] = np[n]
	}
	nvp.prepared = true
}

// Exists checks if the key or name exists. It returns the index of the element if found, -1 if not found.
func (nvp *NameValues) Exists(name string) bool {
	if !nvp.prepared {
		nvp.prepare()
	}
	name = strings.ToLower(name)
	_, exists := nvp.Pair[name]
	return exists
}

// NameValueGet gets the value from the collection of NameValues by name
//
// This function requires version 1.18+
func NameValueGet[T constraints.Ordered | bool](nvs NameValues, name string) T {
	if !nvs.prepared {
		nvs.prepare()
	}

	name = strings.ToLower(name)
	tmp := nvs.Pair[name]

	tpt := any(*new(T))
	value := *new(T)

	// If the value is a string and the inferred type is otherwise
	// try to convert, else just convery via inferred type
	switch t := tmp.(type) {
	case string:
		switch tpt.(type) {
		case int:
			val, _ := strconv.ParseInt(t, 10, 32)
			value = any(int(val)).(T)
		case int64:
			val, _ := strconv.ParseInt(t, 10, 64)
			value = any(val).(T)
		case bool:
			val, _ := strconv.ParseBool(t)
			value = any(val).(T)
		case float32:
			val, _ := strconv.ParseFloat(t, 32)
			value = any(val).(T)
		case float64:
			val, _ := strconv.ParseFloat(t, 64)
			value = any(val).(T)
		default:
			if tmp != nil {
				value = tmp.(T)
			} else {
				value = GetZero[T]()
			}
		}
	default:
		if t != nil {
			value = t.(T)
		} else {
			value = GetZero[T]()
		}
	}

	return value
}

// NameValueGetPtr gets the value from the collection of NameValues by name as pointer
//
// This function requires version 1.18+
func NameValueGetPtr[T constraints.Ordered | bool](nvs NameValues, name string) *T {
	value := NameValueGet[T](nvs, name)
	return &value
}

// String returns the name value as string. The second result returns the existence.
func (nvp *NameValues) String(name string) (string, bool) {
	if !nvp.prepared {
		nvp.prepare()
	}
	name = strings.ToLower(name)
	tmp, exists := nvp.Pair[name]
	value, _ := tmp.(string)
	return value, exists
}

// Strings returns the values as a string array.
// If the value is comma-separated, all elements delimited by the comma
// will be returned as an array of string
func (nvp *NameValues) Strings(name string) []string {
	value, _ := nvp.String(name)
	if strings.Contains(value, ",") {
		return strings.Split(value, ",")
	}
	return []string{value}
}

// Int returns the name value as int. The second result returns the existence.
func (nvp *NameValues) Int(name string) (int, bool) {
	if !nvp.prepared {
		nvp.prepare()
	}
	var value int
	name = strings.ToLower(name)
	tmp, exists := nvp.Pair[name]
	if exists {
		val, _ := strconv.ParseInt(tmp.(string), 10, 32)
		value = int(val)
	}
	return value, exists
}

// Ints returns the values as an int array
func (nvp *NameValues) Ints(name string) []int {
	value, _ := nvp.Int(name)
	return []int{value}
}

// Int64 returns the name value as int64. The second result returns the existence.
func (nvp *NameValues) Int64(name string) (int64, bool) {
	if !nvp.prepared {
		nvp.prepare()
	}
	var value int64
	name = strings.ToLower(name)
	tmp, exists := nvp.Pair[name]
	if exists {
		value, _ = strconv.ParseInt(tmp.(string), 10, 64)
	}
	return value, exists
}

// Int64s returns the values as an int64 array
func (nvp *NameValues) Int64s(name string) []int64 {
	value, _ := nvp.Int64(name)
	return []int64{value}
}

// Plain returns the name value as interface{}. The second result returns the existence.
func (nvp *NameValues) Plain(name string) (interface{}, bool) {
	if !nvp.prepared {
		nvp.prepare()
	}
	name = strings.ToLower(name)
	tmp, exists := nvp.Pair[name]
	return tmp, exists
}

// Bool returns the name value as boolean. It automatically convers 'true', 'yes', '1', '-1' and 'on' to boolean The second result returns the existence.
func (nvp *NameValues) Bool(name string) (bool, bool) {
	value, _ := nvp.String(name)
	return (value == "true" || value == "yes" || value == "1" || value == "-1" || value == "on"), true
}

// Bools returns the values as a boolean array
func (nvp *NameValues) Bools(name string) []bool {
	value, _ := nvp.Bool(name)
	return []bool{value}
}

// Float64 returns the name value as float64. The second result returns the existence.
func (nvp *NameValues) Float64(name string) (float64, bool) {
	if !nvp.prepared {
		nvp.prepare()
	}
	var value float64
	name = strings.ToLower(name)
	tmp, exists := nvp.Pair[name]
	if exists {
		value, _ = strconv.ParseFloat(tmp.(string), 64)
	}
	return value, exists
}

// Float64s returns the values as a float64 array
func (nvp *NameValues) Float64s(name string) []float64 {
	value, _ := nvp.Float64(name)
	return []float64{value}
}

// **************************************************************
//   Pointer outputs
// **************************************************************

// PtrString returns the name value as pointer to string. The second result returns the existence.
func (nvp *NameValues) PtrString(name string) (*string, bool) {
	value, exists := nvp.String(name)
	return &value, exists
}

// PtrInt returns the name value as pointer to int. The second result returns the existence.
func (nvp *NameValues) PtrInt(name string) (*int, bool) {
	value, exists := nvp.Int(name)
	return &value, exists
}

// PtrInt64 returns the name value as pointer to int64. The second result returns the existence.
func (nvp *NameValues) PtrInt64(name string) (*int64, bool) {
	value, exists := nvp.Int64(name)
	return &value, exists
}

// PtrPlain returns the name value as pointer to interface{}. The second result returns the existence.
func (nvp *NameValues) PtrPlain(name string) (*interface{}, bool) {
	value, exists := nvp.Plain(name)
	return &value, exists
}

// PtrBool returns the name value as pointer to bool. The second result returns the existence.
func (nvp *NameValues) PtrBool(name string) (*bool, bool) {
	value, exists := nvp.Bool(name)
	return &value, exists
}

// PtrFloat64 returns the name value as pointer to int64. The second result returns the existence.
func (nvp *NameValues) PtrFloat64(name string) (*float64, bool) {
	value, exists := nvp.Float64(name)
	return &value, exists
}

// **************************************************************
//   Miscellaneous functions
// **************************************************************

// ToInterfaceArray - converts name values to interface array
func (nvp *NameValues) ToInterfaceArray() []interface{} {
	return NameValuesToInterfaceArray(*nvp)
}

// Interpolate - interpolate string with values from with base string
func (nvp *NameValues) Interpolate(base string) (string, []interface{}) {
	return Interpolate(base, *nvp)
}

// SortByKey sort name values by key order array
func (nvp *NameValues) SortByKey(keyOrder *[]string) NameValues {
	return SortByKey(nvp, keyOrder)
}
