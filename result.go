package stdutil

import (
	"runtime"
	"strings"
)

// Status type
type Status string

// Status items
const (
	OK        Status = `OK`
	EXCEPTION Status = `EXCEPTION`
	VALID     Status = `VALID`
	INVALID   Status = `INVALID`
	YES       Status = `YES`
	NO        Status = `NO`
)

// Result - standard result structure
type Result struct {
	MessageManager
	Status       string       `json:"status,omitempty"`        // OK, ERROR, VALID or any status
	Operation    string       `json:"operation,omitempty"`     // Name of the operation / function that returned the result
	TaskID       *string      `json:"task_id,omitempty"`       // ID of the request and of the result
	WorkerID     *string      `json:"worker_id,omitempty"`     // ID of the worker that processed the data
	FocusControl *string      `json:"focus_control,omitempty"` // Control to focus when error was activated
	Page         *int         `json:"page,omitempty"`          // Current Page
	PageCount    *int         `json:"page_count,omitempty"`    // Page Count
	PageSize     *int         `json:"page_size,omitempty"`     // Page Size
	Tag          *interface{} `json:"tag,omitempty"`           // Miscellaneous result
}

// InitResult - initialize result for API query. This is the recommended initialization of this object.
func InitResult() Result {

	res := Result{
		Status: string(EXCEPTION),
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

// Return a status
func (r *Result) Return(status Status) {
	r.Status = string(status)
}

// OK returns true if the status is OK.
func (r *Result) OK() bool {
	return r.Status == string(OK)
}

// Error returns true if the status is EXCEPTION.
func (r *Result) Error() bool {
	return r.Status == string(EXCEPTION)
}

// Valid returns true if the status is VALID.
func (r *Result) Valid() bool {
	return r.Status == string(VALID)
}

// Invalid returns true if the status is INVALID.
func (r *Result) Invalid() bool {
	return r.Status == string(INVALID)
}

// Yes returns true if the status is YES.
func (r *Result) Yes() bool {
	return r.Status == string(YES)
}

// No returns true if the status is No.
func (r *Result) No() bool {
	return r.Status == string(NO)
}
