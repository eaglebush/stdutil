package stdutil

import (
	"strconv"
	"strings"
)

// NameValue - a struct to manage value objects
type NameValue struct {
	Name  string
	Value interface{}
}

// NameValues - a struct to manage value structs
type NameValues struct {
	Pair     []NameValue
	prepared bool
}

// **************************************************************
//   New functions
// **************************************************************
func (nvp *NameValues) prepare() {

	for i := range nvp.Pair {
		nvp.Pair[i].Name = strings.ToLower(nvp.Pair[i].Name)
	}

	nvp.prepared = true
}

// Exists checks if the key or name exists. It returns the index of the element if found, -1 if not found.
func (nvp *NameValues) Exists(name string) int {

	if !nvp.prepared {
		nvp.prepare()
	}

	for i, nv := range nvp.Pair {
		if strings.EqualFold(name, nv.Name) {
			return i
		}
	}
	return -1
}

// String returns the name value as string. The second argument returns the existence.
func (nvp *NameValues) String(name string) (string, bool) {

	if !nvp.prepared {
		nvp.prepare()
	}

	for _, nv := range nvp.Pair {
		if strings.EqualFold(name, nv.Name) {
			return AnyToString(nv.Value), true
		}
	}

	return "", false
}

// Strings returns the values as a string array
func (nvp *NameValues) Strings(name string) []string {

	if !nvp.prepared {
		nvp.prepare()
	}

	str := make([]string, 0)

	for _, nv := range nvp.Pair {
		if strings.EqualFold(name, nv.Name) {
			str = append(str, AnyToString(nv.Value))
		}
	}

	return str
}

// Int returns the name value as int. The second argument returns the existence.
func (nvp *NameValues) Int(name string) (int, bool) {

	if !nvp.prepared {
		nvp.prepare()
	}

	for _, nv := range nvp.Pair {
		if strings.EqualFold(name, nv.Name) {
			v, _ := strconv.Atoi(AnyToString(nv.Value))
			return v, true
		}
	}
	return 0, false
}

// Ints returns the values as an int array
func (nvp *NameValues) Ints(name string) []int {

	if !nvp.prepared {
		nvp.prepare()
	}

	str := make([]int, 0)

	for _, nv := range nvp.Pair {
		if strings.EqualFold(name, nv.Name) {
			v, _ := strconv.Atoi(AnyToString(nv.Value))
			str = append(str, v)
		}
	}

	return str
}

// Int64 returns the name value as int64. The second argument returns the existence.
func (nvp *NameValues) Int64(name string) (int64, bool) {

	if !nvp.prepared {
		nvp.prepare()
	}

	for _, nv := range nvp.Pair {
		if strings.EqualFold(name, nv.Name) {
			v, _ := strconv.ParseInt(AnyToString(nv.Value), 10, 64)
			return v, true
		}
	}
	return 0, false
}

// Int64s returns the values as an int64 array
func (nvp *NameValues) Int64s(name string) []int64 {

	if !nvp.prepared {
		nvp.prepare()
	}

	str := make([]int64, 0)

	for _, nv := range nvp.Pair {
		if strings.EqualFold(name, nv.Name) {
			v, _ := strconv.ParseInt(AnyToString(nv.Value), 10, 64)
			str = append(str, v)
		}
	}

	return str
}

// Plain returns the name value as interface{}. The second argument returns the existence.
func (nvp *NameValues) Plain(name string) (interface{}, bool) {

	if !nvp.prepared {
		nvp.prepare()
	}

	for _, nv := range nvp.Pair {
		if strings.EqualFold(name, nv.Name) {
			return nv.Value, true
		}
	}
	return nil, false
}

// Bool returns the name value as boolean. It automatically convers 'true', 'yes', '1', '-1' and 'on' to boolean The second argument returns the existence.
func (nvp *NameValues) Bool(name string) (bool, bool) {

	if !nvp.prepared {
		nvp.prepare()
	}

	for _, nv := range nvp.Pair {
		if strings.EqualFold(name, nv.Name) {
			vs := AnyToString(nv.Value)
			return (vs == "true" || vs == "yes" || vs == "1" || vs == "-1" || vs == "on"), true
		}
	}
	return false, false
}

// Bools returns the values as a boolean array
func (nvp *NameValues) Bools(name string) []bool {

	if !nvp.prepared {
		nvp.prepare()
	}

	str := make([]bool, 0)

	for _, nv := range nvp.Pair {
		if strings.EqualFold(name, nv.Name) {
			vs := AnyToString(nv.Value)
			str = append(str, (vs == "true" || vs == "yes" || vs == "1" || vs == "-1" || vs == "on"))
		}
	}

	return str
}

// Float64 returns the name value as float64. The second argument returns the existence.
func (nvp *NameValues) Float64(name string) (float64, bool) {

	if !nvp.prepared {
		nvp.prepare()
	}

	for _, nv := range nvp.Pair {
		if strings.EqualFold(name, nv.Name) {
			v, _ := strconv.ParseFloat(AnyToString(nv.Value), 64)
			return v, true
		}
	}
	return 0, false
}

// Float64s returns the values as a float64 array
func (nvp *NameValues) Float64s(name string) []float64 {

	if !nvp.prepared {
		nvp.prepare()
	}

	str := make([]float64, 0)

	for _, nv := range nvp.Pair {
		if strings.EqualFold(name, nv.Name) {
			v, _ := strconv.ParseFloat(AnyToString(nv.Value), 64)
			str = append(str, v)
		}
	}

	return str
}

// **************************************************************
//   Pointer outputs
// **************************************************************

// PtrString returns the name value as pointer to string. The second argument returns the existence.
func (nvp *NameValues) PtrString(name string) (*string, bool) {
	value, exists := nvp.String(name)
	return &value, exists
}

// PtrInt returns the name value as pointer to int. The second argument returns the existence.
func (nvp *NameValues) PtrInt(name string) (*int, bool) {
	value, exists := nvp.Int(name)
	return &value, exists
}

// PtrInt64 returns the name value as pointer to int64. The second argument returns the existence.
func (nvp *NameValues) PtrInt64(name string) (*int64, bool) {
	value, exists := nvp.Int64(name)
	return &value, exists
}

// PtrPlain returns the name value as pointer to interface{}. The second argument returns the existence.
func (nvp *NameValues) PtrPlain(name string) (*interface{}, bool) {
	value, exists := nvp.Plain(name)
	return &value, exists
}

// PtrBool returns the name value as pointer to bool. The second argument returns the existence.
func (nvp *NameValues) PtrBool(name string) (*bool, bool) {
	value, exists := nvp.Bool(name)
	return &value, exists
}

// PtrFloat64 returns the name value as pointer to int64. The second argument returns the existence.
func (nvp *NameValues) PtrFloat64(name string) (*float64, bool) {
	value, exists := nvp.Float64(name)
	return &value, exists
}

// **************************************************************
//   Deprecated functions
// **************************************************************

// KeyExists - checks if the key or name exists
//
// Deprecated: Use Exists() instead
func (nvp *NameValues) KeyExists(name string) bool {
	return nvp.Exists(name) != -1
}

// **************************************************************
//   Miscellaneous functions
// **************************************************************

// ToInterfaceArray - converts name values to interface array
func (nvp *NameValues) ToInterfaceArray() []interface{} {
	return NameValuesToInterfaceArray(*nvp)
}

// ToValidationExpressionArray - converts name values to validation expression array
func (nvp *NameValues) ToValidationExpressionArray() []ValidationExpression {
	return NameValuesToValidationExpressionArray(*nvp)
}

// Interpolate - interpolate string with values from with base string
func (nvp *NameValues) Interpolate(base string) (string, []interface{}) {
	return InterpolateString(base, *nvp)
}

// SortByKeyArray - sort name values by key order array
func (nvp *NameValues) SortByKeyArray(keyOrder *[]string) NameValues {
	return SortByKeyArray(nvp, keyOrder)
}
