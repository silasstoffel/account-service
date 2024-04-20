package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/silasstoffel/account-service/configs"
)

type Logger struct {
	Env     string
	Service string
}

func init() {
	//zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

func NewLogger(config *configs.Config) *Logger {
	return &Logger{
		Env:     config.Env,
		Service: config.App.AppName,
	}
}

func NewLoggerWithService(config *configs.Config, serviceName string) *Logger {
	return &Logger{
		Env:     config.Env,
		Service: serviceName,
	}
}

func (ref *Logger) Info(message string, data map[string]interface{}) {
	log.Info().Fields(ref.appendDefaultInput(data)).Msg(message)
}

func (ref *Logger) Error(message string, err error, data map[string]interface{}) {
	log.Error().Stack().Fields(ref.appendDefaultInput(data)).Err(err).Msg(message)
}

func (ref *Logger) Warn(message string, data map[string]interface{}) {
	log.Warn().Fields(ref.appendDefaultInput(data)).Msg(message)
}

func (ref *Logger) Debug(message string, data map[string]interface{}) {
	log.Debug().Fields(ref.appendDefaultInput(data)).Msg(message)
}

func (ref *Logger) appendDefaultInput(data map[string]interface{}) map[string]interface{} {
	if data == nil {
		data = make(map[string]interface{})
	}
	data["env"] = ref.Env
	data["service"] = ref.Service

	return data
}
