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
	Pair []NameValue
}

// KeyExists - checks if the key or name exists
func (nvp *NameValues) KeyExists(name string) bool {
	nn := strings.ToLower(name)
	for _, nv := range nvp.Pair {
		if nn == strings.ToLower(nv.Name) {
			return true
		}
	}
	return false
}

// ValueString - get the struct value as string
func (nvp *NameValues) ValueString(name string) string {
	nn := strings.ToLower(name)
	for _, nv := range nvp.Pair {
		if nn == strings.ToLower(nv.Name) {
			return AnyToString(nv.Value)
		}
	}
	return ""
}

// ValueInt - get the struct value as int
func (nvp *NameValues) ValueInt(name string) int {
	nn := strings.ToLower(name)
	for _, nv := range nvp.Pair {
		if nn == strings.ToLower(nv.Name) {
			v, _ := strconv.Atoi(AnyToString(nv.Value))
			return v
		}
	}
	return 0
}

// ValuePlain - get the struct value as interface{}
func (nvp *NameValues) ValuePlain(name string) interface{} {
	nn := strings.ToLower(name)
	for _, nv := range nvp.Pair {
		if nn == strings.ToLower(nv.Name) {
			return nv.Value
		}
	}
	return nil
}

// ValueBool - get the struct value as bool
func (nvp *NameValues) ValueBool(name string) bool {
	nn := strings.ToLower(name)
	for _, nv := range nvp.Pair {
		if nn == strings.ToLower(nv.Name) {
			vs := AnyToString(nv.Value)
			return (vs == "true" || vs == "yes" || vs == "1" || vs == "on")
		}
	}
	return false
}

// ValuePtrString - get the struct value of string pointer
func (nvp *NameValues) ValuePtrString(name string) *string {
	result := nvp.ValueString(name)
	return &result
}

// ValuePtrInt - get the struct value as int pointer
func (nvp *NameValues) ValuePtrInt(name string) *int {
	result := nvp.ValueInt(name)
	return &result
}

// ValueInt64 - get the struct value as int64
func (nvp *NameValues) ValueInt64(name string) int64 {
	for _, nv := range nvp.Pair {
		if strings.ToLower(name) == strings.ToLower(nv.Name) {
			v, _ := strconv.ParseInt(AnyToString(nv.Value), 10, 64)
			return v
		}
	}
	return 0
}

// ValuePtrInt64 - get the struct value as int pointer
func (nvp *NameValues) ValuePtrInt64(name string) *int64 {
	result := nvp.ValueInt64(name)
	return &result
}

// ValuePtrPlain - get the struct value as interface{} pointer
func (nvp *NameValues) ValuePtrPlain(name string) *interface{} {
	result := nvp.ValuePlain(name)
	return &result
}

// ValueFloat64 - get the struct value as int
func (nvp *NameValues) ValueFloat64(name string) float64 {
	for _, nv := range nvp.Pair {
		if strings.ToLower(name) == strings.ToLower(nv.Name) {
			v, _ := strconv.ParseFloat(AnyToString(nv.Value), 32)
			return v
		}
	}
	return 0
}

// ValuePtrFloat64 - get the struct value as int as pointer
func (nvp *NameValues) ValuePtrFloat64(name string) *float64 {
	result := nvp.ValueFloat64(name)
	return &result
}

//ValuePtrBool - get the struct value as bool as pointer
func (nvp *NameValues) ValuePtrBool(name string) *bool {
	result := nvp.ValueBool(name)
	return &result
}

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
