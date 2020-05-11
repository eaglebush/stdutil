package stdutil

import "strings"

// MessageType - result message types
type MessageType string

// Constants
const (
	INFO    MessageType = "INFO: "
	WARNING MessageType = "WARNING: "
	ERROR   MessageType = "ERROR: "
)

// MessageManager - a struct to create messages
type MessageManager struct {
	Messages []string
}

// AddInfo - adds an information message
func (r *MessageManager) AddInfo(Message string) {
	addMessage(&r.Messages, strings.TrimSpace(Message), INFO)
}

// AddWarning - adds a warning message
func (r *MessageManager) AddWarning(Message string) {
	addMessage(&r.Messages, strings.TrimSpace(Message), WARNING)
}

// AddError - adds an error message
func (r *MessageManager) AddError(Message string) {
	addMessage(&r.Messages, strings.TrimSpace(Message), ERROR)
}

// Fix - fix messages within an instance
func (r *MessageManager) Fix() {
	r.Messages = fixMessages(&r.Messages)
}

// HasErrors - Checks if the message array has errors
func (r MessageManager) HasErrors() bool {

	for _, msg := range r.Messages {
		if strings.HasPrefix(strings.ToUpper(msg), string(ERROR)) {
			return true
		}
	}

	return false
}

// HasWarnings - Checks if the message array has warnings
func (r MessageManager) HasWarnings() bool {

	for _, msg := range r.Messages {
		if strings.HasPrefix(strings.ToUpper(msg), string(WARNING)) {
			return true
		}
	}

	return false
}

// HasInfos - Checks if the message array has information messages
func (r MessageManager) HasInfos() bool {

	for _, msg := range r.Messages {
		if strings.HasPrefix(strings.ToUpper(msg), string(INFO)) {
			return true
		}
	}

	return false
}

// DominantMessageType - checks for a dominant message
func (r *MessageManager) DominantMessageType() MessageType {
	return getDominantMessageType(&r.Messages)
}

// AppendInfo - appends an information message
func AppendInfo(Messages *[]string, Message string) {
	addMessage(Messages, Message, INFO)
}

// AppendWarning - appends a warning message
func AppendWarning(Messages *[]string, Message string) {
	addMessage(Messages, Message, WARNING)
}

// AppendError - appends an error message
func AppendError(Messages *[]string, Message string) {
	addMessage(Messages, Message, ERROR)
}

// FixMessages - fix all unformatted messages to formatted messages
func FixMessages(Messages *[]string) {
	fixMessages(Messages)
}

// DominantMessageType - get dominant message type
func DominantMessageType(Messages *[]string) MessageType {
	return getDominantMessageType(Messages)
}

// fix messages
func fixMessages(Messages *[]string) []string {

	msgr := *Messages

	for i, msg := range *Messages {
		ms := strings.ToUpper(msg)
		switch true {
		case strings.HasPrefix(ms, string(INFO)):
		case strings.HasPrefix(ms, string(WARNING)):
		case strings.HasPrefix(ms, string(ERROR)):
		default:
			// fix all messages as errors
			msgr[i] = string(ERROR) + strings.TrimSpace(msg)
		}
	}

	return msgr
}

// add new message to the message array
func addMessage(Messages *[]string, Message string, Type MessageType) {

	Message = strings.TrimSpace(Message)
	sm := strings.ToUpper(Message)

	if !strings.HasPrefix(sm, string(Type)) {
		*Messages = append(*Messages, string(Type)+Message)
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

	for _, msg := range *Messages {
		switch true {
		case strings.HasPrefix(msg, string(INFO)):
			nfo++
		case strings.HasPrefix(msg, string(WARNING)):
			wrn++
		case strings.HasPrefix(msg, string(ERROR)):
			err++
		}
	}

	if nfo > wrn && nfo > err {
		return INFO
	}

	if wrn > nfo && wrn > err {
		return WARNING
	}

	if err > nfo && err > wrn {
		return ERROR
	}

	// default is error
	return ERROR
}
