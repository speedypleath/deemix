package deezer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
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
	Results struct {
		Data []struct {
			Token string `json:"TRACK_TOKEN"`
		} `json:"data"`
	} `json:"results"`
}

type UrlResponse struct {
	Data []struct {
		Media []struct {
			Format  string `json:"format"`
			Sources []struct {
				Url string `json:"url"`
			} `json:"sources"`
		} `json:"media"`
	} `json:"data"`
}

type Session struct {
	Sid       string
	ApiToken  string
	UserToken string
	License   string
}

type Response interface {
	PingResponse | UserDataResponse | ListDataResponse | UrlResponse
}

func httpRequest[T Response](method string, url string, sid string, data []byte) *T {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		log.Error(err)
		return nil
	}

	header := fmt.Sprintf("sid=%s", sid)
	req.Header.Set("Cookie", header)

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

	r := new(T)
	err = json.Unmarshal(body, r)
	if err != nil {
		log.Error("Error unmarshalling json")
		return nil
	}

	return r
}

func Ping(session *Session) {
	r := httpRequest[PingResponse]("GET", "https://www.deezer.com/ajax/gw-light.php?method=deezer.ping&api_version=1.0&api_token", "", nil)
	session.Sid = r.Results.Session
}

func UserData(session *Session) {
	r := httpRequest[UserDataResponse]("GET", "https://www.deezer.com/ajax/gw-light.php?method=deezer.getUserData&api_version=1.0&api_token", session.Sid, nil)
	session.ApiToken = r.Results.ApiToken
	session.License = r.Results.User.Options.LicenseToken
}

func GetListData(tracks []string, session *Session) []string {
	url := fmt.Sprintf("https://www.deezer.com/ajax/gw-light.php?method=song.getListData&api_version=1.0&api_token=%s", session.ApiToken)
	body := fmt.Sprintf(`{"sng_ids":%s}`, tracks)
	r := httpRequest[ListDataResponse]("POST", url, session.Sid, []byte(body))
	return funk.Map(r.Results.Data, func(x struct {
		Token string `json:"TRACK_TOKEN"`
	}) string {
		return fmt.Sprintf(`"%s"`, x.Token)
	}).([]string)
}

func GetStreamUrl(tracks []string, session Session) string {
	body := fmt.Sprintf(`{
		"license_token": "%s",
		"media": [{
			"type": "FULL",
			"formats": [{
				"cipher": "BF_CBC_STRIPE",
				"format": "MP3_128"
				}]
			}],
		"track_tokens": %s}`, session.License, tracks)
	r := httpRequest[UrlResponse]("POST", "https://media.deezer.com/v1/get_url", session.Sid, []byte(body))
	return r.Data[0].Media[0].Sources[0].Url
}
