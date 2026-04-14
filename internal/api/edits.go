package api

import (
	"encoding/json"
	"fmt"
)

// Edit represents an edit session.
type Edit struct {
	ID             string `json:"id"`
	ExpiryTimeSeconds string `json:"expiryTimeSeconds,omitempty"`
}

// CreateEdit creates a new edit session for the given package.
func (c *Client) CreateEdit(pkg string) (*Edit, error) {
	resp, err := c.Post(NewEditPath(pkg), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create edit: %w", err)
	}

	var edit Edit
	if err := json.Unmarshal(resp, &edit); err != nil {
		return nil, fmt.Errorf("could not parse edit response: %w", err)
	}
	return &edit, nil
}

// GetEdit retrieves an existing edit session.
func (c *Client) GetEdit(pkg, editID string) (*Edit, error) {
	resp, err := c.Get(EditsPath(pkg, editID), nil)
	if err != nil {
		return nil, fmt.Errorf("could not get edit: %w", err)
	}

	var edit Edit
	if err := json.Unmarshal(resp, &edit); err != nil {
		return nil, fmt.Errorf("could not parse edit response: %w", err)
	}
	return &edit, nil
}

// ValidateEdit validates an edit session without committing.
func (c *Client) ValidateEdit(pkg, editID string) (*Edit, error) {
	path := EditsPath(pkg, editID) + ":validate"
	resp, err := c.Post(path, nil)
	if err != nil {
		return nil, fmt.Errorf("could not validate edit: %w", err)
	}

	var edit Edit
	if err := json.Unmarshal(resp, &edit); err != nil {
		return nil, fmt.Errorf("could not parse edit response: %w", err)
	}
	return &edit, nil
}

// CommitEdit commits an edit session, applying all changes.
func (c *Client) CommitEdit(pkg, editID string) (*Edit, error) {
	path := EditsPath(pkg, editID) + ":commit"
	resp, err := c.Post(path, nil)
	if err != nil {
		return nil, fmt.Errorf("could not commit edit: %w", err)
	}

	var edit Edit
	if err := json.Unmarshal(resp, &edit); err != nil {
		return nil, fmt.Errorf("could not parse edit response: %w", err)
	}
	return &edit, nil
}

// DeleteEdit deletes an edit session, discarding all changes.
func (c *Client) DeleteEdit(pkg, editID string) error {
	return c.Delete(EditsPath(pkg, editID))
}

// WithEdit creates an edit session, executes the provided function within it,
// and commits on success or deletes on error. This is the recommended way to
// make changes that require an edit session.
func (c *Client) WithEdit(pkg string, fn func(editID string) error) (string, error) {
	edit, err := c.CreateEdit(pkg)
	if err != nil {
		return "", err
	}

	if err := fn(edit.ID); err != nil {
		// Best-effort cleanup.
		_ = c.DeleteEdit(pkg, edit.ID)
		return "", err
	}

	committed, err := c.CommitEdit(pkg, edit.ID)
	if err != nil {
		_ = c.DeleteEdit(pkg, edit.ID)
		return "", fmt.Errorf("could not commit edit: %w", err)
	}

	return committed.ID, nil
}
