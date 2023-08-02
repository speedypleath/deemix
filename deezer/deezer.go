package deezer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type PingResponse struct {
	Results struct {
		Session string `json:"SESSION"`
	} `json:"results"`
}

type UserDataResponse struct {
	Results struct {
		Token    string `json:"USER_TOKEN"`
		ApiToken string `json:"checkForm"`
		User     struct {
			Options struct {
				LicenseToken string `json:"license_token"`
			} `json:"OPTIONS"`
		} `json:"USER"`
	} `json:"results"`
}

type ListDataResponse struct {
}

type Response interface {
	PingResponse | UserDataResponse | ListDataResponse
}

func GetRequest[T Response](url string) *T {
	resp, err := http.Get(url)
	if err != nil {
		log.Error(err)
		return nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil
	}

	r := new(T)

	json.Unmarshal(body, r)

	return r
}

func PostRequest[T Response](url string, sid string, data []byte) *T {
	log.Info(url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Error(err)
		return nil
	}

	req.Header.Set("Cookie", fmt.Sprintf("sid=%s", sid))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil
	}

	log.Info(string(body))

	r := new(T)

	json.Unmarshal(body, r)

	return r
}

func Ping() *PingResponse {
	return GetRequest[PingResponse]("https://www.deezer.com/ajax/gw-light.php?method=deezer.ping&api_version=1.0&api_token")
}

func UserData() *UserDataResponse {
	return GetRequest[UserDataResponse]("https://www.deezer.com/ajax/gw-light.php?method=deezer.getUserData&api_version=1.0&api_token")
}

func GetListData(tracks []string) *ListDataResponse {
	ping := Ping()
	userData := UserData()
	url := fmt.Sprintf("https://www.deezer.com/ajax/gw-light.php?method=song.getListData&api_version=1.0&api_token=%s", userData.Results.ApiToken)
	return PostRequest[ListDataResponse](url, ping.Results.Session, []byte(fmt.Sprintf(`{"sng_ids":[%s]`, tracks)))
}
