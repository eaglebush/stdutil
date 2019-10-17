package stdutil

//Result - standard result structure
type Result struct {
	Execution    string //Values: SUCCESS, FAIL
	Status       string //OK, ERROR, VALID or any status
	FocusControl string //Control to focus when error was activated
	Messages     []string
}

// InitResult - initialize result for API query
func InitResult() Result {
	return Result{Execution: "FAIL", Status: "ERROR", Messages: make([]string, 0)}
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
	return r.Status == "OK"
}

//IsStatusError - checks if the status is Error
func (r *Result) IsStatusError() bool {
	return r.Status == "ERROR"
}

//IsStatusValid - checks if the status is Valid
func (r *Result) IsStatusValid() bool {
	return r.Status == "VALID"
}

//IsStatusInvalid - checks if the status is invalid
func (r *Result) IsStatusInvalid() bool {
	return r.Status == "INVALID"
}

//IsStatusYes - checks if the status is Yes
func (r *Result) IsStatusYes() bool {
	return r.Status == "YES"
}

//IsStatusNo - checks if the status is No
func (r *Result) IsStatusNo() bool {
	return r.Status == "NO"
}
