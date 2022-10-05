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

	// AppendError(&mm.Messages, "This is an appended message")

	if !r.OK() {
		for _, m := range r.Messages {
			t.Log(`Fixed`, m)
		}
	}

	t.Log(`Dominant Message`, mm.Prevailing())
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

func TestMessageMix(t *testing.T) {

	res1 := InitResult(NameValue[string]{
		Name:  "prefix",
		Value: "WHERE",
	})

	res2 := InitResult(NameValue[string]{
		Name:  "prefix",
		Value: "WHAT",
	})

	res3 := InitResult(NameValue[string]{
		Name:  "prefix",
		Value: "WHO",
	})

	res4 := InitResult(NameValue[string]{
		Name:  "prefix",
		Value: "HOW",
	})

	res := InitResult(NameValue[string]{
		Name:  "prefix",
		Value: "NEWS",
	})

	res1.AddInfo("This is an info")
	res1.AddWarning("This is a warning")
	res1.AddError("This is an error")

	res1.SetPrefix("WHERE2")

	res1.AddInfo("This is an info")
	res1.AddWarning("This is a warning")
	res1.AddError("This is an error")

	res2.AddInfo("This is an info")
	res2.AddWarning("This is a warning")
	res2.AddError("This is an error")

	res3.AddInfo("This is an info")
	res3.AddWarning("This is a warning")
	res3.AddError("This is an error")

	res4.AddInfo("This is an info")
	res4.AddWarning("This is a warning")
	res4.AddError("This is an error")

	res.AppendInfo(res1, "This is a new info, appending to WHERE group")
	res.AppendInfo(res2, "This is a new info, appending to WHAT group")
	res.AppendInfo(res3, "This is a new info, appending to HOW group")
	res.AppendInfo(res4, "This is a new info, appending to WHO group")

	res.AppendWarning(res1, "This is a new warning, appending to WHERE group")
	res.AppendWarning(res2, "This is a new warning, appending to WHAT group")
	res.AppendWarning(res3, "This is a new warning, appending to HOW group")
	res.AppendWarning(res4, "This is a new warning, appending to WHO group")

	res.AppendError(res1, "This is a new error, appending to WHERE group")
	res.AppendError(res2, "This is a new error, appending to WHAT group")
	res.AppendError(res3, "This is a new error, appending to HOW group")
	res.AppendError(res4, "This is a new error, appending to WHO group")

	for _, m := range res.Messages {
		t.Log(m)
	}

}
