package dto

type SimCardsListMTS struct {
	Data struct {
		Items SimCardsListMTSItems `json:"items"`
	} `json:"data"`
	Status string `json:"status"`
}

type SimCardsListMTSItems []SimCardMTS

type SimCardMTS struct {
	ID          int    `json:"id"`
	Status      string `json:"status"`
	State       string `json:"state"`
	Type        string `json:"type"`
	NetCheck    string `json:"net_check"`
	ContractID  int    `json:"contract_id"`
	ActivatedAt string `json:"activated_at"`
	Counters    []struct {
		Amount      int `json:"amount"`
		RaterType   int `json:"rater_type"`
		DirectionId int `json:"direction_id"`
	} `json:"counters"`
	Sim struct {
		ID         int    `json:"id"`
		Imsi       string `json:"imsi"`
		Iccid      string `json:"iccid"`
		Apn        string `json:"apn"`
		Imei       string `json:"imei"`
		Msisdn     string `json:"msisdn"`
		DeviceName string `json:"device_name"`
		Mode       string `json:"mode"`
		Tariff     int    `json:"tariff"`
		Account    string `json:"account"`
	} `json:"sim"`
}

type SimCardServicesMTS []SimCardServiceMTS

type SimCardServiceMTS struct {
	Value           int     `json:"value"`
	PeriodEnd       *string `json:"period_end"`
	ContractService struct {
		Service struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			ClientName  string `json:"client_name"`
			BaseService struct {
				Key string `json:"key"`
			} `json:"base_service"`
			Parameters struct {
				Volume string `json:"volume"`
			} `json:"parameters"`
		} `json:"service"`
	} `json:"contract_service"`
}
