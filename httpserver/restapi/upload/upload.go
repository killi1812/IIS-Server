package upload

import (
	"errors"
	"fmt"
	"iis_server/config"
	"iis_server/httpserver/httpio"
	"iis_server/xmlvalidator"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/lestrrat-go/libxml2/xsd"
	"go.uber.org/zap"
)

// TODO: Simplify function mby no need tmp file
func HandleUploadFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	parts := strings.Split(r.URL.Path, "/")
	method := parts[len(parts)-1]
	zap.S().Infof("Parsing method %s", method)

	file, header, err := r.FormFile("file")
	if err != nil {
		zap.S().Errorf("Can't read file from form err = %s", err)
		httpio.WriteStandardHTTPResponse(w, http.StatusInternalServerError, nil, err)
	}
	data := make([]byte, header.Size)
	file.Read(data)
	defer file.Close()

	if err := xmlvalidator.Validate(data, method); err != nil {
		var validationErr xsd.SchemaValidationError
		var invalidXmlErr *xmlvalidator.ErrInvalidXML
		valErrs := []string{}
		switch {
		case errors.As(err, &validationErr):
			zap.S().Infof("errs: %+v\n", validationErr.Errors())
			for _, err2 := range validationErr.Errors() {
				valErrs = append(valErrs, err2.Error())
			}
			httpio.WriteStandardHTTPResponse(w, http.StatusOK, valErrs, nil)

		case errors.As(err, &invalidXmlErr):
			zap.S().Infof("Invalid xml err = %v", invalidXmlErr)
			// TODO: return info on where is the invalid xml
			httpio.WriteStandardHTTPResponse(w, http.StatusOK, invalidXmlErr.Error(), nil)

		default:
			zap.S().Errorf("Cannot Validate xml, err = %v", err)
			httpio.WriteStandardHTTPResponse(w, http.StatusInternalServerError, nil, err)
		}
		return
	}

	tempFile, err := os.CreateTemp(config.TMP_FOLDER, "part*")
	if err != nil {
		zap.S().Errorf("Cannot create temp file, err = %v", err)
		httpio.WriteStandardHTTPResponse(w, http.StatusInternalServerError, nil, err)
		return
	}

	tempName := tempFile.Name()
	zap.S().Infof("Temp file name = %v", tempName)
	defer tempFile.Close()

	n, err := io.Copy(tempFile, file)
	if err != nil {
		zap.S().Errorf("Cannot copy content to temp file, err = %v", err)
		// write response
		httpio.WriteStandardHTTPResponse(w, http.StatusInternalServerError, nil, err)
		return

	}
	zap.S().Debugf("Number of bytes written = %v", n)
	tempFile.Close()

	filePath := fmt.Sprintf("%v/%v", config.UPLOAD_FOLDER, header.Filename)
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
