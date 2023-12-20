package app

import (
	"app/internal/dto"
	"app/internal/repository/mts"
	"app/internal/services/httpServices"
	"log"
	"time"
)

func (r *App) importList() map[string]dto.LostSimCard {

	log.Println("Begin importList()")

	accountsList, err := r.spoRepo.GetOperatorsAccountsList()
	if err != nil {
		log.Println(err)
		return nil
	}

	lostSimCards := make(map[string]dto.LostSimCard)

	for aci := range accountsList {

		mtsHttpClient := httpServices.NewMtsHttpClient(&accountsList[aci], r.spoRepo, r.logger)
		mtsRepo := mts.NewMtsRepo(&accountsList[aci], mtsHttpClient, r.logger)
		// Пути выборки симкарт: Основная база и склад
		paths := []string{"sim_cards", "sim_cards/list_store"}
	out:
		for p := range paths {

			if r.shutdown() {
				return nil
			}

			currentPage := uint(1)

			timer := time.NewTimer(0)

			// Перебор страниц
			for {

				if r.shutdown() {
					return nil
				}

				select {
				case <-timer.C:
					simCards, err := mtsRepo.GetSimCards(paths[p], currentPage)

					if simCards == nil {
						break out
					}

					if err != nil {
						log.Println(err)
						break out
					}

					if len(simCards.Data.Items) > 0 {

						savedItemsCount := int(0)

						for simi := range simCards.Data.Items {

							cachedMtsSimCard := r.fastRepo.Get(simCards.Data.Items[simi].Sim.Iccid)

							r.fastRepo.Add(&simCards.Data.Items[simi])

							// Если есть изменения по сим-карте, то обновляем данные в основной БД
							if cachedMtsSimCard == nil ||
								simCards.Data.Items[simi].Status != cachedMtsSimCard.Status ||
								simCards.Data.Items[simi].State != cachedMtsSimCard.State ||
								simCards.Data.Items[simi].Sim.Imei != cachedMtsSimCard.Sim.Imei ||
								simCards.Data.Items[simi].Sim.Imsi != cachedMtsSimCard.Sim.Imsi ||
								simCards.Data.Items[simi].Sim.Account != cachedMtsSimCard.Sim.Account {

								saveSimErr := r.spoRepo.SimCardSave(&simCards.Data.Items[simi], &accountsList[aci])

								if saveSimErr != nil {
									r.logger.Error.Println(saveSimErr.Error())
								}

								savedItemsCount++

							}

							// Если сменился статус сим-карты с активной на блокированную
							if cachedMtsSimCard != nil {

								if cachedMtsSimCard.Status == "active" && simCards.Data.Items[simi].Status != "active" {
									lostSimCards[simCards.Data.Items[simi].Sim.Iccid] = dto.LostSimCard{Status: simCards.Data.Items[simi].Status, Imei: simCards.Data.Items[simi].Sim.Imei}
									r.logger.Info.Printf("Lost sim-card ICCID: %s, IMEI: %s\n", simCards.Data.Items[simi].Sim.Iccid, simCards.Data.Items[simi].Sim.Imei)
								}

								// Если сим-карта активна, удаляем из списка блокированных
								if simCards.Data.Items[simi].Status == "active" {
									delete(lostSimCards, simCards.Data.Items[simi].Sim.Iccid)
								}

								// Закомментировать
								//if simCards[i].Sim.Iccid == "89701015388560972001" || simCards[i].Sim.Iccid == "89701010085380071205" {
								//	lostSimCards[simCards[i].Sim.Iccid] = dto.LostSimCard{Status: simCards[i].Status, Imei: simCards[i].Sim.Imei}
								//}
							}

						}

						log.Printf("Account %s, path %s, page %d, getted items: %d, saved: %d\n",
							accountsList[aci].Username,
							paths[p],
							currentPage,
							len(simCards.Data.Items),
							savedItemsCount)

					} else {
						break out
					}
					currentPage++

					// Интевал между постраничными запросам
					timer.Reset(200 * time.Millisecond)
				default:
					time.Sleep(10 * time.Millisecond)
				}
			}
		}
	}
	log.Println("End importList()")
	return lostSimCards
}
