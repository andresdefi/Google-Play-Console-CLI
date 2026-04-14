package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func editsTestClient(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	c := NewClientWithHTTP("test-token", srv.Client(), srv.URL)
	return c, srv
}

func TestCreateEdit_Success(t *testing.T) {
	c, srv := editsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/edits") {
			t.Errorf("expected path ending with /edits, got %s", r.URL.Path)
		}
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(Edit{ID: "edit-123", ExpiryTimeSeconds: "3600"})
	})
	defer srv.Close()

	edit, err := c.CreateEdit("com.example")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if edit.ID != "edit-123" {
		t.Errorf("expected edit ID 'edit-123', got %q", edit.ID)
	}
	if edit.ExpiryTimeSeconds != "3600" {
		t.Errorf("expected expiryTimeSeconds '3600', got %q", edit.ExpiryTimeSeconds)
	}
}

func TestCreateEdit_Error(t *testing.T) {
	c, srv := editsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(`{"error":{"code":403,"message":"forbidden"}}`))
	})
	defer srv.Close()

	_, err := c.CreateEdit("com.example")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "could not create edit") {
		t.Errorf("expected 'could not create edit' in error, got %q", err.Error())
	}
}

func TestGetEdit_Success(t *testing.T) {
	c, srv := editsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(Edit{ID: "edit-456"})
	})
	defer srv.Close()

	edit, err := c.GetEdit("com.example", "edit-456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if edit.ID != "edit-456" {
		t.Errorf("expected edit ID 'edit-456', got %q", edit.ID)
	}
}

func TestValidateEdit_Success(t *testing.T) {
	c, srv := editsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, ":validate") {
			t.Errorf("expected path ending with :validate, got %s", r.URL.Path)
		}
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(Edit{ID: "edit-789"})
	})
	defer srv.Close()

	edit, err := c.ValidateEdit("com.example", "edit-789")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if edit.ID != "edit-789" {
		t.Errorf("expected edit ID 'edit-789', got %q", edit.ID)
	}
}

func TestCommitEdit_Success(t *testing.T) {
	c, srv := editsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, ":commit") {
			t.Errorf("expected path ending with :commit, got %s", r.URL.Path)
		}
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(Edit{ID: "edit-committed"})
	})
	defer srv.Close()

	edit, err := c.CommitEdit("com.example", "edit-committed")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if edit.ID != "edit-committed" {
		t.Errorf("expected edit ID 'edit-committed', got %q", edit.ID)
	}
}

func TestDeleteEdit_Success(t *testing.T) {
	c, srv := editsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(204)
	})
	defer srv.Close()

	err := c.DeleteEdit("com.example", "edit-del")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteEdit_Error(t *testing.T) {
	c, srv := editsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"edit not found"}}`))
	})
	defer srv.Close()

	err := c.DeleteEdit("com.example", "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestWithEdit_Success(t *testing.T) {
	callCount := 0
	c, srv := editsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		switch {
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/edits"):
			// CreateEdit
			w.WriteHeader(200)
			_ = json.NewEncoder(w).Encode(Edit{ID: "e1"})
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, ":commit"):
			// CommitEdit
			w.WriteHeader(200)
			_ = json.NewEncoder(w).Encode(Edit{ID: "e1"})
		default:
			// Any inner calls from fn
			w.WriteHeader(200)
			_, _ = w.Write([]byte(`{}`))
		}
	})
	defer srv.Close()

	fnCalled := false
	editID, err := c.WithEdit("com.example", func(eid string) error {
		fnCalled = true
		if eid != "e1" {
			t.Errorf("expected edit ID 'e1', got %q", eid)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !fnCalled {
		t.Error("expected fn to be called")
	}
	if editID != "e1" {
		t.Errorf("expected returned edit ID 'e1', got %q", editID)
	}
}

func TestWithEdit_FnError(t *testing.T) {
	deleteCalled := false
	c, srv := editsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/edits"):
			w.WriteHeader(200)
			_ = json.NewEncoder(w).Encode(Edit{ID: "e2"})
		case r.Method == http.MethodDelete:
			deleteCalled = true
			w.WriteHeader(204)
		default:
			w.WriteHeader(200)
			_, _ = w.Write([]byte(`{}`))
		}
	})
	defer srv.Close()

	_, err := c.WithEdit("com.example", func(eid string) error {
		return fmt.Errorf("something went wrong")
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "something went wrong") {
		t.Errorf("expected original error, got %q", err.Error())
	}
	if !deleteCalled {
		t.Error("expected DeleteEdit cleanup to be called")
	}
}

func TestWithEdit_CommitError(t *testing.T) {
	deleteCalled := false
	c, srv := editsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/edits") && !strings.Contains(r.URL.Path, ":commit"):
			w.WriteHeader(200)
			_ = json.NewEncoder(w).Encode(Edit{ID: "e3"})
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, ":commit"):
			w.WriteHeader(500)
			_, _ = w.Write([]byte(`{"error":{"code":500,"message":"commit failed"}}`))
		case r.Method == http.MethodDelete:
			deleteCalled = true
			w.WriteHeader(204)
		default:
			w.WriteHeader(200)
			_, _ = w.Write([]byte(`{}`))
		}
	})
	defer srv.Close()

	_, err := c.WithEdit("com.example", func(eid string) error {
		return nil
	})
	if err == nil {
		t.Fatal("expected error from commit failure")
	}
	if !strings.Contains(err.Error(), "could not commit edit") {
		t.Errorf("expected 'could not commit edit' error, got %q", err.Error())
	}
	if !deleteCalled {
		t.Error("expected DeleteEdit cleanup to be called after commit failure")
	}
}

func TestCreateEdit_InvalidJSON(t *testing.T) {
	c, srv := editsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`not json`))
	})
	defer srv.Close()

	_, err := c.CreateEdit("com.example")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "could not parse edit response") {
		t.Errorf("expected parse error, got %q", err.Error())
	}
}
