package stdutil

import (
	"testing"
)

func TestResultMessage(t *testing.T) {

	r := InitResult()
	r.Messages = append(r.Messages, "   This is the first message")
	r.Messages = append(r.Messages, "         This is the second message")
	r.Messages = append(r.Messages, "This is the third  message                 ")

	if !r.IsStatusOK() {
		for _, m := range r.Messages {
			t.Log(m)
		}
	}

	r.AddInfo("     This is the fourth message      ")
	r.AddInfo("   This is the fifth message   ")
	r.AddInfo("  This is the sixth message  ")

	if !r.IsStatusOK() {
		for _, m := range r.Messages {
			t.Log(m)
		}
	}
}
