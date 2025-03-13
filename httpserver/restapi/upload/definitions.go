package upload

type UploadDocumentRequest struct {
	Id      string                       `json:"id"`
	Type    string                       `json:"type"`
	Payload UploadDocumentRequestPayload `json:"payload"`
}
type UploadDocumentRequestPayload struct {
	Data string `json:"data"`
}

const (
	TEMP_FILE_PREFIX   = "part"
	UPLOAD_FILE_PREFIX = "file"
)
