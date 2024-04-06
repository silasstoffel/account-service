package exception

const (
	UnknownError      = "UNKNOWN_ERROR"
	DbCommandError    = "DATABASE_ERROR"
	HttpClientError   = 400
	HttpUnauthorized  = 401
	HttpNotFoundError = 404
	HttpInternalError = 500
)

type Exception struct {
	Code           string `json:"code"`
	Message        string `json:"message"`
	OriginalError  error  `json:"-"`
	HttpStatusCode int    `json:"-"`
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
		HttpStatusCode: status,
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
