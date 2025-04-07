package userclient

import (
	"TestTask/pkg/config"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type UserClient struct {
	Client *http.Client
	c      *config.Config
	// BaseURL string
}

func NewUserClient(c *config.Config) *UserClient {
	return &UserClient{
		Client: &http.Client{
			Timeout: 25 * time.Second,
		},
		c: c,
	}
}

type User struct {
	Age      int
	National string
	Gender   string
}

type NationalResponse struct {
	Country []Country `json:"country"`
}

type Country struct {
	Country_ID string `json:"country_id"`
}

type GenderResponse struct {
	Gender string `json:"gender"`
}

type AgeResponse struct {
	Age int `json:"age"`
}

func (uc *UserClient) GetNational(name string) string {
	url := fmt.Sprintf("%s%s", uc.c.GetURLApiForReqConfig().GetNationalURL, name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		errors.Wrap(err, "ошибка формирования запроса national")
		return ""
	}

	resp, err := uc.Client.Do(req)
	if err != nil {
		errors.Wrap(err, "ошибка отправки запроса national")
		return ""
	}

	defer resp.Body.Close()

	var nationalResponse NationalResponse
	err = json.NewDecoder(resp.Body).Decode(&nationalResponse)
	if err != nil {
		errors.Wrap(err, "ошибка в декодировании national")
		return ""
	}

	if len(nationalResponse.Country) > 0 {
		return nationalResponse.Country[0].Country_ID
	} else {
		return "Unknown"
	}
}

func (uc *UserClient) GetGender(name string) string {

	url := fmt.Sprintf("%s%s", uc.c.GetURLApiForReqConfig().GetGenderURL, name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		errors.Wrap(err, "ошибка формирования запроса gender")
		return ""
	}

	resp, err := uc.Client.Do(req)
	if err != nil {
		errors.Wrap(err, "ошибка отправки запроса gender")
		return ""
	}

	defer resp.Body.Close()

	var genderResponse GenderResponse
	err = json.NewDecoder(resp.Body).Decode(&genderResponse)
	if err != nil {
		errors.Wrap(err, "ошибка в декодировании gender")
		return ""
	}

	return genderResponse.Gender
}

func (uc *UserClient) GetAge(name string) int {
	url := fmt.Sprintf("%s%s", uc.c.GetURLApiForReqConfig().GetAgeURL, name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		errors.Wrap(err, "ошибка формирования запроса age")
		return 0
	}

	resp, err := uc.Client.Do(req)
	if err != nil {
		errors.Wrap(err, "ошибка отправки запроса age")
		return 0
	}

	defer resp.Body.Close()

	var ageResponse AgeResponse
	err = json.NewDecoder(resp.Body).Decode(&ageResponse)
	if err != nil {
		errors.Wrap(err, "ошибка в декодировании age")
		return 0
	}

	return ageResponse.Age
}
