package httpServices

import (
	"app/internal/dto"
	"app/internal/services/logger"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const UserAgent = "Kvant SPO SimCards Importer Go"
const ContentType = "application/json"
const HttpUnauth = "401"
const HttpBadRequest = "400"

type AccountProvider interface {
	UpdateAccount(*dto.Account)
}

type MtsHttpClient struct {
	account         *dto.Account
	accountProvider AccountProvider
	logger          *logger.Logger
}

func NewMtsHttpClient(account *dto.Account, accountProvider AccountProvider, logger *logger.Logger) *MtsHttpClient {
	return &MtsHttpClient{
		account,
		accountProvider,
		logger,
	}
}

// Выполнение рабочих запросов к серверу МТС
func (r *MtsHttpClient) ExecuteRequest(path string, body string) (string, error) {

	url := fmt.Sprintf("https://%s/api/v0/dashboards/%d/%s", r.account.Host, r.account.DasboardId, path)

	dataBufer := bytes.NewBuffer([]byte(""))

	httpMethod := http.MethodGet

	if len(body) > 0 {
		httpMethod = http.MethodPost
		var jsonStr = []byte(body)
		dataBufer = bytes.NewBuffer(jsonStr)
	}

	req, err := http.NewRequest(httpMethod, url, dataBufer)

	if err != nil {
		fmt.Println("NewRequest error:", err.Error())
	}

	req.Header.Set("Content-Type", ContentType)
	req.Header.Set("User-Agent", UserAgent)
	//req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.account.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		r.logger.Error.Panicln(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return "", errors.New("404")
	}

	if resp.StatusCode == 401 {

		err := r.RefreshToken()

		if err != nil {
			r.logger.Error.Println(err.Error())
			return "", err
		}

		return r.ExecuteRequest(path, body)
	}

	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	return string(respBody), nil
}

// Обновление токенов
func (r *MtsHttpClient) RefreshToken() error {

	if len(r.account.RefreshToken) == 0 {
		return r.Auth()
	}

	log.Printf("Get auth token by refresh for account %s\n\r", r.account.Username)

	url := fmt.Sprintf("https://%s/oauth/token", r.account.Host)

	var refreshData dto.MtsRefreshData
	json.Unmarshal([]byte(r.account.Params), &refreshData)
	refreshData.RefreshToken = r.account.RefreshToken
	refreshData.GrantType = "refresh_token"

	jsonStr, _ := json.Marshal(refreshData)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	if err != nil {
		fmt.Println("NewRequest error:", err.Error())
	}

	req.Header.Set("Content-Type", ContentType)
	req.Header.Set("User-Agent", UserAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		r.logger.Error.Panicln(err)
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 400 {
		return errors.New(HttpBadRequest)
	}

	if resp.StatusCode == 429 {
		return errors.New("429 Too Many Requests")
	}

	if resp.StatusCode == 401 {
		return r.Auth()
	}

	var authResp dto.MtsAuthResponce

	json.Unmarshal(respBody, &authResp)

	r.account.Token = authResp.AccessToken
	r.account.RefreshToken = authResp.RefreshToken
	r.account.OauthExpires = authResp.ExpiresIn
	r.accountProvider.UpdateAccount(r.account)

	return nil
}

// Авторизация на сервере МТС
func (r *MtsHttpClient) Auth() error {

	if len(r.account.Username) == 0 || len(r.account.Password) == 0 {
		return errors.New("empty login or password")
	}

	r.logger.Info.Printf("Get tokens by username and password for account %s\n\r", r.account.Username)

	url := fmt.Sprintf("https://%s/oauth/token", r.account.Host)

	var authData dto.MtsAuthData

	json.Unmarshal([]byte(r.account.Params), &authData)

	authData.Username = r.account.Username
	authData.Password = r.account.Password

	jsonStr, _ := json.Marshal(authData)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	if err != nil {
		fmt.Println("NewRequest error:", err.Error())
	}

	req.Header.Set("Content-Type", ContentType)
	req.Header.Set("User-Agent", UserAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		r.logger.Error.Panicln(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		return errors.New(HttpBadRequest)
	}

	if resp.StatusCode == 401 {
		r.logger.Error.Println("wrong login or password")
		return errors.New(HttpUnauth)
	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	var authResp dto.MtsAuthResponce

	json.Unmarshal(respBody, &authResp)

	r.account.Token = authResp.AccessToken
	r.account.RefreshToken = authResp.RefreshToken
	r.account.OauthExpires = authResp.ExpiresIn
	r.accountProvider.UpdateAccount(r.account)

	return nil
}
