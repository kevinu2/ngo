package Error

import "errors"

type ErrorType uint8

type ErrorCustom struct {
	ErrCode uint16
	Errors  error
}

const (
	ErrorNot200 ErrorType = iota + 1
	ErrorMapEmpty
	ErrorNotMatch
	ErrorMapNotEnough
	ErrorQueueEmpty
	ErrorValueNotChanged
	ErrorNotEqual
	ErrorNotFound
	ErrorNotAdded
	ErrorCountLimit
	ErrorLimitOver
	ErrorNotRequired
	ErrorNotEnough
	ErrorNotAllowed
	ErrorExpired
)

func (et ErrorType) GetMsg(msg string) error {
	switch et {
	case ErrorNot200:
		return errors.New(msg + " Http Code is not 200")
	case ErrorMapEmpty:
		return errors.New(msg + " Map is empty")
	case ErrorNotMatch:
		return errors.New(msg + " Not Match")
	case ErrorMapNotEnough:
		return errors.New(msg + " Map is not enough")
	case ErrorQueueEmpty:
		return errors.New(msg + " Queue is empty")
	case ErrorValueNotChanged:
		return errors.New(msg + " Value dose not change")
	case ErrorNotEqual:
		return errors.New(msg + " Not Equal")
	case ErrorNotFound:
		return errors.New(msg + " Not Found")
	case ErrorNotAdded:
		return errors.New(msg + " Not Added")
	case ErrorCountLimit:
		return errors.New(msg + " Limited 5")
	case ErrorLimitOver:
		return errors.New(msg + " Limit Over")
	case ErrorNotRequired:
		return errors.New(msg + " Not Required")
	case ErrorNotEnough:
		return errors.New(msg + " Not Enough")
	case ErrorNotAllowed:
		return errors.New(msg + " Not Allowed")
	case ErrorExpired:
		return errors.New(msg + " Expired")

	default:
		return nil
	}
}

func (et ErrorType) Get(msg string) ErrorCustom {
	switch et {
	case ErrorNot200:
		return ErrorCustom{ErrCode: 1, Errors: errors.New(msg + " Http Code is not 200")}
	case ErrorMapEmpty:
		return ErrorCustom{ErrCode: 2, Errors: errors.New(msg + " Map is empty")}
	case ErrorNotMatch:
		return ErrorCustom{ErrCode: 3, Errors: errors.New(msg + " Not Match")}
	case ErrorMapNotEnough:
		return ErrorCustom{ErrCode: 4, Errors: errors.New(msg + " Map is not enough")}
	case ErrorQueueEmpty:
		return ErrorCustom{ErrCode: 5, Errors: errors.New(msg + " Queue is empty")}
	case ErrorValueNotChanged:
		return ErrorCustom{ErrCode: 6, Errors: errors.New(msg + " Value dose not change")}
	case ErrorNotEqual:
		return ErrorCustom{ErrCode: 7, Errors: errors.New(msg + " Not Equal")}
	case ErrorNotFound:
		return ErrorCustom{ErrCode: 8, Errors: errors.New(msg + " Not Found")}
	case ErrorNotAdded:
		return ErrorCustom{ErrCode: 9, Errors: errors.New(msg + " Not Added")}
	case ErrorCountLimit:
		return ErrorCustom{ErrCode: 10, Errors: errors.New(msg + " Limited 5")}
	case ErrorLimitOver:
		return ErrorCustom{ErrCode: 11, Errors: errors.New(msg + " Limit Over")}
	case ErrorNotRequired:
		return ErrorCustom{ErrCode: 12, Errors: errors.New(msg + " Not Required")}
	case ErrorNotEnough:
		return ErrorCustom{ErrCode: 13, Errors: errors.New(msg + " Not Enough")}
	case ErrorNotAllowed:
		return ErrorCustom{ErrCode: 14, Errors: errors.New(msg + " Not Allowed")}
	case ErrorExpired:
		return ErrorCustom{ErrCode: 15, Errors: errors.New(msg + " Expired")}
	default:
		return ErrorCustom{ErrCode: 0, Errors: nil}
	}
}

func Err2Code(custom ErrorCustom) (errCode uint16, err error) {
	return custom.ErrCode, custom.Errors
}
