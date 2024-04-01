package opa

// Token define the rego auth token
type Token struct {
	Type     string
	Username string
	Roles    []string
}

// EntityData define the entity owner
type EntityData struct {
	Type  string
	Owner string
}

// PolicyInput input policy data
type PolicyInput struct {
	Token      Token
	EntityData EntityData
}
