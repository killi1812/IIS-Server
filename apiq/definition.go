package apiq

type IGApi interface {
	GetUsernameByUserId(userId string) (*InstagramUsername, error)
	GetUserInfoByUsername(username string) (*UserInfo, error)
}

func IgApiFactory() (IGApi, error) {
	return NewMockApi(), nil
}

type WeatherApi interface {
	GetWeatherForCity(city string) (*City, error)
}
