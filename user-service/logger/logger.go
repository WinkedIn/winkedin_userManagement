package logger

import (
	"context"
	"os"
	"runtime"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/winkedin/user-service/constants"
)

var (
	isTerm         = isatty.IsTerminal(os.Stdout.Fd())
	loggerInstance zerolog.Logger
	once           sync.Once
)

// getFuncName returns the name of the function that called the function that called getFuncName
func getFuncName() string {
	pc, _, _, _ := runtime.Caller(3)
	return runtime.FuncForPC(pc).Name()
}

// InitLogger initializes the logger
func InitLogger() {
	once.Do(func() {
		loggerInstance = zerolog.New(gin.DefaultWriter).
			Output(
				zerolog.ConsoleWriter{
					Out:     gin.DefaultWriter,
					NoColor: !isTerm,
				},
			).
			With().
			Timestamp().
			Logger()
	})
}

// LogFunctionPointWithContext TODO - NEED TO DEVISE A WAY TO LOG FUNCTION PARAMETERS
// LogFunctionPointWithContext logs the entry and exit of a function
func LogFunctionPointWithContext(ctx context.Context, point string) {
	// get func name
	funcName := getFuncName()

	// get request params from context
	requestID, ok := ctx.Value(constants.Key(constants.RequestIdKey)).(string)
	if !ok {
		requestID = "unknown"
	}
	userIP, ok := ctx.Value(constants.Key(constants.UserIP)).(string)
	if !ok {
		userIP = "unknown"
	}

	var message string
	if point == constants.LogFunctionEntry {
		message = "Entering function"
	} else if point == constants.LogFunctionExit {
		message = "Exiting function"
	}

	l := loggerInstance

	evt := l.WithLevel(constants.LogLevelInfo)
	evt.Str(constants.RequestIdKey, requestID).
		Str(constants.UserIP, userIP).
		Str("func", funcName).
		Msg(message)
}

// LogErrorWithContext logs an error
func LogErrorWithContext(ctx context.Context, message string) {
	// get func name
	funcName := getFuncName()

	// get request params from context
	requestID, ok := ctx.Value(constants.Key(constants.RequestIdKey)).(string)
	if !ok {
		requestID = "unknown"
	}
	userIP, ok := ctx.Value(constants.Key(constants.UserIP)).(string)
	if !ok {
		userIP = "unknown"
	}
	requestPath, ok := ctx.Value(constants.Key(constants.RequestPath)).(string)
	if !ok {
		requestPath = "unknown"
	}
	requestMethod, ok := ctx.Value(constants.Key(constants.RequestMethod)).(string)
	if !ok {
		requestMethod = "unknown"
	}

	l := loggerInstance

	evt := l.WithLevel(constants.LogLevelInfo)
	evt.Str(constants.RequestIdKey, requestID).
		Str(constants.RequestMethod, requestMethod).
		Str(constants.UserIP, userIP).
		Str(constants.RequestPath, requestPath).
		Str(constants.RequestMethod, requestMethod).
		Str("func", funcName).
		Msg(message)
}
