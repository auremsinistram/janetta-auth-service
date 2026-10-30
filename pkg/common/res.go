package common

type Response struct {
	Success bool `json:"success"`

	Data any `json:"data,omitempty"`

	Error *ErrorResponse `json:"error,omitempty"`
}

type ErrorResponse struct {
	Code int `json:"code"`

	Message string `json:"message"`
}

func (r *Response) MakeSuccess(data any) {
	r.Success = true
	r.Data = data
}

func (r *Response) MakeFailure(code int, message string) {
	r.Success = false
	r.Error = &ErrorResponse{
		Code:    code,
		Message: message,
	}
}
