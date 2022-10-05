package stdutil

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/narsilworks/livenote"
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

	ln        livenote.LiveNote // Internal note
	eventVerb string            // event verb related to the name of the operation
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
func InitResult(args ...NameValue[string]) Result {

	res := Result{
		Status: string(EXCEPTION),
		ln:     livenote.LiveNote{},
	}

	res.Messages = make([]string, 0)

	for _, nv := range args {

		// check if it is a valid status, ignore if not
		// go to next value if valid
		if strings.EqualFold(nv.Name, `status`) {
			switch nv.Value {
			case string(OK), string(EXCEPTION), string(VALID), string(INVALID), string(YES), string(NO):
				res.Status = nv.Value
				continue
			}
		}

		if strings.EqualFold(nv.Name, `prefix`) {
			res.MessagePrefix = nv.Value
			res.ln.Prefix = nv.Value // set default prefix for livenote
			continue
		}

		if strings.EqualFold(nv.Name, `message`) {
			if res.Status == string(EXCEPTION) {
				res.AddError(nv.Value)
			} else {
				res.AddInfo(nv.Value)
			}
			continue
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
			res.eventVerb = res.Operation
		}
	}

	return res
}

// MessageManager returns the internal message manager
func (r *Result) MessageManager() *livenote.LiveNote {
	return &r.ln
}

// Return a status
func (r *Result) Return(status Status) Result {
	r.Status = string(status)
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

// AddInfo adds an information message and returns itself
func (r *Result) AddInfo(message ...string) Result {
	// add message
	r.ln.AddInfo(message...)

	// get current notes to update the messages
	nts := r.ln.Notes()
	r.Messages = make([]string, 0, len(nts))
	for _, n := range nts {
		r.Messages = append(r.Messages, n.ToString())
	}
	return *r
}

// AddInfof adds a formatted information message and returns itself
func (r *Result) AddInfof(format string, a ...interface{}) Result {
	return r.AddInfo(fmt.Sprintf(format, a...))
}

// AddWarning - adds a warning message and returns itself
func (r *Result) AddWarning(message ...string) Result {
	// add message
	r.ln.AddWarning(message...)

	// get current notes to update the messages
	nts := r.ln.Notes()
	r.Messages = make([]string, 0, len(nts))
	for _, n := range nts {
		r.Messages = append(r.Messages, n.ToString())
	}
	return *r
}

// AddWarningf adds a formatted warning message and returns itself
func (r *Result) AddWarningf(format string, a ...interface{}) Result {
	return r.AddWarning(fmt.Sprintf(format, a...))
}

// AddError adds an error message and returns itself
func (r *Result) AddError(message ...string) Result {
	// add message
	r.ln.AddError(message...)

	// get current notes to update the messages
	nts := r.ln.Notes()
	r.Messages = make([]string, 0, len(nts))
	for _, n := range nts {
		r.Messages = append(r.Messages, n.ToString())
	}
	return *r
}

// AddErrorf adds a formatted error message and returns itself
func (r *Result) AddErrorf(format string, a ...interface{}) Result {
	return r.AddError(fmt.Sprintf(format, a...))
}

// AddErr - adds a real error and returns itself
func (r *Result) AddErr(err error) Result {
	r.AddError(err.Error())
	return *r
}

// AppendError copies the messages of the Result parameter and append the current message
func (r *Result) AppendError(rs Result, message ...string) Result {

	for _, n := range rs.ln.Notes() {
		r.ln.Append(n)
	}

	if len(message) == 0 {
		r.AddError(message...)
	}

	return *r
}

// AppendErr copies the messages of the Result parameter and append an error message
func (r *Result) AppendErr(rs Result, err error) Result {

	for _, n := range rs.ln.Notes() {
		r.ln.Append(n)
	}

	return r.AddErr(err)
}

// AppendErrorf copies the messages of the Result parameter and append a formatted error message
func (r *Result) AppendErrorf(rs Result, format string, a ...interface{}) Result {

	for _, n := range rs.ln.Notes() {
		r.ln.Append(n)
	}

	return r.AddErrorf(format, a...)
}

// AppendInfo copies the messages of the Result parameter and append the current message
func (r *Result) AppendInfo(rs Result, message ...string) Result {

	for _, n := range rs.ln.Notes() {
		r.ln.Append(n)
	}

	if len(message) == 0 {
		r.AddInfo(message...)
	}

	return *r
}

// AppendInfof copies the messages of the Result parameter and append a formatted information message
func (r *Result) AppendInfof(rs Result, format string, a ...interface{}) Result {

	for _, n := range rs.ln.Notes() {
		r.ln.Append(n)
	}

	return r.AddInfof(format, a...)
}

// AppendWarning copies the messages of the Result parameter and append the current message
func (r *Result) AppendWarning(rs Result, message ...string) Result {

	for _, n := range rs.ln.Notes() {
		r.ln.Append(n)
	}

	if len(message) == 0 {
		r.AddWarning(message...)
	}

	return *r
}

// AppendWarningf copies the messages of the Result parameter and append a formatted warning message
func (r *Result) AppendWarningf(rs Result, format string, a ...interface{}) Result {

	for _, n := range rs.ln.Notes() {
		r.ln.Append(n)
	}

	return r.AddWarningf(format, a...)
}

// EventID returns the past tense of Operation
func (r *Result) EventID() string {
	ev := r.eventVerb
	if ev == "" {
		return "unknown"
	}

	// simple past tenser
	if !strings.HasSuffix(ev, "e") {
		return ev + "ed"
	}

	return ev + "d"
}

// ToString adds a formatted error message and returns itself
func (r *Result) MessagesToString() string {
	return r.ln.ToString()
}

// SetPrefix changes the prefix
func (r *Result) SetPrefix(pfx string) {
	r.ln.Prefix = pfx
	r.MessagePrefix = pfx
}

// RowsAffectedInfo - a function to simplify adding information for rows affected
func (r *Result) RowsAffectedInfo(rowsaff int64) {
	if rowsaff != 0 {
		r.AddInfo(fmt.Sprintf("%d rows affected", rowsaff))
	} else {
		r.AddInfo("No rows affected")
	}
}
