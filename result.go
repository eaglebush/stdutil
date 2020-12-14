package stdutil

import (
	"runtime"
	"strings"
)

//Result - standard result structure
type Result struct {
	MessageManager
	TaskID       *string      `json:"task_id,omitempty"`       // ID of the request and of the result
	WorkerID     *string      `json:"worker_id,omitempty"`     // ID of the worker that processed the data
	Execution    string       `json:"execution,omitempty"`     // Values: SUCCESS, FAIL
	Status       string       `json:"status,omitempty"`        // OK, ERROR, VALID or any status
	Operation    string       `json:"operation,omitempty"`     // Name of the operation / function that returned the result
	FocusControl string       `json:"focus_control,omitempty"` // Control to focus when error was activated
	Page         *int         `json:"page,omitempty"`          // Current Page
	PageCount    *int         `json:"page_count,omitempty"`    // Page Count
	PageSize     *int         `json:"page_size,omitempty"`     // Page Size
	Tag          *interface{} `json:"tag,omitempty"`           // Miscellaneous result
}

// InitResult - initialize result for API query. This is the recommended initialization of this object.
func InitResult() Result {

	res := Result{
		Execution: "FAIL",
		Status:    "ERROR",
	}
	res.Messages = make([]string, 0)

	// Auto-detect function that called this function
	if pc, _, _, ok := runtime.Caller(1); ok {
		if details := runtime.FuncForPC(pc); details != nil {
			nm := details.Name()
			if pos := strings.LastIndex(nm, `.`); pos != -1 {
				nm = nm[pos+1:]
			}
			res.Operation = strings.ToLower(nm)
		}
	}

	return res
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
	FixMessages(&r.Messages)
	return r.Status == "OK"
}

//IsStatusError - checks if the status is Error
func (r *Result) IsStatusError() bool {
	FixMessages(&r.Messages)
	return r.Status == "ERROR"
}

//IsStatusValid - checks if the status is Valid
func (r *Result) IsStatusValid() bool {
	FixMessages(&r.Messages)
	return r.Status == "VALID"
}

//IsStatusInvalid - checks if the status is invalid
func (r *Result) IsStatusInvalid() bool {
	FixMessages(&r.Messages)
	return r.Status == "INVALID"
}

//IsStatusYes - checks if the status is Yes
func (r *Result) IsStatusYes() bool {
	FixMessages(&r.Messages)
	return r.Status == "YES"
}

//IsStatusNo - checks if the status is No
func (r *Result) IsStatusNo() bool {
	FixMessages(&r.Messages)
	return r.Status == "NO"
}
