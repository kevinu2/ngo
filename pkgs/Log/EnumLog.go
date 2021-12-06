package Log

type LogType uint8

const (
	LogOutputStd LogType = iota + 1
	LogOutputFile
	LogOutputBoth
)

func (lt LogType) Type() string {
	switch lt {
	case LogOutputStd:
		return "std"
	case LogOutputFile:
		return "file"
	case LogOutputBoth:
		return "both"
	default:
		return "std"

	}
}

type LogLevel uint8

const (
	LogDebug LogLevel = iota + 1
	LogInfo
	LogWarn
	LogError
	LogDPanic
	LogPanic
	LogFatal
)

func (ll LogLevel) Level() string {
	switch ll {
	case LogDebug:
		return "Debug"
	case LogInfo:
		return "Info"
	case LogWarn:
		return "Warn"
	case LogError:
		return "Error"
	case LogDPanic:
		return "DPanic"
	case LogPanic:
		return "Panic"
	case LogFatal:
		return "Fatal"
	default:
		return "Info"

	}
}
