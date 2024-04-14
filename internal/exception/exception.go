package exception

type Exception struct {
	Code          string `json:"code"`
	Message       string `json:"message"`
	OriginalError error  `json:"-"`
	StatusCode    int    `json:"-"`
}

func New(code string, err *error) *Exception {
	m, c := GetMessageByCode(code)
	return &Exception{
		Code:          code,
		Message:       m,
		OriginalError: *err,
		StatusCode:    c,
	}
}

func NewUnknown(err *error) *Exception {
	m, c := GetMessageByCode(UnknownError)
	return &Exception{
		Code:          UnknownError,
		Message:       m,
		OriginalError: *err,
		StatusCode:    c,
	}
}

func NewDbCommandError(err *error) *Exception {
	m, c := GetMessageByCode(DbCommandError)
	return &Exception{
		Code:          DbCommandError,
		Message:       m,
		OriginalError: *err,
		StatusCode:    c,
	}
}

func (e *Exception) GetStatusCode() int {
	return e.StatusCode
}

func (e *Exception) Error() string {
	return e.Message
}
