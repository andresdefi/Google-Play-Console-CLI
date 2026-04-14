package exitcode

const (
	Success = 0
	Error   = 1
	Usage   = 2
	Auth    = 3
	API     = 4
	Config  = 5
)

// FromHTTPStatus maps an HTTP status code to an exit code.
func FromHTTPStatus(status int) int {
	switch {
	case status >= 200 && status < 300:
		return Success
	case status == 401, status == 403:
		return Auth
	case status == 404:
		return API
	case status == 409:
		return API
	case status >= 400 && status < 500:
		return API
	case status >= 500:
		return API
	default:
		return Error
	}
}
