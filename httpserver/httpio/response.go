package httpio

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type StandardHTTPResponse struct {
	StatusCode int    `json:"statusCode"`
	Status     string `json:"status"`
	Error      string `json:"error"`
	TimeStamp  string `json:"timeStamp"`
	Data       any    `json:"data"`
}

func GetStandardHTTPResponse(Data any, err error) StandardHTTPResponse {
	var httpResponse StandardHTTPResponse
	var code int

	if err == nil {
		httpResponse.StatusCode = 0
		httpResponse.Status = "success"

	} else {
		errW := errors.Unwrap(err)
		if errW == nil {
			errW = err
		}
		httpResponse.StatusCode = code
		httpResponse.Status = errW.Error()

		httpResponse.Error = err.Error()
	}

	httpResponse.TimeStamp = time.Now().Format(time.RFC3339Nano)
	httpResponse.Data = Data

	return httpResponse
}

func WriteStandardHTTPResponse(w http.ResponseWriter, httpStatus int, payload any, err error) {
	httpResponse := GetStandardHTTPResponse(payload, err)

	// log response
	// - log all data
	zap.S().Debugf("Response to client, http status = %v, response = %+v", httpStatus, httpResponse)
	// - log status only
	// zap.S().Debugf("Response to client, http status = %v, status = %v, status code = %v, status ext = %v",
	// 	httpStatus, httpResponse.Status, httpResponse.StatusCode, httpResponse.StatusExt)

	// write header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	// write data
	jsonResp, err := json.MarshalIndent(httpResponse, "", " ")
	if err != nil {
		zap.S().Errorf("Cannot convert response to json, err = %v", err)
		jsonResp = []byte{}
	}

	_, err = w.Write(jsonResp)
	if err != nil {
		zap.S().Errorf("Cannot write response to client, err = %v", err)
		return
	}
}
