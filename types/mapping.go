package types

type OpensearchMapping struct {
	Analyzer    string                       `json:"analyzer,omitempty"`
	Type        string                       `json:"type,omitempty"`
	Properties  map[string]OpensearchMapping `json:"properties,omitempty"`
	Fields      map[string]OpensearchMapping `json:"fields,omitempty"`
	Index       *bool                        `json:"index,omitempty"`
	IgnoreAbove int                          `json:"ignore_above,omitempty"` // used with keyword field type
}
