package app

import (
	"app/internal/dto"
	"app/internal/repository/spo"
	"log"
)

// Создание уведомлений по поводу заблокированных сим-карт
func (r *App) createNotifications(lostSimCards map[string]dto.LostSimCard) {

	log.Println("createNotifications")

	if len(lostSimCards) == 0 {
		return
	}

	presets, err := r.spoRepo.GetSuperadminsPresets()

	if len(presets) == 0 {
		r.logger.Info.Println("Empty presets")
		return
	}

	if err != nil {
		r.logger.Error.Println(err.Error())
	}

	notifications := make(map[string]dto.NotifAction)

	// Перебор пресетов
	for i := range presets {

		if r.shutdown() {
			return
		}

		preset := presets[i]

		r.logger.Info.Print("Preset ")
		r.logger.Info.Println(preset.Id)

		if preset.Data.Actions.TelegramBot && preset.ChatId != nil {

			notifications[*preset.ChatId] = dto.NotifAction{
				Handler: "common\\actions\\handlers\\simCards\\TelegramBotHandler",
				Target:  *preset.ChatId,
				Data:    lostSimCards,
				Preset:  preset.Id,
				Status:  spo.NOTIF_STATUS_CREATED,
			}
		}

		for _, email := range preset.Data.Actions.Emails {
			notifications[email] = dto.NotifAction{
				Handler: "common\\actions\\handlers\\simCards\\EmailsHandler",
				Target:  email,
				Data:    lostSimCards,
				Preset:  preset.Id,
				Status:  spo.NOTIF_STATUS_CREATED,
			}
		}

		for _, phone := range preset.Data.Actions.CallPhones {
			notifications[phone] = dto.NotifAction{
				Handler: "common\\actions\\handlers\\simCards\\CallPhonesHandler",
				Target:  phone,
				Data:    lostSimCards,
				Preset:  preset.Id,
				Status:  spo.NOTIF_STATUS_CREATED,
			}
		}

	}

	if len(notifications) > 0 {
		err := r.spoRepo.AddNotifActions(notifications)
		if err != nil {
			r.logger.Error.Println(err.Error())
		}
	}

}
