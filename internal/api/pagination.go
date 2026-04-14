package api

import (
	"encoding/json"
	"fmt"
)

// PagedResponse represents an API response that may have additional pages.
type PagedResponse struct {
	NextPageToken string          `json:"nextPageToken,omitempty"`
	Data          json.RawMessage `json:"-"`
}

// ListAll follows pagination tokens to retrieve all pages from a list endpoint.
// The mergeFn is called with each page's raw response; it should extract and
// accumulate items. Returns when there are no more pages.
func (c *Client) ListAll(path string, params map[string]string, mergeFn func(json.RawMessage) error) error {
	if params == nil {
		params = make(map[string]string)
	}

	for {
		resp, err := c.Get(path, params)
		if err != nil {
			return err
		}

		if err := mergeFn(resp); err != nil {
			return fmt.Errorf("could not process page: %w", err)
		}

		// Check for next page token.
		var page struct {
			NextPageToken string `json:"nextPageToken"`
		}
		if err := json.Unmarshal(resp, &page); err != nil {
			// If we can't parse pagination, assume single page.
			return nil
		}

		if page.NextPageToken == "" {
			return nil
		}

		params["pageToken"] = page.NextPageToken
	}
}

// ListAllRaw collects all raw page responses into a slice.
func (c *Client) ListAllRaw(path string, params map[string]string) ([]json.RawMessage, error) {
	var pages []json.RawMessage
	err := c.ListAll(path, params, func(page json.RawMessage) error {
		pages = append(pages, page)
		return nil
	})
	return pages, err
}
