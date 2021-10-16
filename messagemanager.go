package stdutil

import "strings"

// MessageType - result message types
type MessageType string

// Constants
const (
	MsgInfo      MessageType = "INFO"
	MsgWarn      MessageType = "WARNING"
	MsgError     MessageType = "ERROR"
	MsgFatal     MessageType = "FATAL"
	MsgApp       MessageType = ""
	DelimMsgType string      = `: `
)

// MessageManager - a struct to create messages
type MessageManager struct {
	Messages []string `json:"messages,omitempty"`
}

// AddInfo - adds an information message
func (r *MessageManager) AddInfo(Message ...string) {
	for _, m := range Message {
		addMessage(&r.Messages, m, MsgInfo)
	}
}

// AddWarning - adds a warning message
func (r *MessageManager) AddWarning(Message ...string) {
	for _, m := range Message {
		addMessage(&r.Messages, m, MsgWarn)
	}
}

// AddError - adds an error message
func (r *MessageManager) AddError(Message ...string) {
	for _, m := range Message {
		addMessage(&r.Messages, m, MsgError)
	}
}

// AddFatal - adds a fatal error message
func (r *MessageManager) AddFatal(Message ...string) {
	for _, m := range Message {
		addMessage(&r.Messages, m, MsgFatal)
	}
}

// AddAppMsg - adds an error message
func (r *MessageManager) AddAppMsg(Message ...string) {
	for _, m := range Message {
		addMessage(&r.Messages, m, MsgApp)
	}
}

// Fix - fix messages within an instance
func (r *MessageManager) Fix() {
	r.Messages = fixMessages(&r.Messages)
}

// HasErrors - Checks if the message array has errors
func (r MessageManager) HasErrors() bool {
	for _, msg := range r.Messages {
		if strings.HasPrefix(strings.ToUpper(msg), string(MsgError)+DelimMsgType) {
			return true
		}
	}
	return false
}

// HasWarnings - Checks if the message array has warnings
func (r MessageManager) HasWarnings() bool {

	for _, msg := range r.Messages {
		if strings.HasPrefix(strings.ToUpper(msg), string(MsgWarn)+DelimMsgType) {
			return true
		}
	}

	return false
}

// HasInfos - Checks if the message array has information messages
func (r MessageManager) HasInfos() bool {
	for _, msg := range r.Messages {
		if strings.HasPrefix(strings.ToUpper(msg), string(MsgInfo)+DelimMsgType) {
			return true
		}
	}
	return false
}

// PrevailingType - checks for a dominant message
func (r *MessageManager) PrevailingType() MessageType {
	return getDominantMessageType(&r.Messages)
}

// ToString return the messages as a carriage/return delimited string
func (r *MessageManager) ToString() string {
	return strings.Join(r.Messages, "\r\n")
}

// AppendInfo - appends an information message
func AppendInfo(Messages *[]string, Message ...string) {
	for _, m := range Message {
		addMessage(Messages, m, MsgInfo)
	}
}

// AppendWarning - appends a warning message
func AppendWarning(Messages *[]string, Message ...string) {
	for _, m := range Message {
		addMessage(Messages, m, MsgWarn)
	}
}

// AppendError - appends an error message
func AppendError(Messages *[]string, Message ...string) {
	for _, m := range Message {
		addMessage(Messages, m, MsgError)
	}
}

// FixMessages - fix all unformatted messages to formatted messages
func FixMessages(Messages *[]string) {
	fixMessages(Messages)
}

// DominantMessageType - get dominant message type. App messages will be deleted
func DominantMessageType(Messages *[]string) MessageType {
	return getDominantMessageType(Messages)
}

// trims message from spaces
func fixMessages(Messages *[]string) []string {

	msgr := *Messages

	for i, msg := range *Messages {
		msgr[i] = strings.TrimSpace(msg)
	}

	return msgr
}

// add new message to the message array
func addMessage(Messages *[]string, Message string, Type MessageType) {

	Message = strings.TrimSpace(Message)
	td := string(Type) + DelimMsgType

	if !strings.HasPrefix(strings.ToUpper(Message), td) && Type != MsgApp {
		*Messages = append(*Messages, td+Message)
		return
	}

	*Messages = append(*Messages, Message)
}

// get dominant message
func getDominantMessageType(Messages *[]string) MessageType {
	fixMessages(Messages)

	nfo := 0
	wrn := 0
	err := 0
	ftl := 0

	for _, msg := range *Messages {
		switch true {
		case strings.HasPrefix(msg, string(MsgInfo)+DelimMsgType):
			nfo++
		case strings.HasPrefix(msg, string(MsgWarn)+DelimMsgType):
			wrn++
		case strings.HasPrefix(msg, string(MsgError)+DelimMsgType):
			err++
		case strings.HasPrefix(msg, string(MsgFatal)+DelimMsgType):
			ftl++
		}
	}

	// fatal errors always dominate
	if ftl > 0 {
		return MsgFatal
	}

	if nfo > wrn && nfo > err {
		return MsgInfo
	}

	if wrn > nfo && wrn > err {
		return MsgWarn
	}

	if err > nfo && err > wrn {
		return MsgError
	}

	return MsgApp
}
