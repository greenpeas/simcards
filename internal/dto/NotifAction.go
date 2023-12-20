package dto

type NotifAction struct {
	Handler string                 `db:"HANDLER"`
	Target  string                 `db:"TARGET"`
	Data    map[string]LostSimCard `db:"DATA"`
	Preset  string                 `db:"PRESET"`
	Status  uint                   `db:"STATUS"`
}
