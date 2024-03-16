package domain

const (
	UnknownError   = "UNKNOWN_ERROR"
	DbCommandError = "DATABASE_ERROR"
)

type Error struct {
	Code          string `json:"code"`
	Message       string `json:"message"`
	OriginalError error  `json:"-"`
}

type ShortError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewError(code string, message string, originalError error) *Error {
	return &Error{
		Code:          code,
		Message:       message,
		OriginalError: originalError,
	}
}

func (e *Error) ToDomain() ShortError {
	return ShortError{
		Code:    e.Code,
		Message: e.Message,
	}
}

func (e *Error) Error() string {
	return e.Message
}
