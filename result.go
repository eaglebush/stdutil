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
	Messages     []string     `json:"messages,omitempty"`      // Accumulated messages as a result from Add methods. Do not append messages using append()
	Status       string       `json:"status,omitempty"`        // OK, ERROR, VALID or any status
	Operation    string       `json:"operation,omitempty"`     // Name of the operation / function that returned the result
	TaskID       *string      `json:"task_id,omitempty"`       // ID of the request and of the result
	WorkerID     *string      `json:"worker_id,omitempty"`     // ID of the worker that processed the data
	FocusControl *string      `json:"focus_control,omitempty"` // Control to focus when error was activated
	Page         *int         `json:"page,omitempty"`          // Current Page
	PageCount    *int         `json:"page_count,omitempty"`    // Page Count
	PageSize     *int         `json:"page_size,omitempty"`     // Page Size
	Tag          *interface{} `json:"tag,omitempty"`           // Miscellaneous result
	mm           *MessageManager
}

// InitResult - initialize result for API query. This is the recommended initialization of this object.
// In the variadic argument, the first slice will be its status, the rest will be added to the messages
func InitResult(args ...string) Result {

	res := Result{
		Status: string(EXCEPTION),
	}

	res.mm = &MessageManager{}

	res.Messages = make([]string, 0)

	if ln := len(args); ln > 0 {

		res.Status = args[0]

		if ln > 1 {
			for i := 1; i < len(args); i++ {
				if res.Status == string(EXCEPTION) {
					res.AddError(args[i])
				} else {
					res.AddInfo(args[i])
				}
			}
		}
	}

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
func (r *Result) Return(status Status) Result {
	r.Status = string(status)

	if r.mm == nil {
		r.mm = &MessageManager{}
	}

	r.Messages = r.mm.Messages

	return *r
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

// MessageManager returns the internal message manager for further manipulation
func (r *Result) MessageManager() *MessageManager {
	if r.mm == nil {
		r.mm = &MessageManager{}
	}
	return r.mm
}

// AddInfo - adds an information message and returns itself
func (r *Result) AddInfo(Message ...string) Result {
	if r.mm == nil {
		r.mm = &MessageManager{}
	}
	r.mm.AddInfo(Message...)
	r.Messages = r.mm.Messages

	return *r
}

// AddWarning - adds a warning message and returns itself
func (r *Result) AddWarning(Message ...string) Result {
	if r.mm == nil {
		r.mm = &MessageManager{}
	}
	r.mm.AddWarning(Message...)
	r.Messages = r.mm.Messages

	return *r
}

// AddError - adds an error message and returns itself
func (r *Result) AddError(Message ...string) Result {
	if r.mm == nil {
		r.mm = &MessageManager{}
	}
	r.mm.AddError(Message...)
	r.Messages = r.mm.Messages

	return *r
}

// AddErr - adds a real error and returns itself
func (r *Result) AddErr(err error) Result {
	if r.mm == nil {
		r.mm = &MessageManager{}
	}
	r.mm.AddError(err.Error())
	r.Messages = r.mm.Messages

	return *r
}
