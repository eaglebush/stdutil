package stdutil

import (
	"strings"
)

// NameValues - a struct to manage value structs
type NameValues struct {
	Pair     map[string]any
	prepared bool
}

func (nvp *NameValues) prepare() {

	for n := range nvp.Pair {
		ln := strings.ToLower(n)
		nvp.Pair[ln] = nvp.Pair[n]
		delete(nvp.Pair, n)
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

// String returns the name value as string. The second argument returns the existence.
func (nvp *NameValues) String(name string) (string, bool) {

	if !nvp.prepared {
		nvp.prepare()
	}

	name = strings.ToLower(name)
	tmp, exists := nvp.Pair[name]
	value, _ := tmp.(string)
	return value, exists
}

// Strings returns the values as a string array
func (nvp *NameValues) Strings(name string) []string {

	if !nvp.prepared {
		nvp.prepare()
	}

	value, _ := nvp.String(name)
	return []string{value}
}

// Int returns the name value as int. The second argument returns the existence.
func (nvp *NameValues) Int(name string) (int, bool) {

	if !nvp.prepared {
		nvp.prepare()
	}

	name = strings.ToLower(name)
	tmp, exists := nvp.Pair[name]
	value, _ := tmp.(int)
	return value, exists
}

// Ints returns the values as an int array
func (nvp *NameValues) Ints(name string) []int {

	if !nvp.prepared {
		nvp.prepare()
	}

	value, _ := nvp.Int(name)
	return []int{value}
}

// Int64 returns the name value as int64. The second argument returns the existence.
func (nvp *NameValues) Int64(name string) (int64, bool) {

	if !nvp.prepared {
		nvp.prepare()
	}

	name = strings.ToLower(name)
	tmp, exists := nvp.Pair[name]
	value, _ := tmp.(int64)
	return value, exists
}

// Int64s returns the values as an int64 array
func (nvp *NameValues) Int64s(name string) []int64 {

	if !nvp.prepared {
		nvp.prepare()
	}

	value, _ := nvp.Int64(name)
	return []int64{value}
}

// Plain returns the name value as interface{}. The second argument returns the existence.
func (nvp *NameValues) Plain(name string) (interface{}, bool) {

	if !nvp.prepared {
		nvp.prepare()
	}

	name = strings.ToLower(name)
	tmp, exists := nvp.Pair[name]
	return tmp, exists
}

// Bool returns the name value as boolean. It automatically convers 'true', 'yes', '1', '-1' and 'on' to boolean The second argument returns the existence.
func (nvp *NameValues) Bool(name string) (bool, bool) {

	if !nvp.prepared {
		nvp.prepare()
	}

	value, _ := nvp.String(name)
	return (value == "true" || value == "yes" || value == "1" || value == "-1" || value == "on"), true
}

// Bools returns the values as a boolean array
func (nvp *NameValues) Bools(name string) []bool {

	if !nvp.prepared {
		nvp.prepare()
	}

	value, _ := nvp.Bool(name)
	return []bool{value}
}

// Float64 returns the name value as float64. The second argument returns the existence.
func (nvp *NameValues) Float64(name string) (float64, bool) {

	if !nvp.prepared {
		nvp.prepare()
	}

	name = strings.ToLower(name)
	tmp, exists := nvp.Pair[name]
	value, _ := tmp.(float64)
	return value, exists
}

// Float64s returns the values as a float64 array
func (nvp *NameValues) Float64s(name string) []float64 {

	if !nvp.prepared {
		nvp.prepare()
	}

	value, _ := nvp.Float64(name)
	return []float64{value}
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
//   Miscellaneous functions
// **************************************************************

// ToInterfaceArray - converts name values to interface array
func (nvp *NameValues) ToInterfaceArray() []interface{} {
	return NameValuesToInterfaceArray(*nvp)
}

// ToVerifyExpressionArray - converts name values to verify expression array
func (nvp *NameValues) ToVerifyExpressionArray() []VerifyExpression {
	return NameValuesToVerifyExpressionArray(*nvp)
}

// Interpolate - interpolate string with values from with base string
func (nvp *NameValues) Interpolate(base string) (string, []interface{}) {
	return InterpolateString(base, *nvp)
}

// SortByKeyArray - sort name values by key order array
func (nvp *NameValues) SortByKeyArray(keyOrder *[]string) NameValues {
	return SortByKeyArray(nvp, keyOrder)
}
