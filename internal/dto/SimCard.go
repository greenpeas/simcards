package dto

type SimCard struct {
	SimId             int                 `db:"SIM_ID"`
	ICCID             string              `db:"ICCID"`
	Balance           *int                `db:"BALANCE"`
	Services          *SimCardServicesMTS `db:"SERVICES"`
	EnabledServiceKey *string             `db:"ENABLED_SERVICE_KEY"`
	APN               string              `db:"APN"`
	RemainingValue    *int                `db:"REMAINING_VALUE"`
	GivenValue        *int                `db:"GIVEN_VALUE"`
	ResiduePercantage *uint16             `db:"RESIDUE_PERCENTAGE"`
	ServiceUpdateDate *string             `db:"SERVICE_UPDATE_DATE"`
	Data              struct {
		Id  int `db:"id"`
		Sim struct {
			IMEI   string `db:"imei"`
			IMSI   string `db:"imsi"`
			ICCID  string `db:"iccid"`
			MSISDN string `db:"msisdn"`
		} `db:"sim"`
		Status string `db:"status"`
		State  string `db:"state"`
	} `db:"DATA"`
}
