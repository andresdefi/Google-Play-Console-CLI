package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/andresdefi/gpc/internal/config"
	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2/google"
)

const (
	keyringService = "gpc-cli"
	keyringUser    = "oauth-token"

	// Scope required for Google Play Developer API.
	androidPublisherScope = "https://www.googleapis.com/auth/androidpublisher"
)

// ServiceAccountInfo holds the parsed service account identity.
type ServiceAccountInfo struct {
	ClientEmail string `json:"client_email"`
	ProjectID   string `json:"project_id"`
}

// Login stores the service account key file path and validates it can produce a token.
func Login(keyFilePath string) (*ServiceAccountInfo, error) {
	data, err := os.ReadFile(keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not read key file: %w", err)
	}

	// Validate the key file by attempting to create credentials.
	creds, err := google.CredentialsFromJSONWithType(context.Background(), data, google.ServiceAccount, androidPublisherScope)
	if err != nil {
		return nil, fmt.Errorf("invalid service account key file: %w", err)
	}

	// Verify we can get a token.
	_, err = creds.TokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("could not obtain access token: %w", err)
	}

	var info ServiceAccountInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("could not parse service account info: %w", err)
	}

	// Store key file path in config.
	cfg, err := config.Load()
	if err != nil {
		cfg = &config.Config{}
	}
	cfg.KeyFilePath = keyFilePath
	if err := config.Save(cfg); err != nil {
		return nil, fmt.Errorf("could not save config: %w", err)
	}

	// Store the key file content in keyring for secure access.
	if err := keyring.Set(keyringService, keyringUser, string(data)); err != nil {
		// Fallback: key file path is already in config, warn but don't fail.
		fmt.Fprintf(os.Stderr, "Warning: could not store credentials in system keychain: %v\n", err)
		fmt.Fprintf(os.Stderr, "Credentials will be read from the key file at: %s\n", keyFilePath)
	}

	return &info, nil
}

// GetToken returns a valid OAuth2 access token, refreshing if necessary.
func GetToken() (string, error) {
	data, err := getKeyData()
	if err != nil {
		return "", err
	}

	creds, err := google.CredentialsFromJSONWithType(context.Background(), data, google.ServiceAccount, androidPublisherScope)
	if err != nil {
		return "", fmt.Errorf("invalid stored credentials: %w", err)
	}

	tok, err := creds.TokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("could not obtain access token: %w", err)
	}

	return tok.AccessToken, nil
}

// GetStatus returns the stored service account info, or an error if not logged in.
func GetStatus() (*ServiceAccountInfo, error) {
	data, err := getKeyData()
	if err != nil {
		return nil, err
	}

	var info ServiceAccountInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("could not parse stored credentials: %w", err)
	}
	return &info, nil
}

// Logout removes stored credentials from keyring and config.
func Logout() error {
	// Remove from keyring (ignore error if not found).
	_ = keyring.Delete(keyringService, keyringUser)

	// Clear key file path from config.
	cfg, err := config.Load()
	if err != nil {
		return nil
	}
	cfg.KeyFilePath = ""
	return config.Save(cfg)
}

// MaskEmail masks a service account email for display.
func MaskEmail(email string) string {
	parts := strings.SplitN(email, "@", 2)
	if len(parts) != 2 {
		return email
	}
	name := parts[0]
	if len(name) > 8 {
		name = name[:4] + "..." + name[len(name)-4:]
	}
	return name + "@" + parts[1]
}

// getKeyData retrieves the service account key data from keyring or file.
func getKeyData() ([]byte, error) {
	// Try keyring first.
	data, err := keyring.Get(keyringService, keyringUser)
	if err == nil && data != "" {
		return []byte(data), nil
	}

	// Fallback to reading from key file path in config.
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("not authenticated: %w", err)
	}

	if cfg.KeyFilePath == "" {
		return nil, fmt.Errorf("not authenticated; run 'gpc auth login' first")
	}

	fileData, err := os.ReadFile(cfg.KeyFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not read key file at %s: %w", cfg.KeyFilePath, err)
	}

	return fileData, nil
}
