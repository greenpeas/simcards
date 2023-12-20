package dto

type MtsAuthResponce struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    uint32 `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
