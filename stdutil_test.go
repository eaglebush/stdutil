package stdutil

import (
	"log"
	"testing"
	"time"
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
		NameValue{
			Name:  "prefix",
			Value: "SampleFunc",
		},
		NameValue{
			Name:  "message",
			Value: "This is a first message",
		},
		NameValue{
			Name:  "message",
			Value: "This is a second message",
		},
		NameValue{
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

	newString = NewString("NEW")
	newInt = NewInt(100)
	newFloat = NewFloat64(42.0)

	t.Log("With New:", newString, newInt, newFloat)

	t.Log("No new 1.8: ", newString18, newInt18, newFloat18)
	newString18 = New("NEW")
	newInt18 = New(100)
	newFloat18 = New(42.0)

	t.Log("With New 1.8:", newString18, newInt18, newFloat18)

}
