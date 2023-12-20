package app

import (
	"app/internal/dto"
	"app/internal/repository/mts"
	"app/internal/services/httpServices"
	"log"
	"strconv"
	"time"
)

// Получение детализации по сим-картам
func (r *App) getDetails(accountsList []dto.Account) {

	log.Println("Begin getDetails")

	updatedCount := 0

	for acci := range accountsList {

		if r.shutdown() {
			return
		}

		mtsHttpClient := httpServices.NewMtsHttpClient(&accountsList[acci], r.spoRepo, r.logger)
		mtsRepo := mts.NewMtsRepo(&accountsList[acci], mtsHttpClient, r.logger)
		accountSimCards, err := r.spoRepo.GetAccountSimCardsIds(&accountsList[acci])

		if err != nil {
			r.logger.Error.Println(err.Error())
		}

	out:
		for j := range accountSimCards {

			if r.shutdown() {
				return
			}

			// Обновить Сим-карту по итогу выявления изменений
			doUpdate := false

			//if accountSimCards[j].SimId != 3636929 {
			//	continue
			//}

			balance, balErr := mtsRepo.GetBalance(accountSimCards[j].SimId)
			if balErr != nil {

				// Если в М2М сим-карта не найдена, удаляем её из БД и из кэша
				if balErr.Error() == "404" {
					log.Println("Delete SIM-card " + accountSimCards[j].ICCID)
					r.logger.Info.Println("Delete SIM-card " + accountSimCards[j].ICCID)
					delErr := r.spoRepo.DeleteSimCard(accountSimCards[j].ICCID)
					r.fastRepo.Del(accountSimCards[j].ICCID)
					if delErr != nil {
						r.logger.Info.Println("Ошибка удаления сим-карты " + accountSimCards[j].ICCID)
					}
					continue
				}

				log.Printf("error: %s, ", balErr)
				r.logger.Error.Printf("GetBalance request error for %s ICCID: %s: %s",
					accountsList[acci].Username, accountSimCards[j].ICCID, balErr)

				break out
			}

			if accountSimCards[j].Balance == nil || balance != *accountSimCards[j].Balance {
				doUpdate = true
			}

			accountSimCards[j].Balance = &balance

			time.Sleep(200 * time.Millisecond)

			services, servErr := mtsRepo.GetServices(accountSimCards[j].SimId)
			if servErr != nil {
				log.Printf("error: %s\n", servErr)
				r.logger.Error.Printf("GetServices request error for %s ICCID: %s: %s",
					accountsList[acci].Username, accountSimCards[j].ICCID, servErr)

				break out
			}

			// Сверяем сервисы (как минимум их кол-во)

			if accountSimCards[j].Services == nil || len(services) != len(*accountSimCards[j].Services) {
				doUpdate = true
			}

			if doUpdate {

				log.Printf("Account %s, SIM-card %s: update\n",
					accountsList[acci].Username, accountSimCards[j].ICCID)

				accountSimCards[j].Services = &services

				// Определяем ключ актуальной для нас услуги
				// Если APN == iot, то это передача данных, иначе пакеты NIDD

				needKey := "ts.scef_package"

				if accountSimCards[j].APN == "iot" {
					needKey = "ts.data_nbiot_package"
				}

				accountSimCards[j].EnabledServiceKey = &needKey

				for i := range services {

					serviceKey := services[i].ContractService.Service.BaseService.Key
					value := services[i].Value

					volume, err := strconv.Atoi(services[i].ContractService.Service.Parameters.Volume)
					if err != nil {
						volume = 1
					}

					if volume == 0 {
						volume = 1
					}

					// Если передача данных, то конвертируем МБайты в Байты
					if serviceKey == "ts.data_nbiot_package" {
						volume = volume * 1024 * 1024
					}

					// Если текущая услуга подходит к режиму сим-карты (NIDD или обычные данные)
					if serviceKey == needKey {
						accountSimCards[j].RemainingValue = &value
						accountSimCards[j].GivenValue = &volume
						percentage := uint16(float64(value * 100 / volume))
						accountSimCards[j].ResiduePercantage = &percentage

						PeriodEnd := services[i].PeriodEnd
						if len(*PeriodEnd) == 0 {
							PeriodEnd = nil
						}

						accountSimCards[j].ServiceUpdateDate = PeriodEnd
					}
				}

				upderr := r.spoRepo.UpdSimCard(&accountSimCards[j])

				if upderr != nil {
					r.logger.Error.Println(upderr)
				}

				updatedCount++

			} else {
				log.Printf("Account %s, SIM-card %s: skip\n",
					accountsList[acci].Username, accountSimCards[j].ICCID)
			}

			time.Sleep(200 * time.Millisecond)

		}

		log.Printf("Account %s, SIM-cards getted: %d, Details updated: %d\n",
			accountsList[acci].Username, len(accountSimCards), updatedCount)

	}

	log.Println("End getDetails")
}
