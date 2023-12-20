package dto

type MtsAuthData struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	ClientId     uint   `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}
