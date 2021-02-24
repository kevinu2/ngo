package log

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	zapper "go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.SugaredLogger

func InitLogger(logPath string, level zapper.Level) *zap.SugaredLogger {

	//logger,_:=zap.NewProduction()
	//defer logger.Sync()

	hook := lumberjack.Logger{
		Filename:   logPath, // ⽇志⽂件路径
		MaxSize:    1024,    // megabytes
		MaxBackups: 10,      // 最多保留10个备份
		MaxAge:     30,      //days
		Compress:   true,    // 是否压缩 disabled by default
	}
	w := zapper.AddSync(&hook)
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapper.ISO8601TimeEncoder
	encoderConfig.CallerKey = "l"
	core := zapper.NewCore(
		zapper.NewConsoleEncoder(encoderConfig),
		w,
		level,
	)
	//logger := zap.NewDevelopment()
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapper.ErrorLevel)).Sugar()

	return logger
}

func InitLoggerLevel(level string) *zap.SugaredLogger {
	switch level {
	case "Debug":
		return InitLogger(getLogPath(), zapper.DebugLevel)
	case "Info":
		return InitLogger(getLogPath(), zapper.InfoLevel)
	case "Warn":
		return InitLogger(getLogPath(), zapper.WarnLevel)
	case "Error":
		return InitLogger(getLogPath(), zapper.ErrorLevel)
	case "DPanic":
		return InitLogger(getLogPath(), zapper.DPanicLevel)
	case "Panic":
		return InitLogger(getLogPath(), zapper.PanicLevel)
	case "Fatal":
		return InitLogger(getLogPath(), zapper.FatalLevel)
	default:
		return InitLogger(getLogPath(), zapper.InfoLevel)
	}
}

func Logger() *zap.SugaredLogger {
	if logger == nil {
		logger = InitLoggerLevel(viper.GetString("log.level"))
	}
	return logger
}

func getLogPath() string {
	logPath := viper.GetString("log.path")
	logFile := viper.GetString("log.file")
	return logPath + "/" + logFile
}
