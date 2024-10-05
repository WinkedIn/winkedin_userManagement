package constants

import "github.com/rs/zerolog"

type Key string

const (
	RequestIdKey     string = "request_id"
	UserIP           string = "user_ip"
	RequestPath      string = "request_path"
	RequestMethod    string = "request_method"
	LogFunctionEntry string = "function_entry"
	LogFunctionExit  string = "function_exit"
	LogLevelDebug           = zerolog.DebugLevel
	LogLevelInfo            = zerolog.InfoLevel
	LogLevelWarn            = zerolog.WarnLevel
	LogLevelError           = zerolog.ErrorLevel
	LogLevelFatal           = zerolog.FatalLevel
	LogLevelPanic           = zerolog.PanicLevel
)
