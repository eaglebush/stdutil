package stdutil

type NameValue[T any] struct {
	Name  string `json:"name,omitempty"`
	Value T      `json:"value,omitempty"`
}
