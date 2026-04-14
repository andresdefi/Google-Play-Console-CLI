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
	if got := FromHTTPStatus(404); got != NotFound {
		t.Errorf("FromHTTPStatus(404) = %d, want %d (NotFound)", got, NotFound)
	}
}

func TestFromHTTPStatus_409(t *testing.T) {
	if got := FromHTTPStatus(409); got != Conflict {
		t.Errorf("FromHTTPStatus(409) = %d, want %d (Conflict)", got, Conflict)
	}
}

func TestFromHTTPStatus_400_BadRequest(t *testing.T) {
	if got := FromHTTPStatus(400); got != 10 {
		t.Errorf("FromHTTPStatus(400) = %d, want 10", got)
	}
}

func TestFromHTTPStatus_422_Unprocessable(t *testing.T) {
	if got := FromHTTPStatus(422); got != 32 {
		t.Errorf("FromHTTPStatus(422) = %d, want 32", got)
	}
}

func TestFromHTTPStatus_429_TooManyRequests(t *testing.T) {
	if got := FromHTTPStatus(429); got != 39 {
		t.Errorf("FromHTTPStatus(429) = %d, want 39", got)
	}
}

func TestFromHTTPStatus_500_Internal(t *testing.T) {
	if got := FromHTTPStatus(500); got != 60 {
		t.Errorf("FromHTTPStatus(500) = %d, want 60", got)
	}
}

func TestFromHTTPStatus_502_BadGateway(t *testing.T) {
	if got := FromHTTPStatus(502); got != 62 {
		t.Errorf("FromHTTPStatus(502) = %d, want 62", got)
	}
}

func TestFromHTTPStatus_503_ServiceUnavailable(t *testing.T) {
	if got := FromHTTPStatus(503); got != 63 {
		t.Errorf("FromHTTPStatus(503) = %d, want 63", got)
	}
}

func TestFromHTTPStatus_Unknown(t *testing.T) {
	for _, code := range []int{0, 100, 199, 301, 302} {
		if got := FromHTTPStatus(code); got != Error {
			t.Errorf("FromHTTPStatus(%d) = %d, want %d (Error)", code, got, Error)
		}
	}
}

func TestExitCode_NotFound_Constant(t *testing.T) {
	if NotFound != 4 {
		t.Errorf("NotFound should be 4, got %d", NotFound)
	}
}

func TestExitCode_Conflict_Constant(t *testing.T) {
	if Conflict != 5 {
		t.Errorf("Conflict should be 5, got %d", Conflict)
	}
}

func TestExitCode_Config_Constant(t *testing.T) {
	if Config != 6 {
		t.Errorf("Config should be 6, got %d", Config)
	}
}

func TestExitCode_API_Alias(t *testing.T) {
	if API != NotFound {
		t.Errorf("API should equal NotFound (%d), got %d", NotFound, API)
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
	if NotFound != 4 {
		t.Errorf("NotFound should be 4, got %d", NotFound)
	}
	if Conflict != 5 {
		t.Errorf("Conflict should be 5, got %d", Conflict)
	}
	if Config != 6 {
		t.Errorf("Config should be 6, got %d", Config)
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
