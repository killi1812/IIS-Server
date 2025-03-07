package xmlrpc

import (
	"fmt"
	"iis_server/apiq"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func GetWeather(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	city := ps.ByName("city")
	data, err := apiq.GetWeatherForCity(city)
	if err != nil {
		// TODO: handl xmlrpc error return
		return
	}
	fmt.Printf("data: %v\n", data)
}
