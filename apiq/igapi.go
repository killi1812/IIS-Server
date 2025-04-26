package apiq

import (
	"encoding/json"
	"errors"
	"fmt"
	"iis_server/config"
	"net/http"
	"net/url"

	"go.uber.org/zap"
)

type MockApi struct{}

func NewMockApi() *MockApi {
	return &MockApi{}
}

func (api *MockApi) GetUsernameByUserId(userId string) (*InstagramUsername, error) {
	baseURL := config.MOCK_API
	if config.OPTION_API == 2 {
		baseURL = config.REAL_API
	}
	tmpU, err := url.JoinPath(baseURL, "id")
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(tmpU)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("id", userId)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-rapidapi-key", config.RapidApiKey)
	req.Header.Add("x-rapidapi-host", "instagram-looter2.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	fmt.Printf("res: %v\n", res)
	if res.StatusCode != 200 {
		return nil, errors.New("not found")
	}

	var resp InstagramUsername
	json.NewDecoder(res.Body).Decode(&resp)
	defer res.Body.Close()
	zap.S().Debugf("respons from (%s), %+v", u.String(), resp)

	return &resp, nil
}

func (api *MockApi) GetUserInfoByUsername(username string) (*UserInfo, error) {
	baseURL := config.MOCK_API
	if config.OPTION_API == 2 {
		baseURL = config.REAL_API
	}
	tmpU, err := url.JoinPath(baseURL, "profile2")
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(tmpU)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("username", username)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-rapidapi-key", config.RapidApiKey)
	req.Header.Add("x-rapidapi-host", "instagram-looter2.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	fmt.Printf("res: %v\n", res)
	if res.StatusCode != 200 {
		return nil, errors.Join(ErrApiRequest, err)
	}

	var resp UserInfo
	json.NewDecoder(res.Body).Decode(&resp)
	defer res.Body.Close()
	zap.S().Debugf("response from (%s), %+v", u.String(), resp)

	return &resp, nil
}
