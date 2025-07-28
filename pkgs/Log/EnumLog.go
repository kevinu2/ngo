package Log

type Type uint8

const (
	OutputStd Type = iota + 1
	OutputFile
	OutputBoth
)

func (lt Type) Type() string {
	switch lt {
	case OutputStd:
		return "std"
	case OutputFile:
		return "file"
	case OutputBoth:
		return "both"
	default:
		return "std"

	}
}

type Level uint8

const (
	Debug Level = iota + 1
	Info
	Warn
	Error
	DPanic
	Panic
	Fatal
)

func (ll Level) Level() string {
	switch ll {
	case Debug:
		return "Debug"
	case Info:
		return "Info"
	case Warn:
		return "Warn"
	case Error:
		return "Error"
	case DPanic:
		return "DPanic"
	case Panic:
		return "Panic"
	case Fatal:
		return "Fatal"
	default:
		return "Info"

	}
}
