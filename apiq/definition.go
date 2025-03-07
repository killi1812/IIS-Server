package apiq

import (
	"errors"
	"iis_server/config"
)

type IGApi interface {
	GetUsernameByUserId(userId string) (*InstagramUsername, error)
	GetUserInfoByUsername(username string) (*UserInfo, error)
}

func IgApiFactory() (IGApi, error) {
	switch config.OPTION_API {
	case 1:
		return NewMockApi(), nil
		// TODO: Implement a real api
	case 2:
		return nil, nil

	default:
		return nil, errors.New("bad option")
	}
}

type WeatherApi interface {
	GetWeatherForCity(city string) (*City, error)
}
