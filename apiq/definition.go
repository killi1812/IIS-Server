package apiq

import (
	"errors"
	"iis_server/config"
)

type IgApi interface {
	GetUsernameByUserId(userId string) (*InstagramUsername, error)
}

func IgApiFactory() (IgApi, error) {
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
