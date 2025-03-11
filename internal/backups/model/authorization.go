package model

type Session struct {
	Valid    bool   `json:"valid"`
	Totp     bool   `json:"totp"`
	Sid      string `json:sid`
	Csrf     string `json:csrf`
	Validity int    `json:validity`
}

type AuthorizationResponse struct {
	Session Session `json:"session"`
	Took    float32 `json:"took"`
}
