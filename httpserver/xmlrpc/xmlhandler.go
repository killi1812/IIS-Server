package xmlrpc

import (
	"bytes"
	"encoding/xml"
	"iis_server/apiq"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type RPCRequest struct {
	XMLName    xml.Name `xml:"methodCall"`
	MethodName string   `xml:"methodName"`
	Param      RPCParam `xml:"param"`
}

type RPCParam struct {
	Value ParamValue `xml:"value"`
}

type ParamValue struct {
	String string `xml:"string"`
}

type RPCResponse struct {
	XMLName xml.Name  `xml:"methodResponse"`
	Params  RPCResult `xml:"params"`
}

type RPCResult struct {
	Param []ResultValue `xml:"param"`
}

type ResultValue struct {
	Location string `xml:"value>Location"`
	Temp     string `xml:"value>Temptriture"`
}

func RegisterEndpoint(router *mux.Router) {
	router.HandleFunc("/weather", xmlRPCHandler).Methods("POST", "OPTIONS")
}

// XML-RPC Handler (Custom Implementation)
func xmlRPCHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req RPCRequest
	if err := xml.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid XML", http.StatusBadRequest)
		return
	}

	if req.MethodName != "GetTemp" {
		http.Error(w, "Method not supported", http.StatusNotImplemented)
		return
	}

	// Process request
	data, err := new(apiq.WeaterService).GetWeatherForCity(req.Param.Value.String)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var resp RPCResponse
	for _, c := range data {
		resp.Params.Param = append(resp.Params.Param, ResultValue{c.GradIme, strings.Trim(c.Podatci.Temp, " ")})
	}

	var buf bytes.Buffer
	xml.NewEncoder(&buf).Encode(resp)
	w.Write(buf.Bytes())
}
