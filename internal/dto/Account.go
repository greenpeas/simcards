package dto

type Account struct {
	Id           string `db:"ID"`
	Host         string `db:"HOST"`
	Token        string `db:"OAUTH_ACCESS_TOKEN"`
	RefreshToken string `db:"OAUTH_REFRESH_TOKEN"`
	OauthExpires uint32 `db:"OAUTH_EXPIRES"`
	Username     string `db:"USERNAME"`
	Password     string `db:"PASSWORD"`
	DasboardId   uint32 `db:"DASHBOARD_ID"`
	Params       string `db:"PARAMS"`
	//UpdateOn     time.Time `db:"UPDATEON"`
}
