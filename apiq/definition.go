package apiq

import (
	"errors"
	"iis_server/config"
)

type IgApi interface {
	// TODO: make real return types
	GetUsernameByUserId(userId string) (*InstagramUsername, error)
}

func IgApiFactory() (IgApi, error) {
	switch config.OPTION_API {
	case 1:
		return NewMockApi(), nil

	case 2:
		return nil, nil

	default:
		return nil, errors.New("bad option")
	}
}
