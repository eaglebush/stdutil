package stdutil

import (
	"testing"
)

func TestResultMessage(t *testing.T) {

	r := InitResult()
	r.Messages = append(r.Messages, "   This is the first message")
	r.Messages = append(r.Messages, "         This is the second message")
	r.Messages = append(r.Messages, "This is the third  message                 ")
	AppendInfo(&r.Messages, "This is an information message!")
	r.AddInfo("This is an information message too!")
	r.AddInfo("This is an information message too, damn!")
	r.AddInfo("This is an information message too, damn you!")
	//r.AddWarning("This is a warning!")

	for _, m := range r.Messages {
		t.Log(`Unfixed`, m)
	}

	//r.Fix()

	if !r.OK() {
		for _, m := range r.Messages {
			t.Log(`Fixed`, m)
		}
	}

	t.Log(`Dominant Message`, r.DominantMessageType())
	t.Log(`Has Error Messages`, r.HasErrors())
	t.Log(`Has Warning Messages`, r.HasWarnings())
	t.Log(`Has Info Messages`, r.HasInfos())

}
