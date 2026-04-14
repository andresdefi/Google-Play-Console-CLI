package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// --- Client creation ---

func TestNewClient(t *testing.T) {
	c := NewClient("test-token")
	if c == nil {
		t.Fatal("expected non-nil client")
	}
	if c.token != "test-token" {
		t.Errorf("expected token 'test-token', got %q", c.token)
	}
	if c.baseURL != baseURL {
		t.Errorf("expected baseURL %q, got %q", baseURL, c.baseURL)
	}
	if c.uploadURL != uploadBaseURL {
		t.Errorf("expected uploadURL %q, got %q", uploadBaseURL, c.uploadURL)
	}
	if c.httpClient == nil {
		t.Fatal("expected non-nil httpClient")
	}
}

func TestNewClientWithHTTP(t *testing.T) {
	hc := &http.Client{}
	c := NewClientWithHTTP("tok", hc, "http://localhost")
	if c.httpClient != hc {
		t.Error("expected custom http client")
	}
	if c.baseURL != "http://localhost" {
		t.Errorf("expected baseURL 'http://localhost', got %q", c.baseURL)
	}
	if c.uploadURL != "http://localhost" {
		t.Errorf("expected uploadURL 'http://localhost', got %q", c.uploadURL)
	}
}

// --- Helper to create a test client against an httptest server ---

func testClient(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	c := NewClientWithHTTP("test-token", srv.Client(), srv.URL)
	return c, srv
}

// --- GET ---

func TestGet_Success(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"ok":true}`))
	})
	defer srv.Close()

	resp, err := c.Get("/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(resp) != `{"ok":true}` {
		t.Errorf("unexpected response: %s", resp)
	}
}

func TestGet_WithParams(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("foo") != "bar" {
			t.Errorf("expected query param foo=bar, got %q", r.URL.Query().Get("foo"))
		}
		if r.URL.Query().Get("baz") != "qux" {
			t.Errorf("expected query param baz=qux, got %q", r.URL.Query().Get("baz"))
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{}`))
	})
	defer srv.Close()

	_, err := c.Get("/test", map[string]string{"foo": "bar", "baz": "qux"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGet_EmptyResponse(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	})
	defer srv.Close()

	resp, err := c.Get("/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != nil {
		t.Errorf("expected nil response for empty body, got %s", resp)
	}
}

// --- POST ---

func TestPost_Success(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"created":true}`))
	})
	defer srv.Close()

	resp, err := c.Post("/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(resp) != `{"created":true}` {
		t.Errorf("unexpected response: %s", resp)
	}
}

func TestPost_JSONBody(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %q", ct)
		}
		body, _ := io.ReadAll(r.Body)
		var m map[string]string
		if err := json.Unmarshal(body, &m); err != nil {
			t.Fatalf("could not unmarshal request body: %v", err)
		}
		if m["name"] != "test" {
			t.Errorf("expected body name=test, got %q", m["name"])
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{}`))
	})
	defer srv.Close()

	_, err := c.Post("/test", map[string]string{"name": "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPost_NilBody_NoContentType(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if ct != "" {
			t.Errorf("expected no Content-Type for nil body, got %q", ct)
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{}`))
	})
	defer srv.Close()

	_, err := c.Post("/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- PUT ---

func TestPut_Success(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"updated":true}`))
	})
	defer srv.Close()

	resp, err := c.Put("/test", map[string]string{"key": "val"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(resp) != `{"updated":true}` {
		t.Errorf("unexpected response: %s", resp)
	}
}

// --- PATCH ---

func TestPatch_Success(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"patched":true}`))
	})
	defer srv.Close()

	resp, err := c.Patch("/test", nil, map[string]string{"key": "val"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(resp) != `{"patched":true}` {
		t.Errorf("unexpected response: %s", resp)
	}
}

func TestPatch_WithParams(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("updateMask") != "price" {
			t.Errorf("expected updateMask=price, got %q", r.URL.Query().Get("updateMask"))
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{}`))
	})
	defer srv.Close()

	_, err := c.Patch("/test", map[string]string{"updateMask": "price"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- DELETE ---

func TestDelete_Success(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(204)
	})
	defer srv.Close()

	err := c.Delete("/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- Error responses ---

func TestGet_APIError(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(`{"error":{"code":400,"message":"bad request","status":"INVALID_ARGUMENT"}}`))
	})
	defer srv.Close()

	_, err := c.Get("/test", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", apiErr.StatusCode)
	}
	if apiErr.Message != "bad request" {
		t.Errorf("expected message 'bad request', got %q", apiErr.Message)
	}
}

func TestGet_401_Unauthorized(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		_, _ = w.Write([]byte(`{"error":{"code":401,"message":"unauthorized"}}`))
	})
	defer srv.Close()

	_, err := c.Get("/test", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr := err.(*APIError)
	if apiErr.StatusCode != 401 {
		t.Errorf("expected 401, got %d", apiErr.StatusCode)
	}
}

func TestGet_403_Forbidden(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(`{"error":{"code":403,"message":"forbidden"}}`))
	})
	defer srv.Close()

	_, err := c.Get("/test", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr := err.(*APIError)
	if apiErr.StatusCode != 403 {
		t.Errorf("expected 403, got %d", apiErr.StatusCode)
	}
}

func TestGet_404_NotFound(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"not found"}}`))
	})
	defer srv.Close()

	_, err := c.Get("/test", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr := err.(*APIError)
	if apiErr.StatusCode != 404 {
		t.Errorf("expected 404, got %d", apiErr.StatusCode)
	}
}

func TestGet_409_Conflict(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(409)
		_, _ = w.Write([]byte(`{"error":{"code":409,"message":"conflict"}}`))
	})
	defer srv.Close()

	_, err := c.Get("/test", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr := err.(*APIError)
	if apiErr.StatusCode != 409 {
		t.Errorf("expected 409, got %d", apiErr.StatusCode)
	}
}

func TestGet_429_TooManyRequests_Retries(t *testing.T) {
	attempts := 0
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts <= 3 {
			w.WriteHeader(429)
			_, _ = w.Write([]byte(`{"error":{"code":429,"message":"rate limited"}}`))
			return
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"ok":true}`))
	})
	defer srv.Close()

	resp, err := c.Get("/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(resp) != `{"ok":true}` {
		t.Errorf("unexpected response: %s", resp)
	}
	if attempts != 4 {
		t.Errorf("expected 4 attempts (1 initial + 3 retries), got %d", attempts)
	}
}

func TestGet_500_ServerError_Retries(t *testing.T) {
	attempts := 0
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts <= 2 {
			w.WriteHeader(500)
			_, _ = w.Write([]byte(`{"error":{"code":500,"message":"internal"}}`))
			return
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"ok":true}`))
	})
	defer srv.Close()

	resp, err := c.Get("/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(resp) != `{"ok":true}` {
		t.Errorf("unexpected response: %s", resp)
	}
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestGet_502_BadGateway_Retries(t *testing.T) {
	attempts := 0
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(502)
		_, _ = w.Write([]byte(`{"error":{"code":502,"message":"bad gateway"}}`))
	})
	defer srv.Close()

	_, err := c.Get("/test", nil)
	if err == nil {
		t.Fatal("expected error after all retries exhausted")
	}
	// 1 initial + 3 retries = 4 attempts
	if attempts != 4 {
		t.Errorf("expected 4 attempts, got %d", attempts)
	}
}

func TestGet_503_ServiceUnavailable_Retries(t *testing.T) {
	attempts := 0
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(503)
		_, _ = w.Write([]byte(`{"error":{"code":503,"message":"unavailable"}}`))
	})
	defer srv.Close()

	_, err := c.Get("/test", nil)
	if err == nil {
		t.Fatal("expected error after all retries exhausted")
	}
	if attempts != 4 {
		t.Errorf("expected 4 attempts, got %d", attempts)
	}
}

func TestPost_APIError(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(422)
		_, _ = w.Write([]byte(`{"error":{"code":422,"message":"unprocessable"}}`))
	})
	defer srv.Close()

	_, err := c.Post("/test", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr := err.(*APIError)
	if apiErr.StatusCode != 422 {
		t.Errorf("expected 422, got %d", apiErr.StatusCode)
	}
}

func TestDelete_APIError(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"not found"}}`))
	})
	defer srv.Close()

	err := c.Delete("/test")
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr := err.(*APIError)
	if apiErr.StatusCode != 404 {
		t.Errorf("expected 404, got %d", apiErr.StatusCode)
	}
}

// --- APIError.Error() ---

func TestAPIError_Error_WithMessage(t *testing.T) {
	e := &APIError{StatusCode: 400, Message: "bad request"}
	expected := "API error 400: bad request"
	if e.Error() != expected {
		t.Errorf("expected %q, got %q", expected, e.Error())
	}
}

func TestAPIError_Error_WithoutMessage(t *testing.T) {
	e := &APIError{StatusCode: 500}
	expected := "API error 500"
	if e.Error() != expected {
		t.Errorf("expected %q, got %q", expected, e.Error())
	}
}

// --- Upload ---

func TestUpload_Success(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "multipart/form-data") {
			t.Errorf("expected multipart/form-data Content-Type, got %q", ct)
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"versionCode":42}`))
	})
	defer srv.Close()

	// Create a temp file to upload.
	tmp := filepath.Join(t.TempDir(), "test.apk")
	if err := os.WriteFile(tmp, []byte("fake-apk-content"), 0644); err != nil {
		t.Fatal(err)
	}

	resp, err := c.Upload("/upload", tmp, "application/octet-stream")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(resp) != `{"versionCode":42}` {
		t.Errorf("unexpected response: %s", resp)
	}
}

func TestUpload_FileNotFound(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("should not reach server")
	})
	defer srv.Close()

	_, err := c.Upload("/upload", "/nonexistent/file.apk", "application/octet-stream")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if !strings.Contains(err.Error(), "could not open file") {
		t.Errorf("expected 'could not open file' error, got %q", err.Error())
	}
}

func TestUpload_APIError(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(`{"error":{"code":400,"message":"invalid APK"}}`))
	})
	defer srv.Close()

	tmp := filepath.Join(t.TempDir(), "test.apk")
	_ = os.WriteFile(tmp, []byte("bad-apk"), 0644)

	_, err := c.Upload("/upload", tmp, "application/octet-stream")
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 400 {
		t.Errorf("expected 400, got %d", apiErr.StatusCode)
	}
}

// --- DownloadToFile ---

func TestDownloadToFile_Success(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte("file-content-here"))
	})
	defer srv.Close()

	dest := filepath.Join(t.TempDir(), "downloaded.bin")
	err := c.DownloadToFile("/download", dest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("could not read downloaded file: %v", err)
	}
	if string(data) != "file-content-here" {
		t.Errorf("unexpected file content: %s", data)
	}
}

func TestDownloadToFile_APIError(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(`{"error":{"code":403,"message":"forbidden"}}`))
	})
	defer srv.Close()

	dest := filepath.Join(t.TempDir(), "downloaded.bin")
	err := c.DownloadToFile("/download", dest)
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 403 {
		t.Errorf("expected 403, got %d", apiErr.StatusCode)
	}
}

// --- Headers ---

func TestUserAgent_Header(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		ua := r.Header.Get("User-Agent")
		if !strings.HasPrefix(ua, "gpc-cli/") {
			t.Errorf("expected User-Agent starting with 'gpc-cli/', got %q", ua)
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{}`))
	})
	defer srv.Close()

	_, _ = c.Get("/test", nil)
}

func TestAuthorization_Header(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		expected := "Bearer test-token"
		if auth != expected {
			t.Errorf("expected Authorization %q, got %q", expected, auth)
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{}`))
	})
	defer srv.Close()

	_, _ = c.Get("/test", nil)
}

// --- Path builders ---

func TestAppsPath(t *testing.T) {
	expected := "/applications/com.example.app"
	if got := AppsPath("com.example.app"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestEditsPath(t *testing.T) {
	expected := "/applications/com.example/edits/edit123"
	if got := EditsPath("com.example", "edit123"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestNewEditPath(t *testing.T) {
	expected := "/applications/com.example/edits"
	if got := NewEditPath("com.example"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestTracksPath(t *testing.T) {
	expected := "/applications/com.example/edits/e1/tracks"
	if got := TracksPath("com.example", "e1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestTrackPath(t *testing.T) {
	expected := "/applications/com.example/edits/e1/tracks/production"
	if got := TrackPath("com.example", "e1", "production"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestListingsPath(t *testing.T) {
	expected := "/applications/com.example/edits/e1/listings"
	if got := ListingsPath("com.example", "e1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestListingPath(t *testing.T) {
	expected := "/applications/com.example/edits/e1/listings/en-US"
	if got := ListingPath("com.example", "e1", "en-US"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestImagesPath(t *testing.T) {
	expected := "/applications/com.example/edits/e1/listings/en-US/phoneScreenshots"
	if got := ImagesPath("com.example", "e1", "en-US", "phoneScreenshots"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestImagePath(t *testing.T) {
	expected := "/applications/com.example/edits/e1/listings/en-US/phoneScreenshots/img1"
	if got := ImagePath("com.example", "e1", "en-US", "phoneScreenshots", "img1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestDetailsPath(t *testing.T) {
	expected := "/applications/com.example/edits/e1/details"
	if got := DetailsPath("com.example", "e1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestTestersPath(t *testing.T) {
	expected := "/applications/com.example/edits/e1/testers/alpha"
	if got := TestersPath("com.example", "e1", "alpha"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestAPKsPath(t *testing.T) {
	expected := "/applications/com.example/edits/e1/apks"
	if got := APKsPath("com.example", "e1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestBundlesPath(t *testing.T) {
	expected := "/applications/com.example/edits/e1/bundles"
	if got := BundlesPath("com.example", "e1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestDeobfuscationFilesPath(t *testing.T) {
	expected := "/applications/com.example/edits/e1/apks/42/deobfuscationFiles/proguard"
	if got := DeobfuscationFilesPath("com.example", "e1", 42, "proguard"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestExpansionFilesPath(t *testing.T) {
	expected := "/applications/com.example/edits/e1/apks/42/expansionFiles/main"
	if got := ExpansionFilesPath("com.example", "e1", 42, "main"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestCountryAvailabilityPath(t *testing.T) {
	expected := "/applications/com.example/edits/e1/countryAvailability/production"
	if got := CountryAvailabilityPath("com.example", "e1", "production"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestInAppProductsPath(t *testing.T) {
	expected := "/applications/com.example/inappproducts"
	if got := InAppProductsPath("com.example"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestInAppProductPath(t *testing.T) {
	expected := "/applications/com.example/inappproducts/sku123"
	if got := InAppProductPath("com.example", "sku123"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestSubscriptionsPath(t *testing.T) {
	expected := "/applications/com.example/subscriptions"
	if got := SubscriptionsPath("com.example"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestSubscriptionPath(t *testing.T) {
	expected := "/applications/com.example/subscriptions/sub1"
	if got := SubscriptionPath("com.example", "sub1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestBasePlansPath(t *testing.T) {
	expected := "/applications/com.example/subscriptions/sub1/basePlans"
	if got := BasePlansPath("com.example", "sub1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestBasePlanPath(t *testing.T) {
	expected := "/applications/com.example/subscriptions/sub1/basePlans/bp1"
	if got := BasePlanPath("com.example", "sub1", "bp1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestOffersPath(t *testing.T) {
	expected := "/applications/com.example/subscriptions/sub1/basePlans/bp1/offers"
	if got := OffersPath("com.example", "sub1", "bp1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestOfferPath(t *testing.T) {
	expected := "/applications/com.example/subscriptions/sub1/basePlans/bp1/offers/offer1"
	if got := OfferPath("com.example", "sub1", "bp1", "offer1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestOneTimeProductsPath(t *testing.T) {
	expected := "/applications/com.example/oneTimeProducts"
	if got := OneTimeProductsPath("com.example"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestOneTimeProductPath(t *testing.T) {
	expected := "/applications/com.example/oneTimeProducts/prod1"
	if got := OneTimeProductPath("com.example", "prod1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestReviewsPath(t *testing.T) {
	expected := "/applications/com.example/reviews"
	if got := ReviewsPath("com.example"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestReviewPath(t *testing.T) {
	expected := "/applications/com.example/reviews/rev1"
	if got := ReviewPath("com.example", "rev1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestOrdersPath(t *testing.T) {
	expected := "/applications/com.example/orders"
	if got := OrdersPath("com.example"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestOrderPath(t *testing.T) {
	expected := "/applications/com.example/orders/order1"
	if got := OrderPath("com.example", "order1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestPurchaseProductPath(t *testing.T) {
	expected := "/applications/com.example/purchases/products/prod1/tokens/tok1"
	if got := PurchaseProductPath("com.example", "prod1", "tok1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestPurchaseSubscriptionPath(t *testing.T) {
	expected := "/applications/com.example/purchases/subscriptions/sub1/tokens/tok1"
	if got := PurchaseSubscriptionPath("com.example", "sub1", "tok1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestVoidedPurchasesPath(t *testing.T) {
	expected := "/applications/com.example/purchases/voidedpurchases"
	if got := VoidedPurchasesPath("com.example"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestDeviceTierConfigsPath(t *testing.T) {
	expected := "/applications/com.example/deviceTierConfigs"
	if got := DeviceTierConfigsPath("com.example"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestAppRecoveriesPath(t *testing.T) {
	expected := "/applications/com.example/appRecoveries"
	if got := AppRecoveriesPath("com.example"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestExternalTransactionsPath(t *testing.T) {
	expected := "/applications/com.example/externalTransactions"
	if got := ExternalTransactionsPath("com.example"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestGeneratedAPKsPath(t *testing.T) {
	expected := "/applications/com.example/generatedApks/42"
	if got := GeneratedAPKsPath("com.example", 42); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestSystemAPKVariantsPath(t *testing.T) {
	expected := "/applications/com.example/systemApks/42/variants"
	if got := SystemAPKVariantsPath("com.example", 42); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestInternalSharingAPKPath(t *testing.T) {
	expected := "/applications/internalappsharing/com.example/artifacts/apk"
	if got := InternalSharingAPKPath("com.example"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestDataSafetyPath(t *testing.T) {
	expected := "/applications/com.example/dataSafety"
	if got := DataSafetyPath("com.example"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestPricingConvertPath(t *testing.T) {
	expected := "/applications/com.example/pricing:convertRegionPrices"
	if got := PricingConvertPath("com.example"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestUsersPath(t *testing.T) {
	expected := "/developers/dev1/users"
	if got := UsersPath("dev1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestGrantsPath(t *testing.T) {
	expected := "/developers/dev1/users/user1/grants"
	if got := GrantsPath("dev1", "user1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestReleasesListPath(t *testing.T) {
	expected := "/applications/com.example/tracks/production/releases"
	if got := ReleasesListPath("com.example", "production"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestDeviceTierConfigPath(t *testing.T) {
	expected := "/applications/com.example/deviceTierConfigs/cfg1"
	if got := DeviceTierConfigPath("com.example", "cfg1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestAppRecoveryPath(t *testing.T) {
	expected := "/applications/com.example/appRecoveries/rec1"
	if got := AppRecoveryPath("com.example", "rec1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestExternalTransactionPath(t *testing.T) {
	expected := "/applications/com.example/externalTransactions/txn1"
	if got := ExternalTransactionPath("com.example", "txn1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestGeneratedAPKDownloadPath(t *testing.T) {
	expected := "/applications/com.example/generatedApks/42/downloads/dl1:download"
	if got := GeneratedAPKDownloadPath("com.example", 42, "dl1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestSystemAPKVariantPath(t *testing.T) {
	expected := "/applications/com.example/systemApks/42/variants/v1"
	if got := SystemAPKVariantPath("com.example", 42, "v1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestInternalSharingBundlePath(t *testing.T) {
	expected := "/applications/internalappsharing/com.example/artifacts/bundle"
	if got := InternalSharingBundlePath("com.example"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestUserPath(t *testing.T) {
	expected := "/developers/dev1/users/user1"
	if got := UserPath("dev1", "user1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestGrantPath(t *testing.T) {
	expected := "/developers/dev1/users/user1/grants/grant1"
	if got := GrantPath("dev1", "user1", "grant1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestPurchaseProductV2Path(t *testing.T) {
	expected := "/applications/com.example/purchases/productsv2/tokens/tok1"
	if got := PurchaseProductV2Path("com.example", "tok1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestPurchaseSubscriptionV2Path(t *testing.T) {
	expected := "/applications/com.example/purchases/subscriptionsv2/tokens/tok1"
	if got := PurchaseSubscriptionV2Path("com.example", "tok1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestPurchaseOptionsPath(t *testing.T) {
	expected := "/applications/com.example/oneTimeProducts/prod1/purchaseOptions"
	if got := PurchaseOptionsPath("com.example", "prod1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestOTPOffersPath(t *testing.T) {
	expected := "/applications/com.example/oneTimeProducts/prod1/purchaseOptions/po1/offers"
	if got := OTPOffersPath("com.example", "prod1", "po1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestOTPOfferPath(t *testing.T) {
	expected := "/applications/com.example/oneTimeProducts/prod1/purchaseOptions/po1/offers/off1"
	if got := OTPOfferPath("com.example", "prod1", "po1", "off1"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
