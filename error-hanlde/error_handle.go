package error_hanlde

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
	UnauthorizedReason              = "UNAUTHORIZED"
	ForbiddenReason                 = "FORBIDDEN"
)

type AppErrorRetry string

const (
	None         AppErrorRetry = "NONE"
	Retry        AppErrorRetry = "RETRY"
	WaitAndRetry AppErrorRetry = "WAIT_AND_RETRY"
	Cancel       AppErrorRetry = "CANCEL"
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

func createExtensions(reason, code string, level AppErrorRetry) *ErrorExtensionValues {
	return &ErrorExtensionValues{
		Reason: reason,
		Code:   code,
		Level:  string(level),
	}
}

type ErrorExtensionParams struct {
	Reason   string
	Code     string
	AppError AppError
}

func CreateExtensionForAppError(params ErrorExtensionParams) *ErrorExtensionValues {
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
