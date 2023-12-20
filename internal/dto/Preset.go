package dto

// Пресет системы оповещений
type Preset struct {
	Id     string     `db:"ID"`
	Data   PresetData `db:"DATA"`
	ChatId *string    `db:"CHAT_ID"`
}

type PresetData struct {
	Actions PresetDataActions `json:"actions"`
}

type PresetDataActions struct {
	TelegramBot bool     `json:"telegramBot"`
	Emails      []string `json:"emails"`
	CallPhones  []string `json:"callPhones"`
}
