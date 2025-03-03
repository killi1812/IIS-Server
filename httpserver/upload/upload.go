package upload

import (
	"fmt"
	"iis_server/config"
	"iis_server/httpserver/httpio"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

func HandleUploadFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	parts := strings.Split(r.URL.Path, "/")
	method := parts[len(parts)-1]
	zap.S().Infof("Parsing method %s", method)

	// request header
	uploadFileName := r.Header.Get("FileName")

	// log client request
	zap.S().Debugf("Client request = %+v", r)
	zap.S().Debugf("Upload file name = %+v", uploadFileName)

	// create temp file
	tempFile, err := os.CreateTemp(config.TMP_FOLDER, "part*")
	if err != nil {
		zap.S().Errorf("Cannot create temp file, err = %v", err)
		// write response
		httpio.WriteStandardHTTPResponse(w, http.StatusInternalServerError, nil, err)
		return
	}

	// get temp file name
	tempName := tempFile.Name()
	zap.S().Debugf("Temp file name = %v", tempName)
	defer tempFile.Close()

	// copy to temp file
	n, err := io.Copy(tempFile, r.Body)
	if err != nil {
		zap.S().Errorf("Cannot copy content to temp file, err = %v", err)

		// write response
		httpio.WriteStandardHTTPResponse(w, http.StatusInternalServerError, nil, err)
		return

	}
	zap.S().Debugf("Number of bytes written = %v", n)
	tempFile.Close()

	// get name of upload file
	filePath := ""

	if uploadFileName != "" {
		// set file path in format = folder + name
		filePath = fmt.Sprintf("%v/%v", config.UPLOAD_FOLDER, uploadFileName)
	} else {
		// create temp file in upload folder
		uploadFile, err := os.CreateTemp(config.UPLOAD_FOLDER, UPLOAD_FILE_PREFIX+"*")
		if err != nil {
			zap.S().Errorf("Cannot create temp file, err = %v", err)

			// write response
			httpio.WriteStandardHTTPResponse(w, http.StatusInternalServerError, nil, err)
			return
		}

		filePath = uploadFile.Name()
		uploadFile.Close()
	}
	zap.S().Debugf("Upload file path = %v", filePath)

	err = os.Rename(tempName, filePath)
	if err != nil {
		zap.S().Errorf("Cannot rename temp file, err = %v", err)

		// write response
		httpio.WriteStandardHTTPResponse(w, http.StatusInternalServerError, nil, err)
		return
	}
	zap.S().Debugf("Renamed temp file to %v", filePath)

	// write response
	respPayload := UploadFileResponsePayload{
		FileName: filePath,
	}

	httpio.WriteStandardHTTPResponse(w, http.StatusOK, respPayload, nil)
}
