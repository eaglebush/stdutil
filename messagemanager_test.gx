package stdutil

import "testing"

func TestMessageManager(t *testing.T) {

	mm := &MessageManager{}

	// Add messages directly to the string array
	mm.Messages = append(mm.Messages, "   This is the first message not added thru any Add methods")
	mm.Messages = append(mm.Messages, "         This is the second message not added thru any Add methods")
	mm.Messages = append(mm.Messages, "This is the third  message not added thru any Add methods                ")

	// Add through methods
	mm.AddInfo("This is an information message!")
	mm.AddInfo("This is an information message too, damn!")
	mm.AddInfo("This is an information message too, damn you!")

	mm.AddWarning("This is a warning!")
	mm.AddWarning("This is a warning too!")
	mm.AddWarning("This is a warning too, damn!")

	mm.AddError("This is an error!")
	mm.AddError("This is an error too!")
	mm.AddError("This is an error too, damn!")

	mm.AddFatal("This is a fatal error!")
	mm.AddFatal("This is a fatal error too!")
	mm.AddFatal("This is a fatal error too, damn!")

	mm.AddAppMsg("This is an application message!")
	mm.AddAppMsg("This is an application message too, damn!")
	mm.AddAppMsg("This is an application message too, damn you!")

	for _, m := range mm.Messages {
		t.Log(`Unfixed`, m)
	}

	mm.Fix()
	for _, m := range mm.Messages {
		t.Log(`Fixed`, m)
	}

	t.Log(`Dominant Message`, mm.PrevailingType())
	t.Log(`Has Error Messages`, mm.HasErrors())
	t.Log(`Has Warning Messages`, mm.HasWarnings())
	t.Log(`Has Info Messages`, mm.HasInfos())

}
