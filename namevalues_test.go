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

	vs := NameValueGet[string](nvs, "name")
	t.Log(vs)

	vb := NameValueGet[bool](nvs, "active")
	t.Log(vb)

	vi := NameValueGet[int](nvs, "age")
	t.Log(vi)

	vm := NameValueGet[bool](nvs, "man")
	t.Log(vm)
}
