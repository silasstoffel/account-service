package helper

type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func InvalidInputFormat() *Response {
	return &Response{
		Code:    "INVALID_INPUT_FORMAT",
		Message: "Invalid input format",
	}
}
