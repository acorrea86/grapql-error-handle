package error_hanlde

import (
	"strconv"
)

type AppError string
type ValidationError2 string

const (
	NotFound               AppError = "Could not find resource"
	Unauthorized           AppError = "Unauthorized"
	Forbidden              AppError = "Forbidden"
	ErrorWithoutExtensions AppError = "No Extensions"
	AnyHow                 AppError = "Transparent"
	ServerError            AppError = "ServerError"
	MaxFileSizeError       AppError = "File size exceeds the maximum limit {0}"
	ContentTypeError       AppError = "Content Type not allowed {0}"
	DataSourceError        AppError = "Could not get info from datasource"
	ValidationError        AppError = "ValidationError:  {reason:?}  {code:?}"
)

type AppErrorRetry int

const (
	None         AppErrorRetry = 0
	Retry        AppErrorRetry = 4
	WaitAndRetry AppErrorRetry = 4
	Cancel       AppErrorRetry = -1
)

type ErrorExtensionValues struct {
	Reason string
	Code   string
	Level  string
}

func CreateExtensions(e *ErrorExtensionValues, reason, code string, level AppErrorRetry) *ErrorExtensionValues {
	e.Reason = reason
	e.Code = code
	e.Level = strconv.Itoa(int(level))
	return e
}

func (a AppError) CreateExtensionForAppError(err *error, reason, code *string) *ErrorExtensionValues {
	e := ErrorExtensionValues{}
	switch a {
	case NotFound:
		return CreateExtensions(&e, "Could not find resource", "NOT_FOUND", None)
	case ServerError:
		return CreateExtensions(&e, *reason, "SERVER_ERROR", Cancel)
	case DataSourceError:
		return CreateExtensions(&e, *reason, "DATA_SOURCE_ERROR", WaitAndRetry)
	case ValidationError:
		return CreateExtensions(&e, *reason, *code, None)
	case MaxFileSizeError:
		return CreateExtensions(&e, *reason, "MAX_FILE_SIZE_ERROR", Cancel)
	case ContentTypeError:
		return CreateExtensions(&e, *reason, "CONTENT_TYPE_ERROR", Cancel)
	case AnyHow:
		err := *err
		return CreateExtensions(&e, err.Error(), "SERVER_ERROR", Cancel)
	case ErrorWithoutExtensions:
		return nil
	case Unauthorized:
		return CreateExtensions(&e, "UNAUTHORIZED", "UNAUTHORIZED", Cancel)
	case Forbidden:
		return CreateExtensions(&e, "FORBIDDEN", "FORBIDDEN", Cancel)
	}
	return nil
}
