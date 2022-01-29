package stdutil

import (
	"fmt"
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
	Messages      []string     `json:"messages"`                // Accumulated messages as a result from Add methods. Do not append messages using append()
	Status        string       `json:"status"`                  // OK, ERROR, VALID or any status
	Operation     string       `json:"operation,omitempty"`     // Name of the operation / function that returned the result
	TaskID        *string      `json:"task_id,omitempty"`       // ID of the request and of the result
	WorkerID      *string      `json:"worker_id,omitempty"`     // ID of the worker that processed the data
	FocusControl  *string      `json:"focus_control,omitempty"` // Control to focus when error was activated
	Page          *int         `json:"page,omitempty"`          // Current Page
	PageCount     *int         `json:"page_count,omitempty"`    // Page Count
	PageSize      *int         `json:"page_size,omitempty"`     // Page Size
	Tag           *interface{} `json:"tag,omitempty"`           // Miscellaneous result
	MessagePrefix string       `json:"prefix,omitempty"`        // Prefix of the message to return
	mm            *MessageManager
}

// InitResult - initialize result for API query. This is the recommended initialization of this object.
// The variadic arguments of std.NameValue data type will be intepreted as follows:
//
// To set the initial status, set NameValue.Name to "status" and set NameValue.Value to a valid status
// if the value is not a valid status, it will be ignored */
//
// To set a message prefix, set NameValue.Name to "prefix" and set NameValue.Value to a string value.
//
// To add a message, set NameValue.Name to "message" and set NameValue.Value to a valid message.
// Depending on the current status (default is EXCEPTION), the message type is automatically set to that type
func InitResult(args ...NameValue) Result {

	res := Result{
		Status: string(EXCEPTION),
	}

	res.mm = &MessageManager{}

	res.Messages = make([]string, 0)

	if ln := len(args); ln > 0 {

		for i := 0; i < ln; i++ {

			nm := strings.TrimSpace(args[i].Name)

			// check if it is a valid status, ignore if not
			// go to next value if valid
			if strings.EqualFold(nm, `status`) {
				nvs, ok := args[i].Value.(string)
				if ok {
					switch nvs {
					case string(OK), string(EXCEPTION), string(VALID), string(INVALID), string(YES), string(NO):
						res.Status = nvs
						continue
					}
				}
			}

			if strings.EqualFold(nm, `prefix`) {
				nvs, ok := args[i].Value.(string)
				if ok {
					res.mm.MessagePrefix = nvs
					res.MessagePrefix = nvs
					continue
				}
			}

			if strings.EqualFold(nm, `message`) {
				nvs, ok := args[i].Value.(string)
				if ok {
					if res.Status == string(EXCEPTION) {
						res.AddError(nvs)
					} else {
						res.AddInfo(nvs)
					}
					continue
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
		r.mm = &MessageManager{
			MessagePrefix: r.MessagePrefix,
		}
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
		r.mm = &MessageManager{
			MessagePrefix: r.MessagePrefix,
		}
	}
	return r.mm
}

// AddInfo - adds an information message and returns itself
func (r *Result) AddInfo(message ...string) Result {
	if r.mm == nil {
		r.mm = &MessageManager{
			MessagePrefix: r.MessagePrefix,
		}
	}

	r.mm.MessagePrefix = r.MessagePrefix
	r.mm.AddInfo(message...)
	r.Messages = r.mm.Messages

	return *r
}

// AddInfof adds a formatted information message and returns itself
func (r *Result) AddInfof(format string, a ...interface{}) Result {
	return r.AddInfo(fmt.Sprintf(format, a...))
}

// AddWarning - adds a warning message and returns itself
func (r *Result) AddWarning(message ...string) Result {
	if r.mm == nil {
		r.mm = &MessageManager{
			MessagePrefix: r.MessagePrefix,
		}
	}

	r.mm.MessagePrefix = r.MessagePrefix
	r.mm.AddWarning(message...)
	r.Messages = r.mm.Messages

	return *r
}

// AddWarningf adds a formatted warning message and returns itself
func (r *Result) AddWarningf(format string, a ...interface{}) Result {
	return r.AddWarning(fmt.Sprintf(format, a...))
}

// AddError - adds an error message and returns itself
func (r *Result) AddError(message ...string) Result {

	if r.mm == nil {
		r.mm = &MessageManager{
			MessagePrefix: r.MessagePrefix,
		}
	}

	r.mm.MessagePrefix = r.MessagePrefix
	r.mm.AddError(message...)
	r.Messages = r.mm.Messages

	return *r
}

// AddErrorf adds a formatted error message and returns itself
func (r *Result) AddErrorf(format string, a ...interface{}) Result {
	return r.AddError(fmt.Sprintf(format, a...))
}

// AddErr - adds a real error and returns itself
func (r *Result) AddErr(err error) Result {

	if r.mm == nil {
		r.mm = &MessageManager{
			MessagePrefix: r.MessagePrefix,
		}
	}

	r.mm.MessagePrefix = r.MessagePrefix
	r.mm.AddError(err.Error())
	r.Messages = r.mm.Messages

	return *r
}

// ToString adds a formatted error message and returns itself
func (r *Result) MessagesToString() string {
	return r.mm.ToString()
}
