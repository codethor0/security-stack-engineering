package l2_identity

// Policy encodes zero trust decision logic.
type Policy struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Priority   int         `json:"priority"`
	Conditions []Condition `json:"conditions"`
	Action     string      `json:"action"`
}

// Condition is a single policy condition.
type Condition struct {
	Attribute string      `json:"attribute"`
	Operator  string      `json:"operator"`
	Value     interface{} `json:"value"`
}

// loadDefaultPolicies returns default zero trust policies.
func loadDefaultPolicies() []Policy {
	return []Policy{
		{
			ID:       "pol-allow-privileged-mfa",
			Name:     "Privileged access requires MFA",
			Priority: 100,
			Conditions: []Condition{
				{Attribute: "privileged", Operator: "eq", Value: true},
				{Attribute: "mfa_used", Operator: "eq", Value: true},
			},
			Action: "allow",
		},
		{
			ID:       "pol-deny-orphaned",
			Name:     "Deny orphaned privileged identities",
			Priority: 200,
			Conditions: []Condition{
				{Attribute: "owner", Operator: "eq", Value: ""},
				{Attribute: "privileged", Operator: "eq", Value: true},
			},
			Action: "deny",
		},
		{
			ID:       "pol-default-allow",
			Name:     "Default allow for known identities",
			Priority: 1000,
			Conditions: []Condition{
				{Attribute: "identity_found", Operator: "eq", Value: true},
			},
			Action: "allow",
		},
	}
}
