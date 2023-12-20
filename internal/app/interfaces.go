package app

import (
	"app/internal/dto"
)

// Репозиторий быстрой БД, например Redis
type FastRepo interface {

	// Получение сим-карты из кеша
	Get(string) string

	// Добавление сим-карты в кеш
	Add(string, string)

	// Удаление сим-карты из кэша
	Del(string)
}

// Репозиторий СПО
type SpoRepo interface {

	// Получение пресетов суперадминистраторов системы
	GetSuperadminsPresets() ([]dto.Preset, error)

	// Сохранение сим-карты в базу данных
	SimCardSave(simCard *dto.SimCardMTS, account *dto.Account) error

	AddNotifActions(map[string]dto.NotifAction) error

	AddNotifAction(*dto.NotifAction) error

	GetOperatorsAccountsList() ([]dto.Account, error)

	GetAccountSimCardsIds(*dto.Account) ([]dto.SimCard, error)

	UpdSimCard(*dto.SimCard) error

	UpdateAccount(*dto.Account)

	DeleteSimCard(string) error
}
