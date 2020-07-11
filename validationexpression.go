package stdutil

// ValidationExpression - a struct for validation
type ValidationExpression struct {
	Name     string `json:"name,omitempty"`     // name of the database table column
	Value    string `json:"value,omitempty"`    // value of the column
	Operator string `json:"operator,omitempty"` // operator of the validation
}
