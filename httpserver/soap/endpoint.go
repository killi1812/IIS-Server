package soap

import (
	"encoding/xml"
	"iis_server/apiq"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// SOAP request structure
type SOAPRequest struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		XMLName xml.Name `xml:"Body"`
		GetData struct {
			XMLName xml.Name `xml:"GetData"`
			Input   string   `xml:"input"`
		} `xml:"GetData"`
	} `xml:"Body"`
}

// SOAP response structure
type SOAPResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Xmlns   string   `xml:"xmlns,attr"`
	Body    struct {
		XMLName         xml.Name `xml:"Body"`
		GetDataResponse struct {
			XMLName xml.Name      `xml:"GetDataResponse"`
			User    apiq.UserInfo `xml:"UserInfo"`
		} `xml:"GetDataResponse"`
	} `xml:"Body"`
}

func RegisterEnpint(router *mux.Router) {
	router.HandleFunc("/soap", handleSOAPRequest).Methods("POST", "OPTIONS")
}

func handleSOAPRequest(w http.ResponseWriter, r *http.Request) {
	zap.S().Debugf("query: %+v", r.URL.Query())

	if r.URL.Query().Get("wsdl") != "" {
		data, err := os.ReadFile("httpserver/soap/userInfo.wsdl")
		if err != nil {
			// TODO: return some err
			zap.S().DPanicf("Error requesting wsdl query:%+v, err = %+v", r.URL.Query(), err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}

	body, _ := io.ReadAll(r.Body)
	var req SOAPRequest
	if err := xml.Unmarshal(body, &req); err != nil {
		// TODO: return error
		zap.S().Errorf("Error reading request")
		soapErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	api, err := apiq.IgApiFactory()
	if err != nil {
		// TODO: return error
		zap.S().DPanicf("Error creating api")
		return
	}

	// TODO: add test username
	data, err := api.GetUserInfoByUsername("")
	if err != nil {
		// TODO: return error
		zap.S().Errorf("Error Retriving data from an api")
		soapErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	soapResponse(w, http.StatusOK, data)
}

func soapErrorResponse(w http.ResponseWriter, code int, data any) {
	// TODO: see how to encode errors
}

func soapResponse(w http.ResponseWriter, code int, data *apiq.UserInfo) {
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")

	response := SOAPResponse{
		Xmlns: "http://example.com/soap",
		Body: struct {
			XMLName         xml.Name `xml:"Body"`
			GetDataResponse struct {
				XMLName xml.Name      `xml:"GetDataResponse"`
				User    apiq.UserInfo `xml:"UserInfo"`
			} `xml:"GetDataResponse"`
		}{
			GetDataResponse: struct {
				XMLName xml.Name      `xml:"GetDataResponse"`
				User    apiq.UserInfo `xml:"UserInfo"`
			}{
				User: *data,
			},
		},
	}

	zap.S().Debugf("Responding with code = %d, response = %+v", code, data)
	xml.NewEncoder(w).Encode(response)
}
