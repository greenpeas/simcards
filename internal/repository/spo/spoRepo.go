package spo

import (
	"app/internal/dto"
	"app/internal/services/logger"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const USER_ROLE_SUPER_ADMIN = 2
const EVENT_SEVERITY_HIGH = 8
const NOTIF_STATUS_CREATED = 0

type SpoRepo struct {
	dbpool *pgxpool.Pool
	logger *logger.Logger
	ctx    context.Context
}

func NewSpoRepo(dbpool *pgxpool.Pool, logger *logger.Logger, ctx context.Context) *SpoRepo {
	return &SpoRepo{
		dbpool,
		logger,
		ctx,
	}
}

// Получение списка аккаунтов операторов M2M
func (r *SpoRepo) GetOperatorsAccountsList() ([]dto.Account, error) {

	sql := `SELECT "ID","HOST", 
	COALESCE("OAUTH_ACCESS_TOKEN",'') as "OAUTH_ACCESS_TOKEN", 
	COALESCE("OAUTH_REFRESH_TOKEN",'') as "OAUTH_REFRESH_TOKEN",
	COALESCE("OAUTH_EXPIRES",0) as "OAUTH_EXPIRES",
	"USERNAME", "PASSWORD", "PARAMS", 
	("PARAMS"->>'dashboard_id')::int as "DASHBOARD_ID"
	FROM "REF_TELECOM_OPERATOR_ACCOUNTS"
	WHERE "IS_ACTIVE" is true
	ORDER BY "USERNAME"`

	rows, _ := r.dbpool.Query(r.ctx, sql)
	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[dto.Account])

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return items, nil
}

// Получение списка пресетов суперадминистраторов
func (r *SpoRepo) GetSuperadminsPresets() ([]dto.Preset, error) {

	sql := `SELECT p."ID", p."DATA", tlg."CHAT_ID"
FROM "REG_ABONENT_ACTIONS_PRESETS" p
INNER JOIN "REF_ABONENTS" a ON a."ID" = p."ABONENT" AND a."IS_ACTIVE" = 1 AND a."WEB_ROLE" = $1
LEFT JOIN "REG_TELEGRAM" tlg ON tlg."ABONENT" = p."ABONENT"
INNER JOIN "REG_ABONENTS_SETTINGS" stgn ON stgn."ABONENT" = p."ABONENT"
WHERE p."SEVERITY" & $2 > 0
AND (current_time at time zone stgn."TIMEZONE")::time(0) between p."TIME_START" and p."TIME_STOP"`

	rows, _ := r.dbpool.Query(r.ctx, sql, USER_ROLE_SUPER_ADMIN, EVENT_SEVERITY_HIGH)
	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[dto.Preset])

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return items, nil
}

// Сохранение симкарт в БД
func (r *SpoRepo) SimCardSave(simCard *dto.SimCardMTS, account *dto.Account) error {

	qs := `INSERT INTO "REF_SIM_CARDS" ("ICCID", "DATA", "ACCOUNT_ID")
	VALUES ($1, $2, $3)
	ON CONFLICT("ICCID") DO UPDATE SET "DATA"=EXCLUDED."DATA","UPDATEON"=now()`

	jsonData, err := json.Marshal(simCard)
	if err != nil {
		r.logger.Error.Println(err)
	}

	_, insertError := r.dbpool.Exec(r.ctx, qs, simCard.Sim.Iccid, jsonData, account.Id)
	if insertError != nil {
		return insertError
	}

	return nil
}

// Удаление сим-карты
func (r *SpoRepo) DeleteSimCard(iccid string) error {

	qs := `DELETE FROM "REF_SIM_CARDS" WHERE "ICCID" = $1`

	_, err := r.dbpool.Exec(r.ctx, qs, iccid)
	if err != nil {
		return err
	}

	return nil
}

func (r *SpoRepo) AddNotifAction(action *dto.NotifAction) error {

	qs := `INSERT INTO "REG_ACTIONS" ("HANDLER", "TARGET", "DATA", "PRESET", "STATUS") VALUES ($1, $2, $3, $4, $5)`

	_, err := r.dbpool.Exec(r.ctx, qs, action.Handler, action.Target, action.Data, action.Preset, action.Status)
	if err != nil {
		return err
	}

	return nil
}

func (r *SpoRepo) AddNotifActions(notifications map[string]dto.NotifAction) error {

	dollars := make([]string, 0, len(notifications))

	insertparams := []interface{}{}

	part := 0

	for _, notif := range notifications {
		dollarsPart := fmt.Sprintf("($%d,$%d,$%d,$%d,$%d)", part+1, part+2, part+3, part+4, part+5)
		dollars = append(dollars, dollarsPart)
		insertparams = append(insertparams, notif.Handler, notif.Target, notif.Data, notif.Preset, notif.Status)
		part += 5
	}

	dollarsString := strings.Join(dollars, ",")

	qs := fmt.Sprintf(`INSERT INTO "REG_ACTIONS" ("HANDLER", "TARGET", "DATA", "PRESET", "STATUS") VALUES %s`, dollarsString)

	_, insertError := r.dbpool.Exec(r.ctx, qs, insertparams...)
	if insertError != nil {
		return insertError
	}

	return nil
}

// Получение списка id симкарт аккаунта
func (r *SpoRepo) GetAccountSimCardsIds(account *dto.Account) ([]dto.SimCard, error) {

	sql := `SELECT 
	("DATA"->>'id')::int as "SIM_ID",
	"DATA",
	"ICCID",
	"BALANCE",
	"SERVICES",
	"ENABLED_SERVICE_KEY",
	"DATA"->'sim'->>'apn' as "APN",
	"REMAINING_VALUE",
	"GIVEN_VALUE",
	"RESIDUE_PERCENTAGE",
	"SERVICE_UPDATE_DATE"::text
	FROM "REF_SIM_CARDS" WHERE "ACCOUNT_ID" = $1
	ORDER BY "SERVICES" DESC NULLS FIRST`

	rows, err := r.dbpool.Query(r.ctx, sql, account.Id)

	if err != nil {
		return nil, err
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[dto.SimCard])

	if err != nil {
		return nil, err
	}

	return items, nil
}

// Получение списка id симкарт аккаунта
func (r *SpoRepo) GetSimCard(iccid string) (*dto.SimCard, error) {

	sql := `SELECT 
	("DATA"->>'id')::int as "SIM_ID",
	"DATA",
	"ICCID",
	"BALANCE",
	"SERVICES",
	"ENABLED_SERVICE_KEY",
	"DATA"->'sim'->>'apn' as "APN",
	"REMAINING_VALUE",
	"GIVEN_VALUE",
	"RESIDUE_PERCENTAGE",
	"SERVICE_UPDATE_DATE"::text
	FROM "REF_SIM_CARDS" WHERE "ICCID" = $1
	ORDER BY "SERVICES" DESC NULLS FIRST`

	rows, err := r.dbpool.Query(r.ctx, sql, iccid)

	if err != nil {
		return nil, err
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[dto.SimCard])

	if err != nil {
		return nil, err
	}

	return &items[0], nil
}

// Обновление сим-карты
func (r *SpoRepo) UpdSimCard(simCard *dto.SimCard) error {

	sql := `UPDATE "REF_SIM_CARDS" SET
	"BALANCE" = $2,
	"SERVICES" = $3,
	"ENABLED_SERVICE_KEY" = $4,
	"REMAINING_VALUE" = $5,
	"GIVEN_VALUE" = $6,
	"RESIDUE_PERCENTAGE" = $7,
	"SERVICE_UPDATE_DATE" = $8,
	"UPDATEON" = NOW()
	WHERE "ICCID" = $1`

	_, insertError := r.dbpool.Exec(r.ctx, sql, simCard.ICCID,
		simCard.Balance,
		simCard.Services,
		simCard.EnabledServiceKey,
		simCard.RemainingValue,
		simCard.GivenValue,
		simCard.ResiduePercantage,
		simCard.ServiceUpdateDate)
	if insertError != nil {
		return insertError
	}

	return nil
}

func (r *SpoRepo) UpdateAccount(acount *dto.Account) {

	namedArgs := pgx.NamedArgs{
		"OAUTH_ACCESS_TOKEN":  acount.Token,
		"OAUTH_REFRESH_TOKEN": acount.RefreshToken,
		"OAUTH_EXPIRES":       acount.OauthExpires,
		"ID":                  acount.Id,
	}

	q := `UPDATE "REF_TELECOM_OPERATOR_ACCOUNTS" SET
	"OAUTH_ACCESS_TOKEN" = @OAUTH_ACCESS_TOKEN,
	"OAUTH_REFRESH_TOKEN"= @OAUTH_REFRESH_TOKEN,
	"OAUTH_EXPIRES"      = @OAUTH_EXPIRES,
	"UPDATEON"           = NOW()
	WHERE "ID" = @ID`

	_, err := r.dbpool.Exec(r.ctx, q, namedArgs)

	if err != nil {
		fmt.Println(err)
	}
}
