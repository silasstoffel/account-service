package loggerContract

type Logger interface {
	Info(message string, data map[string]interface{})
	Error(message string, err error, data map[string]interface{})
	Warn(message string, data map[string]interface{})
	Debug(message string, data map[string]interface{})
}

const (
	JsonFormat = "json"
	TextFormat = "text"
	InfoLevel  = "info"
	ErrorLevel = "error"
	DebugLevel = "debug"
)
