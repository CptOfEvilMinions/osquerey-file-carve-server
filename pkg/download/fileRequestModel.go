package download

// FileRequest
type FileRequest struct {
	Token         string `json:"token"`
	TokenAccessor string `json:"token_accessor"`
	FileCarveGUID string `json:"file_carve_guid"`
}
