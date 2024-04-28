package helper

type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func InvalidInputFormat() *Response {
	return &Response{
		Code:    "input.invalid_format",
		Message: "Invalid input format",
	}
}

func ValidationFailure(m string) *Response {
	return &Response{Code: "validation.failure", Message: m}
}
