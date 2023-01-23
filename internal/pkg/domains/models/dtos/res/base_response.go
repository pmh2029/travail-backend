package res

type BaseResponse struct {
	Status string         `json:"status"`
	Data   interface{}    `json:"data"`
	Error  *ErrorResponse `json:"error,omitempty"`
}

// ErrorResponse struct
type ErrorResponse struct {
	ErrorCode    int    `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message"`
}
