package stdutil

import (
	"testing"

	ssd "github.com/shopspring/decimal"
)

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

func TestNVString(t *testing.T) {
	var (
		str    string
		strs   []string
		strp   *string
		exists bool
	)
	mv := map[string]any{
		"key1": "This is a string",
		"key2": "This is a string, this is madness",
	}
	nvs := NameValues{
		Pair: mv,
	}

	// Get value of key
	str, exists = nvs.String("key1")
	t.Log("key1", str, exists)
	strp, exists = nvs.PtrString("key1")
	t.Log("key1-p", strp, exists)

	// Get value of unknown key
	str, exists = nvs.String("unknown")
	t.Log("unknown", str, exists)
	strp, exists = nvs.PtrString("unknown")
	t.Log("unknown-p", strp, exists)

	// Get value as string
	strs = nvs.Strings("key2")
	t.Log("key2", strs)
}

func TestNVInt(t *testing.T) {
	var (
		val    int
		vals   []int
		valp   *int
		exists bool
	)
	mv := map[string]any{
		"key1a": "1028",
		"key1b": 1028,
	}
	nvs := NameValues{
		Pair: mv,
	}

	// Get value of key (string value)
	val, exists = nvs.Int("key1a")
	t.Log("key1a", val, exists)
	val, exists = nvs.Int("key1b")
	t.Log("key1b", val, exists)

	valp, exists = nvs.PtrInt("key1a")
	t.Log("key1a-p", valp, exists)
	valp, exists = nvs.PtrInt("key1b")
	t.Log("key1b-p", valp, exists)

	// Get value of unknown key
	val, exists = nvs.Int("unknown")
	t.Log("unknown", val, exists)
	valp, exists = nvs.PtrInt("unknown")
	t.Log("unknown-p", valp, exists)

	// Get value as string
	vals = nvs.Ints("key1a")
	t.Log("key1a-s", vals)
	vals = nvs.Ints("key1b")
	t.Log("key1b-s", vals)
}

func TestNVInt64(t *testing.T) {
	var (
		val    int64
		vals   []int64
		valp   *int64
		exists bool
	)
	mv := map[string]any{
		"key1a": "102810281028",
		"key1b": 102810281028,
	}
	nvs := NameValues{
		Pair: mv,
	}

	// Get value of key (string value)
	val, exists = nvs.Int64("key1a")
	t.Log("key1a", val, exists)
	val, exists = nvs.Int64("key1b")
	t.Log("key1b", val, exists)

	valp, exists = nvs.PtrInt64("key1a")
	t.Log("key1a-p", valp, exists)
	valp, exists = nvs.PtrInt64("key1b")
	t.Log("key1b-p", valp, exists)

	// Get value of unknown key
	val, exists = nvs.Int64("unknown")
	t.Log("unknown", val, exists)
	valp, exists = nvs.PtrInt64("unknown")
	t.Log("unknown-p", valp, exists)

	// Get value as string
	vals = nvs.Int64s("key1a")
	t.Log("key1a-s", vals)
	vals = nvs.Int64s("key1b")
	t.Log("key1b-s", vals)
}

func TestNVFloat64(t *testing.T) {
	var (
		val    float64
		vals   []float64
		valp   *float64
		exists bool
	)
	mv := map[string]any{
		"key1a": "1028.4321",
		"key1b": 1028.4321,
	}
	nvs := NameValues{
		Pair: mv,
	}

	// Get value of key (string value)
	val, exists = nvs.Float64("key1a")
	t.Log("key1a", val, exists)
	val, exists = nvs.Float64("key1b")
	t.Log("key1b", val, exists)

	valp, exists = nvs.PtrFloat64("key1a")
	t.Log("key1a-p", valp, exists)
	valp, exists = nvs.PtrFloat64("key1b")
	t.Log("key1b-p", valp, exists)

	// Get value of unknown key
	val, exists = nvs.Float64("unknown")
	t.Log("unknown", val, exists)
	valp, exists = nvs.PtrFloat64("unknown")
	t.Log("unknown-p", valp, exists)

	// Get value as string
	vals = nvs.Float64s("key1a")
	t.Log("key1a-s", vals)
	vals = nvs.Float64s("key1b")
	t.Log("key1b-s", vals)
}

func TestNVBool(t *testing.T) {
	var (
		val    bool
		vals   []bool
		valp   *bool
		exists bool
	)
	mv := map[string]any{
		"key1a": "true",
		"key1b": true,
		"key2a": "off",
		"key2b": "0",
	}
	nvs := NameValues{
		Pair: mv,
	}

	// Get value of key (string value)
	val, exists = nvs.Bool("key1a")
	t.Log("key1a", val, exists)
	val, exists = nvs.Bool("key1b")
	t.Log("key1b", val, exists)
	val, exists = nvs.Bool("key2a")
	t.Log("key2a", val, exists)
	val, exists = nvs.Bool("key2b")
	t.Log("key2b", val, exists)

	valp, exists = nvs.PtrBool("key1a")
	t.Log("key1a-p", valp, exists)
	valp, exists = nvs.PtrBool("key1b")
	t.Log("key1b-p", valp, exists)
	valp, exists = nvs.PtrBool("key2a")
	t.Log("key2a-p", valp, exists)
	valp, exists = nvs.PtrBool("key2b")
	t.Log("key2b-p", valp, exists)

	// Get value of unknown key
	val, exists = nvs.Bool("unknown")
	t.Log("unknown", val, exists)
	valp, exists = nvs.PtrBool("unknown")
	t.Log("unknown-p", valp, exists)

	// Get value as string
	vals = nvs.Bools("key1a")
	t.Log("key1a-s", vals)
	vals = nvs.Bools("key1b")
	t.Log("key1b-s", vals)
	vals = nvs.Bools("key2a")
	t.Log("key2a-s", vals)
	vals = nvs.Bools("key2b")
	t.Log("key2b-s", vals)
}

func TestNVDecimal(t *testing.T) {
	var (
		val    ssd.Decimal
		vals   []ssd.Decimal
		valp   *ssd.Decimal
		exists bool
	)
	dec, _ := ssd.NewFromString("10281028.4321")
	mv := map[string]any{
		"key1a": "10,281,028.4321",
		"key1b": 10281028.4321,
		"key1c": dec,
	}
	nvs := NameValues{
		Pair: mv,
	}

	// Get value of key (string value)
	val, exists = nvs.Decimal("key1a")
	t.Log("key1a", val, exists)
	val, exists = nvs.Decimal("key1b")
	t.Log("key1b", val, exists)
	val, exists = nvs.Decimal("key1c")
	t.Log("key1c", val, exists)

	valp, exists = nvs.PtrDecimal("key1a")
	t.Log("key1a-p", valp, exists)
	valp, exists = nvs.PtrDecimal("key1b")
	t.Log("key1b-p", valp, exists)
	valp, exists = nvs.PtrDecimal("key1c")
	t.Log("key1c-p", valp, exists)

	// Get value of unknown key
	val, exists = nvs.Decimal("unknown")
	t.Log("unknown", val, exists)
	valp, exists = nvs.PtrDecimal("unknown")
	t.Log("unknown-p", valp, exists)

	// Get value as string
	vals = nvs.Decimals("key1a")
	t.Log("key1a-s", vals)
	vals = nvs.Decimals("key1b")
	t.Log("key1b-s", vals)
	vals = nvs.Decimals("key1c")
	t.Log("key1b-c", vals)
}
