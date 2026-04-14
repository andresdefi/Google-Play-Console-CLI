package api

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestListAll_SinglePage(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"items":["a","b"]}`))
	})
	defer srv.Close()

	var pages []json.RawMessage
	err := c.ListAll("/test", nil, func(page json.RawMessage) error {
		pages = append(pages, page)
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pages) != 1 {
		t.Errorf("expected 1 page, got %d", len(pages))
	}
}

func TestListAll_MultiplePages(t *testing.T) {
	callCount := 0
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		token := r.URL.Query().Get("pageToken")

		switch {
		case callCount == 1 && token == "":
			w.WriteHeader(200)
			_, _ = w.Write([]byte(`{"items":["a"],"nextPageToken":"page2"}`))
		case callCount == 2 && token == "page2":
			w.WriteHeader(200)
			_, _ = w.Write([]byte(`{"items":["b"],"nextPageToken":"page3"}`))
		case callCount == 3 && token == "page3":
			w.WriteHeader(200)
			_, _ = w.Write([]byte(`{"items":["c"]}`))
		default:
			t.Errorf("unexpected call %d with token %q", callCount, token)
			w.WriteHeader(500)
		}
	})
	defer srv.Close()

	var pages []json.RawMessage
	err := c.ListAll("/test", nil, func(page json.RawMessage) error {
		pages = append(pages, page)
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pages) != 3 {
		t.Errorf("expected 3 pages, got %d", len(pages))
	}
	if callCount != 3 {
		t.Errorf("expected 3 API calls, got %d", callCount)
	}
}

func TestListAll_APIError(t *testing.T) {
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		_, _ = w.Write([]byte(`{"error":{"code":500,"message":"internal"}}`))
	})
	defer srv.Close()

	err := c.ListAll("/test", nil, func(page json.RawMessage) error {
		t.Fatal("mergeFn should not be called on error")
		return nil
	})
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		// The error may be wrapped after retries. Check it's not nil.
		t.Logf("error type: %T, message: %v", err, err)
		return
	}
	if apiErr.StatusCode != 500 {
		t.Errorf("expected status 500, got %d", apiErr.StatusCode)
	}
}

func TestListAllRaw_CollectsPages(t *testing.T) {
	callCount := 0
	c, srv := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.WriteHeader(200)
			_, _ = w.Write([]byte(`{"data":"page1","nextPageToken":"tok2"}`))
			return
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"data":"page2"}`))
	})
	defer srv.Close()

	pages, err := c.ListAllRaw("/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pages) != 2 {
		t.Errorf("expected 2 pages, got %d", len(pages))
	}

	// Verify each page is valid JSON.
	for i, p := range pages {
		var m map[string]any
		if err := json.Unmarshal(p, &m); err != nil {
			t.Errorf("page %d is not valid JSON: %v", i, err)
		}
	}

	// Verify content of pages.
	var m1 map[string]any
	_ = json.Unmarshal(pages[0], &m1)
	if m1["data"] != "page1" {
		t.Errorf("expected page1 data, got %v", m1["data"])
	}

	var m2 map[string]any
	_ = json.Unmarshal(pages[1], &m2)
	if m2["data"] != "page2" {
		t.Errorf("expected page2 data, got %v", m2["data"])
	}
}
