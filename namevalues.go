package stdutil

import (
	"strconv"
	"strings"
)

//NameValue - a struct to manage value objects
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

	nn := strings.ToLower(name)

	for i, nv := range nvp.Pair {
		if nn == nv.Name {
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

	nn := strings.ToLower(name)

	for _, nv := range nvp.Pair {
		if nn == nv.Name {
			return AnyToString(nv.Value), true
		}
	}

	return "", false
}

// Int returns the name value as int. The second argument returns the existence.
func (nvp *NameValues) Int(name string) (int, bool) {

	if !nvp.prepared {
		nvp.prepare()
	}

	nn := strings.ToLower(name)

	for _, nv := range nvp.Pair {
		if nn == nv.Name {
			v, _ := strconv.Atoi(AnyToString(nv.Value))
			return v, true
		}
	}
	return 0, false
}

// Int64 returns the name value as int64. The second argument returns the existence.
func (nvp *NameValues) Int64(name string) (int64, bool) {

	if !nvp.prepared {
		nvp.prepare()
	}

	nn := strings.ToLower(name)

	for _, nv := range nvp.Pair {
		if nn == nv.Name {
			v, _ := strconv.ParseInt(AnyToString(nv.Value), 10, 64)
			return v, true
		}
	}
	return 0, false
}

// Plain returns the name value as interface{}. The second argument returns the existence.
func (nvp *NameValues) Plain(name string) (interface{}, bool) {

	if !nvp.prepared {
		nvp.prepare()
	}

	nn := strings.ToLower(name)

	for _, nv := range nvp.Pair {
		if nn == nv.Name {
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

	nn := strings.ToLower(name)

	for _, nv := range nvp.Pair {
		if nn == nv.Name {
			vs := AnyToString(nv.Value)
			return (vs == "true" || vs == "yes" || vs == "1" || vs == "-1" || vs == "on"), true
		}
	}
	return false, false
}

// Float64 returns the name value as float64. The second argument returns the existence.
func (nvp *NameValues) Float64(name string) (float64, bool) {

	if !nvp.prepared {
		nvp.prepare()
	}

	nn := strings.ToLower(name)

	for _, nv := range nvp.Pair {
		if nn == nv.Name {
			v, _ := strconv.ParseFloat(AnyToString(nv.Value), 64)
			return v, true
		}
	}
	return 0, false
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

// ValueString - get the struct value as string
//
// Deprecated: Use String() instead
func (nvp *NameValues) ValueString(name string) string {
	value, _ := nvp.String(name)
	return value
}

// ValueInt - get the struct value as int
//
// Deprecated: Use Int() instead
func (nvp *NameValues) ValueInt(name string) int {
	value, _ := nvp.Int(name)
	return value
}

// ValueInt64 - get the struct value as int64
//
// Deprecated: Use Int64() instead
func (nvp *NameValues) ValueInt64(name string) int64 {
	value, _ := nvp.Int64(name)
	return value
}

// ValuePlain - get the struct value as interface{}
//
// Deprecated: Use Plain() instead
func (nvp *NameValues) ValuePlain(name string) interface{} {
	value, _ := nvp.Plain(name)
	return value
}

// ValueBool - get the struct value as bool
//
// Deprecated: Use Bool() instead
func (nvp *NameValues) ValueBool(name string) bool {
	value, _ := nvp.Bool(name)
	return value
}

// ValueFloat64 - get the struct value as int
//
// Deprecated: Use Float64() instead
func (nvp *NameValues) ValueFloat64(name string) float64 {
	value, _ := nvp.Float64(name)
	return value
}

// **************************************************************
//   Deprecated pointer functions
// **************************************************************

// ValuePtrString - get the struct value of string pointer
//
// Deprecated: Use PtrString() instead
func (nvp *NameValues) ValuePtrString(name string) *string {
	value, _ := nvp.String(name)
	return &value
}

// ValuePtrInt - get the struct value as int pointer
//
// Deprecated: Use PtrInt() instead
func (nvp *NameValues) ValuePtrInt(name string) *int {
	value, _ := nvp.Int(name)
	return &value
}

// ValuePtrInt64 - get the struct value as int pointer
//
// Deprecated: Use PtrInt64() instead
func (nvp *NameValues) ValuePtrInt64(name string) *int64 {
	value, _ := nvp.Int64(name)
	return &value
}

// ValuePtrPlain - get the struct value as interface{} pointer
//
// Deprecated: Use PtrPlain() instead
func (nvp *NameValues) ValuePtrPlain(name string) *interface{} {
	value, _ := nvp.Plain(name)
	return &value
}

//ValuePtrBool - get the struct value as bool as pointer
//
// Deprecated: Use PtrBool() instead
func (nvp *NameValues) ValuePtrBool(name string) *bool {
	value, _ := nvp.Bool(name)
	return &value
}

// ValuePtrFloat64 - get the struct value as int as pointer
//
// Deprecated: Use PtrFloat64() instead
func (nvp *NameValues) ValuePtrFloat64(name string) *float64 {
	value, _ := nvp.Float64(name)
	return &value
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
