package stdutil

import (
	"testing"

	"github.com/eaglebush/stdutil"
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
	str = stdutil.StripLeading(str, 37)
	t.Log(str, len(str))
}

func TestGenerateRandomString(t *testing.T) {
	for i := 0; i < 5; i++ {
		t.Logf("GenerateFull: %s", stdutil.GenerateFull(10))
		t.Logf("GenerateText: %s", stdutil.GenerateText(10))
		t.Logf("GenerateSeries: %s", stdutil.GenerateSeries(6))
		t.Logf("GenerateAlpha (Normal): %s", stdutil.GenerateAlpha(10, false))
		t.Logf("GenerateAlpha (Lower): %s", stdutil.GenerateAlpha(10, true))
		t.Logf("-------------------------------------------------------")
	}
	t.Logf("GenerateFull: %s", stdutil.GenerateFull(32))
}