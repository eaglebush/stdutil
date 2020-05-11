package stdutil

import (
	"strings"
)

//Result - standard result structure
type Result struct {
	Execution    string   // Values: SUCCESS, FAIL
	Status       string   // OK, ERROR, VALID or any status
	FocusControl string   // Control to focus when error was activated
	Messages     []string // Messages in the result
	fixed        bool
}

// ResultMessageType - result message types
type ResultMessageType string

// Constants
const (
	INFO    ResultMessageType = "INFO: "
	WARNING ResultMessageType = "WARNING: "
	ERROR   ResultMessageType = "ERROR: "
)

// InitResult - initialize result for API query
func InitResult() Result {
	return Result{Execution: "FAIL", Status: "ERROR", Messages: make([]string, 0), fixed: false}
}

//Success - sets the execution status to SUCCESS
func (r *Result) Success() {
	r.Execution = "SUCCESS"
}

//Fail - sets the execution to FAIL
func (r *Result) Fail() {
	r.Execution = "FAIL"
}

//StatusOK - sets the Status to OK
func (r *Result) StatusOK() {
	r.Status = "OK"
}

//StatusError - sets the Status to Error
func (r *Result) StatusError() {
	r.Status = "ERROR"
}

//StatusValid - sets the Status to Valid
func (r *Result) StatusValid() {
	r.Status = "VALID"
}

//StatusInvalid - sets the Status to Invalid
func (r *Result) StatusInvalid() {
	r.Status = "INVALID"
}

//StatusYes - sets the Status to Yes
func (r *Result) StatusYes() {
	r.Status = "YES"
}

//StatusNo - sets the Status  to No
func (r *Result) StatusNo() {
	r.Status = "NO"
}

//IsStatusOK - checks if the status is OK
func (r *Result) IsStatusOK() bool {
	r.fixMessages()
	return r.Status == "OK"
}

//IsStatusError - checks if the status is Error
func (r *Result) IsStatusError() bool {
	r.fixMessages()
	return r.Status == "ERROR"
}

//IsStatusValid - checks if the status is Valid
func (r *Result) IsStatusValid() bool {
	r.fixMessages()
	return r.Status == "VALID"
}

//IsStatusInvalid - checks if the status is invalid
func (r *Result) IsStatusInvalid() bool {
	r.fixMessages()
	return r.Status == "INVALID"
}

//IsStatusYes - checks if the status is Yes
func (r *Result) IsStatusYes() bool {
	r.fixMessages()
	return r.Status == "YES"
}

//IsStatusNo - checks if the status is No
func (r *Result) IsStatusNo() bool {
	r.fixMessages()
	return r.Status == "NO"
}

// AddInfo - adds an information message
func (r *Result) AddInfo(Message string) {
	r.addMessage(strings.TrimSpace(Message), INFO)
}

// AddWarning - adds a warning message
func (r *Result) AddWarning(Message string) {
	r.addMessage(strings.TrimSpace(Message), WARNING)
}

// AddError - adds an error message
func (r *Result) AddError(Message string) {
	r.addMessage(strings.TrimSpace(Message), ERROR)
}

// HasErrors - Checks if the message array has errors
func (r Result) HasErrors() bool {

	for _, msg := range r.Messages {
		if strings.HasPrefix(strings.ToUpper(msg), string(ERROR)) {
			return true
		}
	}

	return false
}

// HasWarnings - Checks if the message array has warnings
func (r Result) HasWarnings() bool {

	for _, msg := range r.Messages {
		if strings.HasPrefix(strings.ToUpper(msg), string(WARNING)) {
			return true
		}
	}

	return false
}

// HasInfos - Checks if the message array has information messages
func (r Result) HasInfos() bool {

	for _, msg := range r.Messages {
		if strings.HasPrefix(strings.ToUpper(msg), string(INFO)) {
			return true
		}
	}

	return false
}

// DominantMessage - checks for a dominant message
func (r *Result) DominantMessage() ResultMessageType {
	r.fixMessages()

	nfo := 0
	wrn := 0
	err := 0

	for _, msg := range r.Messages {
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

// add new message to the message array
func (r *Result) addMessage(Message string, Type ResultMessageType) {
	sm := strings.ToUpper(Message)

	if !strings.HasPrefix(sm, string(Type)) {
		r.Messages = append(r.Messages, string(Type)+Message)
		return
	}

	r.Messages = append(r.Messages, Message)
}

// fix all messages that does not have this format
func (r *Result) fixMessages() {
	if r.fixed {
		return
	}

	for i, msg := range r.Messages {
		ms := strings.ToUpper(msg)
		switch true {
		case strings.HasPrefix(ms, string(INFO)):
		case strings.HasPrefix(ms, string(WARNING)):
		case strings.HasPrefix(ms, string(ERROR)):
		default:
			// fix all messages as errors
			r.Messages[i] = string(ERROR) + strings.TrimSpace(msg)
		}
	}

	r.fixed = true
}
