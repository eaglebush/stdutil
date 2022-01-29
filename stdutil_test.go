package stdutil

import (
	"log"
	"testing"
)

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
