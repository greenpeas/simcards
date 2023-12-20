package dto

type MtsRefreshData struct {
	RefreshToken string `json:"refresh_token"`
	ClientId     uint   `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}
