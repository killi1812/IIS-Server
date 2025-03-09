package xmlrpc

/*
import (
	"encoding/xml"
	"iis_server/apiq"
	"net/http"
)

func (_ *WeaterService) GetWeather(w http.ResponseWriter, r *http.Request) {
	data, err := apiq.GetWeatherForCity("Zagreb")
	if err != nil {
		// TODO: handl xmlrpc error return
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := xml.MarshalIndent(data, " ", " ")
	if err != err {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
*/
