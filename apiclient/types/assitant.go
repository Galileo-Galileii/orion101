package types

type Assistant struct {
	Metadata
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Icons       AgentIcons `json:"icons"`
	EntityID    string     `json:"entityID"`
}

type AssistantList List[Assistant]

type AssistantTool struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon,omitempty"`
	Enabled     bool   `json:"enabled,omitempty"`
	Builtin     bool   `json:"builtin,omitempty"`
}

type AssistantToolList struct {
	ReadOnly bool            `json:"readOnly,omitempty"`
	Items    []AssistantTool `json:"items"`
}
