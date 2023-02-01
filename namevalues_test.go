package stdutil

import "testing"

func TestNameValues(t *testing.T) {

	nvs := NameValues{
		Pair: map[string]any{
			"name":   "Zaldy",
			"band":   "Razzie",
			"active": false,
			"age":    "48",
			"man":    "true",
		},
	}

	vs := Value[string](nvs, "name")
	t.Log(vs)

	vb := Value[bool](nvs, "active")
	t.Log(vb)

	vi := Value[int](nvs, "age")
	t.Log(vi)

	vm := Value[bool](nvs, "man")
	t.Log(vm)
}
