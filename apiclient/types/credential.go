package types

type Credential struct {
	ContextID string   `json:"contextID,omitempty"`
	Name      string   `json:"name,omitempty"`
	EnvVars   []string `json:"envVars,omitempty"`
	ExpiresAt *Time    `json:"expiresAt,omitempty"`
}

type CredentialList List[Credential]
