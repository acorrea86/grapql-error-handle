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

var decisionMapLevel = map[AppError]AppErrorRetry{
	NotFound:         None,
	ServerError:      Cancel,
	DataSourceError:  WaitAndRetry,
	ValidationError:  None,
	MaxFileSizeError: Cancel,
	ContentTypeError: Cancel,
	AnyHow:           Cancel,
	Unauthorized:     Cancel,
	Forbidden:        Cancel,
}

type ErrorExtensionValues struct {
	Reason string
	Code   string
	Level  string
}

type ErrorExtensionParams struct {
	Reason   string
	Code     string
	AppError AppError
}

func createExtensions(reason, code string, level AppErrorRetry) *ErrorExtensionValues {
	return &ErrorExtensionValues{
		Reason: reason,
		Code:   code,
		Level:  strconv.Itoa(int(level)),
	}
}

func CreateExtensionForAppError(params ErrorExtensionParams) *ErrorExtensionValues {
	switch params.AppError {
	case NotFound:
		return createExtensions(params.Reason, "NOT_FOUND", None)
	case ServerError:
		return createExtensions(params.Reason, "SERVER_ERROR", Cancel)
	case DataSourceError:
		return createExtensions(params.Reason, "DATA_SOURCE_ERROR", WaitAndRetry)
	case ValidationError:
		return createExtensions(params.Reason, params.Code, None)
	case MaxFileSizeError:
		return createExtensions(params.Reason, "MAX_FILE_SIZE_ERROR", Cancel)
	case ContentTypeError:
		return createExtensions(params.Reason, "CONTENT_TYPE_ERROR", Cancel)
	case AnyHow:
		//return createExtensions(params.Err.Error(), "SERVER_ERROR", Cancel)
		return createExtensions(params.Reason, "SERVER_ERROR", Cancel)
	case ErrorWithoutExtensions:
		return nil
	case Unauthorized:
		return createExtensions(params.Reason, "UNAUTHORIZED", Cancel)
		//return createExtensions("UNAUTHORIZED", "UNAUTHORIZED", Cancel)
	case Forbidden:
		//return createExtensions("FORBIDDEN", "FORBIDDEN", Cancel)
		return createExtensions(params.Reason, "FORBIDDEN", Cancel)
	}
	return nil
}

func CreateExtensionForAppErrorWithMap(params ErrorExtensionParams) *ErrorExtensionValues {
	code := ""
	retry := None

	if params.AppError == ErrorWithoutExtensions {
		return nil
	}

	decisionMapCode := map[AppError]string{
		NotFound:         "NOT_FOUND",
		ServerError:      "SERVER_ERROR",
		DataSourceError:  "DATA_SOURCE_ERROR",
		ValidationError:  params.Code,
		MaxFileSizeError: "MAX_FILE_SIZE_ERROR",
		ContentTypeError: "CONTENT_TYPE_ERROR",
		AnyHow:           "SERVER_ERROR",
		Unauthorized:     "UNAUTHORIZED",
		Forbidden:        "FORBIDDEN",
	}

	for key, decision := range decisionMapCode {
		if key == params.AppError {
			code = decision
			break
		}
	}

	for key, appErrorRetry := range decisionMapLevel {
		if key == params.AppError {
			retry = appErrorRetry
			break
		}
	}

	return createExtensions(params.Reason, code, retry)

}
