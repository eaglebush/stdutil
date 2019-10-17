package stdutil

// ValidationExpression - a struct for validation
type ValidationExpression struct {
	Name     string `json:"name,omitempty"`
	Value    string `json:"value,omitempty"`
	Operator string `json:"operator,omitempty"`
}
