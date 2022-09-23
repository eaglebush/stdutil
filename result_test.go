package stdutil

import (
	"errors"
	"testing"
)

func TestResultMessage(t *testing.T) {

	r := InitResult()
	r.Messages = append(r.Messages, "   This is the first message not added thru any Add methods")
	r.Messages = append(r.Messages, "         This is the second message not added thru any Add methods")
	r.Messages = append(r.Messages, "This is the third  message not added thru any Add methods                ")

	r.Status = string(MsgWarn)

	r.AddInfo("This is an information message too!")
	r.AddInfo("This is an information message too, damn!")
	x := r.AddInfo("This is an information message too, damn you!")

	r.AddWarning("This is a warning!")

	for _, m := range r.Messages {
		t.Log(`Unfixed`, m)
	}

	// Result returned from r.AddInfo
	for _, m := range x.Messages {
		t.Logf("Result returned: %s, Status: %s", m, x.Status)
	}

	mm := r.MessageManager()
	mm.Fix()

	// AppendError(&mm.Messages, "This is an appended message")

	if !r.OK() {
		for _, m := range r.Messages {
			t.Log(`Fixed`, m)
		}
	}

	t.Log(`Dominant Message`, mm.PrevailingType())
	t.Log(`Has Error Messages`, mm.HasErrors())
	t.Log(`Has Warning Messages`, mm.HasWarnings())
	t.Log(`Has Info Messages`, mm.HasInfos())

}

func TestInitResult(t *testing.T) {
	res := InitResult(NameValue[string]{
		Name:  "prefix",
		Value: "stdutil",
	})

	res.AddError("This is an error")
	res.AddError("This is another error")
	res.AddError("This is another significant error")

	err := errors.New("This is an err")
	res.AddErr(err)

	for _, m := range res.Messages {
		t.Log(m)
	}

}
