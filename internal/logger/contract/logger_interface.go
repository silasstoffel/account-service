package loggerContract

type Logger interface {
	Info(message, data interface{})
	Error(message, err error, data interface{})
	Warn(message, data interface{})
	Debug(message, data interface{})
}

const (
	JsonFormat = "json"
	TextFormat = "text"
	InfoLevel  = "info"
	ErrorLevel = "error"
	DebugLevel = "debug"
)
