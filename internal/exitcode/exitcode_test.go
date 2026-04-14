package exitcode

import (
	"testing"
)

func TestFromHTTPStatus_2xx(t *testing.T) {
	for _, code := range []int{200, 201, 204, 299} {
		if got := FromHTTPStatus(code); got != Success {
			t.Errorf("FromHTTPStatus(%d) = %d, want %d (Success)", code, got, Success)
		}
	}
}

func TestFromHTTPStatus_401(t *testing.T) {
	if got := FromHTTPStatus(401); got != Auth {
		t.Errorf("FromHTTPStatus(401) = %d, want %d (Auth)", got, Auth)
	}
}

func TestFromHTTPStatus_403(t *testing.T) {
	if got := FromHTTPStatus(403); got != Auth {
		t.Errorf("FromHTTPStatus(403) = %d, want %d (Auth)", got, Auth)
	}
}

func TestFromHTTPStatus_404(t *testing.T) {
	if got := FromHTTPStatus(404); got != API {
		t.Errorf("FromHTTPStatus(404) = %d, want %d (API)", got, API)
	}
}

func TestFromHTTPStatus_409(t *testing.T) {
	if got := FromHTTPStatus(409); got != API {
		t.Errorf("FromHTTPStatus(409) = %d, want %d (API)", got, API)
	}
}

func TestFromHTTPStatus_4xx(t *testing.T) {
	for _, code := range []int{400, 422, 429, 451} {
		if got := FromHTTPStatus(code); got != API {
			t.Errorf("FromHTTPStatus(%d) = %d, want %d (API)", code, got, API)
		}
	}
}

func TestFromHTTPStatus_5xx(t *testing.T) {
	for _, code := range []int{500, 502, 503, 504} {
		if got := FromHTTPStatus(code); got != API {
			t.Errorf("FromHTTPStatus(%d) = %d, want %d (API)", code, got, API)
		}
	}
}

func TestFromHTTPStatus_Unknown(t *testing.T) {
	for _, code := range []int{0, 100, 199, 301, 302} {
		if got := FromHTTPStatus(code); got != Error {
			t.Errorf("FromHTTPStatus(%d) = %d, want %d (Error)", code, got, Error)
		}
	}
}

func TestExitError_Error(t *testing.T) {
	e := &ExitError{Code: 1, Message: "something failed"}
	if e.Error() != "something failed" {
		t.Errorf("expected 'something failed', got %q", e.Error())
	}
}

func TestNewExitError(t *testing.T) {
	e := NewExitError(API, "api error: %s", "not found")
	if e.Code != API {
		t.Errorf("expected code %d, got %d", API, e.Code)
	}
	if e.Message != "api error: not found" {
		t.Errorf("expected 'api error: not found', got %q", e.Message)
	}
}

func TestAuthError(t *testing.T) {
	e := AuthError("not authenticated: %s", "missing token")
	if e.Code != Auth {
		t.Errorf("expected code %d, got %d", Auth, e.Code)
	}
	if e.Message != "not authenticated: missing token" {
		t.Errorf("expected 'not authenticated: missing token', got %q", e.Message)
	}
}

func TestAPIErrorExit(t *testing.T) {
	e := APIErrorExit("api failed: %d", 500)
	if e.Code != API {
		t.Errorf("expected code %d, got %d", API, e.Code)
	}
	if e.Message != "api failed: 500" {
		t.Errorf("expected 'api failed: 500', got %q", e.Message)
	}
}

func TestConfigError(t *testing.T) {
	e := ConfigError("config broken: %s", "bad json")
	if e.Code != Config {
		t.Errorf("expected code %d, got %d", Config, e.Code)
	}
	if e.Message != "config broken: bad json" {
		t.Errorf("expected 'config broken: bad json', got %q", e.Message)
	}
}

func TestExitCode_Constants(t *testing.T) {
	if Success != 0 {
		t.Errorf("Success should be 0, got %d", Success)
	}
	if Error != 1 {
		t.Errorf("Error should be 1, got %d", Error)
	}
	if Usage != 2 {
		t.Errorf("Usage should be 2, got %d", Usage)
	}
	if Auth != 3 {
		t.Errorf("Auth should be 3, got %d", Auth)
	}
	if API != 4 {
		t.Errorf("API should be 4, got %d", API)
	}
	if Config != 5 {
		t.Errorf("Config should be 5, got %d", Config)
	}
}

func TestExitError_ImplementsError(t *testing.T) {
	var _ error = &ExitError{}
}

func TestNewExitError_NoArgs(t *testing.T) {
	e := NewExitError(Error, "plain message")
	if e.Message != "plain message" {
		t.Errorf("expected 'plain message', got %q", e.Message)
	}
}
