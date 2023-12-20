package mts

import (
	"app/internal/dto"
	"app/internal/services/httpServices"
	"app/internal/services/logger"
	"encoding/json"
	"fmt"
)

type MtsRepo struct {
	account       *dto.Account
	mtsHttpClient *httpServices.MtsHttpClient
	logger        *logger.Logger
}

func NewMtsRepo(account *dto.Account, mtsHttpClient *httpServices.MtsHttpClient, logger *logger.Logger) *MtsRepo {
	return &MtsRepo{
		account,
		mtsHttpClient,
		logger,
	}
}

// Получить одну страницу сим-карт от оператора
func (r *MtsRepo) GetSimCards(path string, page uint) (*dto.SimCardsListMTS, error) {

	body := fmt.Sprintf(`{"page":%d,"per_page":500}`, page)

	responce, err := r.mtsHttpClient.ExecuteRequest(path, body)

	if err != nil {
		r.logger.Error.Println(err)
		return nil, err
	}

	var result dto.SimCardsListMTS

	json.Unmarshal([]byte(responce), &result)

	return &result, nil
}

// Получение баланса сим-карты
func (r *MtsRepo) GetBalance(id int) (int, error) {

	path := fmt.Sprintf(`sim_cards/%d/get_balance`, id)

	responce, err := r.mtsHttpClient.ExecuteRequest(path, "")

	if err != nil {
		return 0, err
	}

	result := struct {
		Amount int `json:"amount"`
	}{}

	json.Unmarshal([]byte(responce), &result)

	return result.Amount, nil
}

// Получение сервисов сим-карты
func (r *MtsRepo) GetServices(id int) (dto.SimCardServicesMTS, error) {

	path := fmt.Sprintf(`sim_cards/%d/services`, id)

	responce, err := r.mtsHttpClient.ExecuteRequest(path, "")

	if err != nil {
		r.logger.Error.Println(err)
		return dto.SimCardServicesMTS{}, err
	}

	var result struct {
		Data struct {
			Items dto.SimCardServicesMTS `json:"items"`
		} `json:"data"`
	}

	json.Unmarshal([]byte(responce), &result)

	return result.Data.Items, nil
}
