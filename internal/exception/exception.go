package exception

const (
	UnknownError      = "UNKNOWN_ERROR"
	DbCommandError    = "DATABASE_ERROR"
	HttpClientError   = 400
	HttpInternalError = 500
	HttpNotFoundError = 404
)

type Exception struct {
	Code           string `json:"code"`
	Message        string `json:"message"`
	OriginalError  error  `json:"-"`
	httpStatusCode int    `json:"-"`
}

type ShortException struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func New(code string, message string, originalError error, httpStatus int) *Exception {
	status := 500
	if httpStatus != 0 {
		status = httpStatus
	}
	return &Exception{
		Code:           code,
		Message:        message,
		OriginalError:  originalError,
		httpStatusCode: status,
	}
}

func (e *Exception) ToDomain() ShortException {
	return ShortException{
		Code:    e.Code,
		Message: e.Message,
	}
}

func (e *Exception) Error() string {
	return e.Message
}
