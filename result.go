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
	TaskID        *string      `json:"task_id,omitempty"`       // ID of the task and of the result
	WorkerID      *string      `json:"worker_id,omitempty"`     // ID of the worker that processed the data
	FocusControl  *string      `json:"focus_control,omitempty"` // Control to focus when error was activated
	Page          *int         `json:"page,omitempty"`          // Current Page
	PageCount     *int         `json:"page_count,omitempty"`    // Page Count
	PageSize      *int         `json:"page_size,omitempty"`     // Page Size
	Tag           *interface{} `json:"tag,omitempty"`           // Miscellaneous result
	MessagePrefix string       `json:"prefix,omitempty"`        // Prefix of the message to return

	ln        livenote.LiveNote // Internal note
	eventVerb string            // event verb related to the name of the operation
	osIsWin   bool
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
		Status:  string(EXCEPTION),
		ln:      livenote.LiveNote{},
		osIsWin: runtime.GOOS == "windows",
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

// Return sets the current status of a result
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

// AddInfo adds a formatted information message and returns itself
func (r *Result) AddInfo(fmtMsg string, a ...interface{}) Result {
	r.ln.AddInfo(fmt.Sprintf(fmtMsg, a...))
	r.updateMessage()
	return *r
}

// AddWarning adds a formatted warning message and returns itself
func (r *Result) AddWarning(fmtMsg string, a ...interface{}) Result {
	r.ln.AddWarning(fmt.Sprintf(fmtMsg, a...))
	r.updateMessage()
	return *r
}

// AddErrorf adds a formatted error message and returns itself
func (r *Result) AddError(fmtMsg string, a ...interface{}) Result {
	r.ln.AddError(fmt.Sprintf(fmtMsg, a...))
	r.updateMessage()
	return *r
}

// AddErr adds a error-typed value and returns itself.
func (r *Result) AddErr(err error) Result {
	r.AddError(err.Error())
	return *r
}

// AddErrWithAlt adds an error-typed value, and an alternate error
// message if the err happens to be nil. It returns itself.
func (r *Result) AddErrWithAlt(err error, altMsg string, altMsgValues ...any) Result {
	if err != nil {
		return r.AddErr(err)
	}
	if altMsg != "" {
		return r.AddError(altMsg, altMsgValues...)
	}
	return *r
}

// AddErrorWithAlt appends the messages of a Result.
// And an alternative message if the Result is other than OK or VALID status.
func (r *Result) AddErrorWithAlt(rs Result, altMsg string, altMsgValues ...any) Result {
	if !(rs.OK() || rs.Valid()) {
		for _, n := range rs.ln.Notes() {
			r.ln.Append(n)
		}
		r.updateMessage()
		return *r
	}
	if altMsg == "" {
		return *r
	}
	r.ln.Append(
		livenote.LiveNoteInfo{
			Type:    livenote.Error,
			Message: fmt.Sprintf(altMsg, altMsgValues...),
			Prefix:  r.ln.Prefix,
		})
	r.updateMessage()
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
func (r *Result) AppendError(rs Result, fmtMsg string, a ...interface{}) Result {
	for _, n := range rs.ln.Notes() {
		r.ln.Append(n)
	}
	return r.AddError(fmtMsg, a...)
}

// AppendInfof copies the messages of the Result parameter and append a formatted information message
func (r *Result) AppendInfo(rs Result, fmtMsg string, a ...interface{}) Result {
	for _, n := range rs.ln.Notes() {
		r.ln.Append(n)
	}
	return r.AddInfo(fmtMsg, a...)
}

// AppendWarning copies the messages of the Result parameter and append a formatted warning message
func (r *Result) AppendWarning(rs Result, fmtMsg string, a ...interface{}) Result {
	for _, n := range rs.ln.Notes() {
		r.ln.Append(n)
	}
	return r.AddWarning(fmtMsg, a...)
}

// Stuff adds or appends the messages of a Result.
func (r *Result) Stuff(rs Result) Result {
	for _, n := range rs.ln.Notes() {
		r.ln.Append(n)
	}
	r.updateMessage()
	return *r
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
	// if r.Messages is not empty, it can be because it was unmarshalled from result bytes
	if len(r.Messages) > 0 {
		lf := "\n"
		if r.osIsWin {
			lf = "\r\n"
		}
		sb := strings.Builder{}
		for _, v := range r.Messages {
			vlf := v + lf // prevents escape to the heap
			sb.Write([]byte(vlf))
		}
		return sb.String()
	}
	return r.ln.ToString()
}

// SetPrefix changes the prefix
func (r *Result) SetPrefix(pfx string) {
	r.ln.Prefix = pfx
	r.MessagePrefix = pfx
}

func (r *Result) updateMessage() {
	// get current notes to update the messages
	nts := r.ln.Notes()
	r.Messages = make([]string, 0, len(nts))
	for _, n := range nts {
		r.Messages = append(r.Messages, n.ToString())
	}
}

// RowsAffectedInfo - a function to simplify adding information for rows affected
func (r *Result) RowsAffectedInfo(rowsaff int64) {
	if rowsaff != 0 {
		r.AddInfo(fmt.Sprintf("%d rows affected", rowsaff))
	} else {
		r.AddInfo("No rows affected")
	}
}
