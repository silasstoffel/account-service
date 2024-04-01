package domain

type Permission struct {
	AppId  string `json:"appId"`
	Scope  string `json:"scope"`
	Active bool   `json:"active"`
}

type AccountPermission struct {
	AccountId string `json:"accountId"`
	AppId     string `json:"appId"`
	Scope     string `json:"scope"`
}
