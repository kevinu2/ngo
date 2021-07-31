package Log

import (
	"fmt"
	rotates "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"ngo/constant"
	"ngo/enum"
	"ngo/model"
	"os"
	"time"
)

var l *Log

type Log struct {
	logger *zap.SugaredLogger
	config *model.LogConfig
}

func init() {
	l = New()
}

func New() *Log {
	return new(Log)
}

func AddConfig(logLevel, logPath, logFile, logOutput string) {
	l.addConfig(logLevel, logPath, logFile, logOutput)
}
func (l *Log) addConfig(logLevel, logPath, logFile, logOutput string) {
	l.config = &model.LogConfig{
		LogLevel:  logLevel,
		LogPath:   logPath,
		LogFile:   logFile,
		LogOutput: logOutput,
	}
}

func (l *Log) initLogger(level zapcore.Level) {
	var core zapcore.Core

	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(constant.TimeUtcFormat))
		},
		MessageKey:   "msg",
		LevelKey:     "level",
		TimeKey:      "ts",
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	switch l.config.LogOutput {
	case enum.LogOutputStd.Type():
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),
		)
	case enum.LogOutputFile.Type():
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(getWriter(fmt.Sprintf("%s/%s", l.config.LogPath, l.config.LogFile))), level),
		)
		//logger := zap.NewDevelopment()
	case enum.LogOutputBoth.Type():
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),
			zapcore.NewCore(encoder, zapcore.AddSync(getWriter(fmt.Sprintf("%s/%s", l.config.LogPath, l.config.LogFile))), level),
		)
	default:
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),
		)
	}
	l.logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).Sugar()
	return
}

func (l *Log) InitLoggerLevel() {
	switch l.config.LogLevel {
	case enum.LogDebug.Level():
		l.initLogger(zapcore.DebugLevel)
		return
	case enum.LogInfo.Level():
		l.initLogger(zapcore.InfoLevel)
		return
	case enum.LogWarn.Level():
		l.initLogger(zapcore.WarnLevel)
		return
	case enum.LogError.Level():
		l.initLogger(zapcore.ErrorLevel)
		return
	case enum.LogDPanic.Level():
		l.initLogger(zapcore.DPanicLevel)
		return
	case enum.LogPanic.Level():
		l.initLogger(zapcore.PanicLevel)
		return
	case enum.LogFatal.Level():
		l.initLogger(zapcore.FatalLevel)
		return
	default:
		l.initLogger(zapcore.InfoLevel)
		return
	}
}

func Logger() *zap.SugaredLogger {
	if l.logger == nil {
		l.InitLoggerLevel()
	}
	return l.logger
}

func getWriter(filename string) io.Writer {
	hook, err := rotates.New(
		filename+".%Y%m%d%H",
		rotates.WithLinkName(filename),
		rotates.WithMaxAge(time.Hour*24*7),
		rotates.WithRotationTime(time.Hour),
	)

	if err != nil {
		panic(err)
	}
	return hook
}
