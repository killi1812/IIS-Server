package upload

import (
	"errors"
	"fmt"
	"iis_server/config"
	"iis_server/httpserver/httpio"
	"iis_server/xmlvalidator"
	"net/http"
	"os"
	"strings"

	"github.com/killi1812/libxml2/types"
	"go.uber.org/zap"
)

func HandleUploadFile(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	method := parts[len(parts)-1]
	zap.S().Infof("Parsing method %s", method)

	file, header, err := r.FormFile("file")
	if err != nil {
		zap.S().Errorf("Can't read file from form err = %s", err)
		httpio.WriteStandardHTTPResponse(w, http.StatusInternalServerError, nil, err)
		return
	}

	if !strings.Contains(header.Filename, ".xml") {
		err = errors.New("invalid file extension should be .xml")
		httpio.WriteStandardHTTPResponse(w, http.StatusBadRequest, nil, err)
		return

	}

	data := make([]byte, header.Size)
	file.Read(data)
	defer file.Close()

	if err := xmlvalidator.Validate(data, method); err != nil {
		var validationErr types.SchemaValidationError
		var invalidXmlErr *xmlvalidator.ErrInvalidXML
		var valErrs error = nil
		switch {
		case errors.As(err, &validationErr):
			zap.S().Infof("errs: %+v\n", validationErr.Errors)
			for _, err2 := range validationErr.Errors {
				valErrs = errors.Join(valErrs, err2)
			}
			httpio.WriteStandardHTTPResponse(w, http.StatusOK, nil, valErrs)
			return

		case errors.As(err, &invalidXmlErr):
			zap.S().Infof("Invalid xml err = %v", invalidXmlErr)
			// TODO: return info on where is the invalid xml
			httpio.WriteStandardHTTPResponse(w, http.StatusOK, nil, invalidXmlErr)
			return

		default:
			zap.S().Errorf("Cannot Validate xml, err = %v", err)
			httpio.WriteStandardHTTPResponse(w, http.StatusInternalServerError, nil, err)
			return
		}
	}

	filePath := fmt.Sprintf("%v/%v", config.UPLOAD_FOLDER, header.Filename)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		zap.S().Errorf("Error writing a file %s, err = %s", filePath, err.Error())
		httpio.WriteStandardHTTPResponse(w, http.StatusInternalServerError, nil, err)
		return
	}

	httpio.WriteStandardHTTPResponse(w, http.StatusOK, filePath, nil)
}
