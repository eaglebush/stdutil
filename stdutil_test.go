package stdutil

import (
	"log"
	"sync"
	"testing"
	"time"

	ssd "github.com/shopspring/decimal"
)

func TestNullOrEmpty(t *testing.T) {
	var (
		sam1 string
		sam2 int
		sam3 time.Time
		sam4 *string
		sam5 *time.Time
	)

	sam1 = "ok"
	sam2 = 0
	sam3 = time.Now()

	if IsNullOrEmpty(&sam1) {
		t.Log("String empty ")
	} else {
		t.Log("String Not empty")
	}

	if IsNullOrEmpty(&sam2) {
		t.Log("Int empty")
	} else {
		t.Log("Int Not empty")
	}

	if IsNullOrEmpty(&sam3) {
		t.Log("Time empty")
	} else {
		t.Log("Time Not empty")
	}

	if IsNullOrEmpty(sam4) {
		t.Log("String empty or null")
	} else {
		t.Log("String Not empty")
	}

	if IsNullOrEmpty(sam5) {
		t.Log("Time empty or null")
	} else {
		t.Log("Time Not empty")
	}

	if IsEmpty(&sam1) {
		t.Log("IsEmpty: String empty ")
	} else {
		t.Log("IsEmpty: String Not empty")
	}

	if IsEmpty(&sam2) {
		t.Log("IsEmpty: Int empty")
	} else {
		t.Log("IsEmpty: Int Not empty")
	}

	if IsEmpty(&sam3) {
		t.Log("IsEmpty: Time empty")
	} else {
		t.Log("IsEmpty: Time Not empty")
	}

	if IsEmpty(sam4) {
		t.Log("IsEmpty: String empty or null")
	} else {
		t.Log("IsEmpty: String Not empty")
	}

	if IsEmpty(sam5) {
		t.Log("IsEmpty: Time empty or null")
	} else {
		t.Log("IsEmpty: Time Not empty")
	}
}

func TestIsEmpty(t *testing.T) {
	var teststr *string
	if IsEmpty(teststr) {
		t.Log(`String is empty`)
	}

	teststr = new(string)
	if IsEmpty(teststr) {
		t.Log(`String is empty`)
	}

	*teststr = "Hi"
	if IsEmpty(teststr) {
		t.Log(`String is empty`)
	}
}

func TestStripEndingForwardSlash(t *testing.T) {
	addr := "http://localhost:8000asdsadas/"
	addr = StripEndingForwardSlash(addr)
	t.Log(addr)
	t.Fail()
}

func TestStripTrailing(t *testing.T) {
	str := "First_official_images_from_Adam_Wingard's_'Godzilla_vs._Kong'"
	str = StripTrailing(str, 37)
	t.Log(str, len(str))
}

func TestStripLeading(t *testing.T) {
	str := "First_official_images_from_Adam_Wingard's_'Godzilla_vs._Kong'"
	str = StripLeading(str, 37)
	t.Log(str, len(str))
}

func TestGenerateRandomString(t *testing.T) {
	for i := 0; i < 5; i++ {
		t.Logf("GenerateFull: %s", GenerateFull(10))
		t.Logf("GenerateText: %s", GenerateText(10))
		t.Logf("GenerateSeries: %s", GenerateSeries(6))
		t.Logf("GenerateAlpha (Normal): %s", GenerateAlpha(10, false))
		t.Logf("GenerateAlpha (Lower): %s", GenerateAlpha(10, true))
		t.Logf("-------------------------------------------------------")
	}
	t.Logf("GenerateFull: %s", GenerateFull(32))
}

func TestResult(t *testing.T) {

	res := InitResult(
		NameValue[string]{
			Name:  "prefix",
			Value: "SampleFunc",
		},
		NameValue[string]{
			Name:  "message",
			Value: "This is a first message",
		},
		NameValue[string]{
			Name:  "message",
			Value: "This is a second message",
		},
		NameValue[string]{
			Name:  "message",
			Value: "This is a third message",
		},
	)

	res.MessagePrefix = "WEH"
	res.AddError("This is an error message")
	res.AddInfo("This is an informational message")

	log.Println(res.MessagesToString())

}

func TestNew(t *testing.T) {

	var (
		newString   *string
		newInt      *int
		newFloat    *float64
		newString18 *string
		newInt18    *int
		newFloat18  *float64
	)

	t.Log("No new: ", newString, newInt, newFloat)

	t.Log("No new 1.8: ", newString18, newInt18, newFloat18)
	newString18 = New("NEW")
	newInt18 = New(100)
	newFloat18 = New(42.0)

	t.Log("With New 1.8:", newString18, newInt18, newFloat18)

}

func TestInterpolate(t *testing.T) {
	str, obj := Interpolate(`This is ${name}. Leader of the ${band} band.`, NameValues{
		Pair: map[string]any{
			"name": "Zaldy",
			"band": "Razzie",
		},
	})

	log.Println(str, obj)
}

func TestIn(t *testing.T) {

	type SP string

	const spA SP = "A"
	const spB SP = "B"
	const spC SP = "C"
	const spX SP = "X"

	seek := spX

	if !In(seek, spA, spB, spC) {
		log.Println("Seek parameter does not exist in the variadic parameter")
	}
}

func TestValidateDecimal(t *testing.T) {

	vo :=
		DecimalValidationOptions{
			Null:  true,
			Empty: true,
			Min:   New(ssd.NewFromFloat32(30.49)),
			Max:   New(ssd.NewFromFloat32(30.51)),
		}

	//d := ssd.NewFromFloat(30.48)
	//var d ssd.Decimal

	if err := ValidateDecimal(nil, &vo); err != nil {
		t.Fatal(err)
	}

}

func TestValidateNumeric(t *testing.T) {

	vo :=
		NumericValidationOptions[int]{
			Null:  false,
			Empty: false,
			Min:   30,
			Max:   50,
		}

	d := 0
	//var d ssd.Decimal

	if err := ValidateNumeric(&d, &vo); err != nil {
		t.Fatal(err)
	}

}

func TestBuildSeries(t *testing.T) {
	series := BuildSeries(100, SeriesOptions{
		Prefix: "PF",
		Suffix: "EX",
		Length: 0,
	})
	t.Log(series)
}

func TestInterfaceArray(t *testing.T) {
	arr := ToInterfaceArray(time.Now())
	t.Log(arr)
}

func TestGetElement(t *testing.T) {

	arrs := []string{
		"Aruba",
		"Jamaica",
		"Bahamas",
	}

	var exists bool

	str := Elem(&arrs, 4, &exists)
	t.Logf(`Value: %s, Exists: %t`, str, exists)

	str = Elem(&arrs, 2, &exists)
	t.Logf(`Value: %s, Exists: %t`, str, exists)

	strp := ElemPtr(&arrs, 1, &exists)
	t.Logf(`Value: %s, Exists: %t`, *strp, exists)
}

func TestNull(t *testing.T) {
	// Non pointer string
	var a any
	value := Null[string](a, "actual")
	t.Logf(`Value: %s`, value)

	// Pointer string
	var valstr string
	b := new(string)
	*b = "test"
	valstr = Null[string](b, "actual")
	t.Logf(`Value: %p`, &valstr)
}

func TestNonNullComp(t *testing.T) {
	var (
		p1  *string
		p2  *string
		res int
	)

	p2 = new(string)
	res = NonNullComp(p1, p2)
	if res == -1 {
		t.Log(`One of the parameters is invalid`)
	} else {
		if res == 0 {
			t.Log(`Parameters are equal`)
		} else if res == 1 {
			t.Log(`Parameters are not equal`)
		}
	}

	p1 = new(string)
	res = NonNullComp(p1, p2)
	if res == -1 {
		t.Log(`One of the parameters is invalid`)
	} else {
		if res == 0 {
			t.Log(`Parameters are equal`)
		} else if res == 1 {
			t.Log(`Parameters are not equal`)
		}
	}

	*p2 = "Hi"
	res = NonNullComp(p1, p2)
	if res == -1 {
		t.Log(`One of the parameters is invalid`)
	} else {
		if res == 0 {
			t.Log(`Parameters are equal`)
		} else if res == 1 {
			t.Log(`Parameters are not equal`)
		}
	}

	*p1 = "Hi"
	res = NonNullComp(p1, p2)
	if res == -1 {
		t.Log(`One of the parameters is invalid`)
	} else {
		if res == 0 {
			t.Log(`Parameters are equal`)
		} else if res == 1 {
			t.Log(`Parameters are not equal`)
		}
	}
}

func TestSafeMapWrite(t *testing.T) {
	rw := &sync.RWMutex{}
	m := map[string]int{}
	for i := 0; i < 1000; i++ {
		k := i
		go func() {
			SafeMapWrite(&m, "testing", k, rw)
			read := SafeMapRead(&m, "testing", rw)
			t.Logf("Reading map: %d", read)
		}()
	}
}
