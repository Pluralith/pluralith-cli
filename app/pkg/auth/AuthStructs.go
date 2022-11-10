package auth

type OrgResponse struct {
	Code    int32        `json:"code"`
	Message string       `json:"message"`
	Data    PluralithOrg `json:"data"`
}

type PluralithOrg struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Icon    string   `json:"icon"`
	Invites struct{} `json:"invites"`
	Users   struct{} `json:"users"`
	Plan    struct{} `json:"plan"`
}

type ProjectResponse struct {
	Code    int32            `json:"code"`
	Message string           `json:"message"`
	Data    PluralithProject `json:"data"`
}

type PluralithProject struct {
	ID         string   `json:"id"`
	OrgID      string   `json:"orgId"`
	Name       string   `json:"name"`
	Branches   struct{} `json:"branches"`
	Connectors struct{} `json:"connectors"`
	Invites    struct{} `json:"invites"`
	Users      struct{} `json:"users"`
}
