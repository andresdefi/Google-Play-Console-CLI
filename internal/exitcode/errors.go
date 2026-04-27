package exitcode

import (
	"errors"
	"fmt"
)

// ExitError is an error that carries a specific exit code.
type ExitError struct {
	Code    int
	Message string
}

func (e *ExitError) Error() string {
	return e.Message
}

// NewExitError creates a new ExitError with the given code and message.
func NewExitError(code int, format string, args ...any) *ExitError {
	return &ExitError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// AuthError creates an exit error for authentication failures.
func AuthError(format string, args ...any) *ExitError {
	return NewExitError(Auth, format, args...)
}

// APIErrorExit creates an exit error for API failures.
func APIErrorExit(format string, args ...any) *ExitError {
	return NewExitError(apiExitCode(args...), format, args...)
}

// ConfigError creates an exit error for configuration failures.
func ConfigError(format string, args ...any) *ExitError {
	return NewExitError(Config, format, args...)
}

type httpStatusCoder interface {
	HTTPStatusCode() int
}

func apiExitCode(args ...any) int {
	for _, arg := range args {
		if arg == nil {
			continue
		}
		if coder, ok := arg.(httpStatusCoder); ok {
			if code := FromHTTPStatus(coder.HTTPStatusCode()); code != Error {
				return code
			}
		}
		if err, ok := arg.(error); ok {
			var coder httpStatusCoder
			if errors.As(err, &coder) {
				if code := FromHTTPStatus(coder.HTTPStatusCode()); code != Error {
					return code
				}
			}
		}
	}
	return API
}
