package upload

type UploadDocumentRequest struct {
	Id      string                       `json:"id"`
	Type    string                       `json:"type"`
	Payload UploadDocumentRequestPayload `json:"payload"`
}
type UploadDocumentRequestPayload struct {
	Data string `json:"data"`
}

type UploadDocumentResponsePayload struct {
	FileName string `json:"fileName"`
}

type UploadFileResponsePayload struct {
	FileName string `json:"fileName"`
}

const (
	TEMP_FILE_PREFIX   = "part"
	UPLOAD_FILE_PREFIX = "file"
)
