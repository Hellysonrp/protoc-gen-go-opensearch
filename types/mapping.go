package types

type OpensearchDynamicTemplate struct {
	MatchMappingType string            `json:"match_mapping_type,omitempty"`
	Match            string            `json:"match,omitempty"`
	Unmatch          string            `json:"unmatch,omitempty"`
	Mapping          OpensearchMapping `json:"mapping"`
}

type OpensearchMapping struct {
	Analyzer         string                                 `json:"analyzer,omitempty"`
	Type             string                                 `json:"type,omitempty"`
	Format           string                                 `json:"format,omitempty"`
	Properties       map[string]OpensearchMapping           `json:"properties,omitempty"`
	Fields           map[string]OpensearchMapping           `json:"fields,omitempty"`
	Index            *bool                                  `json:"index,omitempty"`
	IgnoreAbove      int                                    `json:"ignore_above,omitempty"` // used with keyword field type
	DynamicTemplates []map[string]OpensearchDynamicTemplate `json:"dynamic_templates,omitempty"`
}
