package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/andresdefi/gpc/internal/version"
)

const (
	baseURL       = "https://androidpublisher.googleapis.com/androidpublisher/v3"
	uploadBaseURL = "https://androidpublisher.googleapis.com/upload/androidpublisher/v3"

	maxRetries    = 3
	baseBackoff   = 100 * time.Millisecond
	clientTimeout = 30 * time.Second
)

// Client is an HTTP client for the Google Play Developer API.
type Client struct {
	httpClient *http.Client
	token      string
	baseURL    string
	uploadURL  string
}

// NewClient creates a new API client with the given OAuth2 access token.
func NewClient(token string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: clientTimeout},
		token:      token,
		baseURL:    baseURL,
		uploadURL:  uploadBaseURL,
	}
}

// NewClientWithHTTP creates a client with a custom http.Client and base URL.
// If httpClient is nil, a default client with timeout is used.
func NewClientWithHTTP(token string, httpClient *http.Client, base string) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: clientTimeout}
	}
	return &Client{
		httpClient: httpClient,
		token:      token,
		baseURL:    base,
		uploadURL:  base,
	}
}

// APIError represents an error returned by the Google Play API.
type APIError struct {
	StatusCode int
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Status     string `json:"status"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("API error %d", e.StatusCode)
}

type errorResponse struct {
	Error *APIError `json:"error"`
}

// Get performs a GET request to the given path with optional query parameters.
func (c *Client) Get(path string, params map[string]string) (json.RawMessage, error) {
	return c.do(http.MethodGet, path, params, nil, false)
}

// Post performs a POST request with a JSON body.
func (c *Client) Post(path string, body any) (json.RawMessage, error) {
	return c.do(http.MethodPost, path, nil, body, false)
}

// Put performs a PUT request with a JSON body.
func (c *Client) Put(path string, body any) (json.RawMessage, error) {
	return c.do(http.MethodPut, path, nil, body, false)
}

// Patch performs a PATCH request with a JSON body.
func (c *Client) Patch(path string, params map[string]string, body any) (json.RawMessage, error) {
	return c.doWithParams(http.MethodPatch, path, params, body, false)
}

// Delete performs a DELETE request.
func (c *Client) Delete(path string) error {
	_, err := c.do(http.MethodDelete, path, nil, nil, false)
	return err
}

// Upload performs a multipart upload request.
func (c *Client) Upload(path string, filePath string, contentType string) (json.RawMessage, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("could not create form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("could not copy file content: %w", err)
	}
	_ = writer.Close()

	fullURL := c.uploadURL + path
	req, err := http.NewRequest(http.MethodPost, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", "gpc-cli/"+version.Version)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("upload request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, parseAPIError(resp.StatusCode, respBody)
	}

	return json.RawMessage(respBody), nil
}

// DownloadToFile performs a GET request and writes the response body to a file.
func (c *Client) DownloadToFile(path string, destPath string) error {
	fullURL := c.baseURL + path
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("User-Agent", "gpc-cli/"+version.Version)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("download request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return parseAPIError(resp.StatusCode, body)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("could not create output file: %w", err)
	}
	defer func() { _ = out.Close() }()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("could not write file: %w", err)
	}

	return nil
}

func (c *Client) do(method, path string, params map[string]string, body any, isUpload bool) (json.RawMessage, error) {
	return c.doWithParams(method, path, params, body, isUpload)
}

func (c *Client) doWithParams(method, path string, params map[string]string, body any, isUpload bool) (json.RawMessage, error) {
	var respBody json.RawMessage
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			backoff := baseBackoff * time.Duration(1<<uint(attempt-1))
			time.Sleep(backoff)
		}

		var reqBody io.Reader
		if body != nil {
			data, err := json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("could not marshal request body: %w", err)
			}
			reqBody = bytes.NewReader(data)
		}

		base := c.baseURL
		if isUpload {
			base = c.uploadURL
		}

		fullURL := base + path
		if len(params) > 0 {
			q := url.Values{}
			for k, v := range params {
				q.Set(k, v)
			}
			fullURL += "?" + q.Encode()
		}

		req, err := http.NewRequest(method, fullURL, reqBody)
		if err != nil {
			return nil, fmt.Errorf("could not create request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+c.token)
		req.Header.Set("User-Agent", "gpc-cli/"+version.Version)
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			if isRetryable(0) {
				continue
			}
			return nil, fmt.Errorf("request failed: %w", err)
		}

		data, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("could not read response: %w", err)
			continue
		}

		if resp.StatusCode >= 400 {
			apiErr := parseAPIError(resp.StatusCode, data)
			if isRetryable(resp.StatusCode) && attempt < maxRetries {
				lastErr = apiErr
				continue
			}
			return nil, apiErr
		}

		if len(data) > 0 {
			respBody = json.RawMessage(data)
		}
		return respBody, nil
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", maxRetries, lastErr)
}

func parseAPIError(statusCode int, body []byte) *APIError {
	var errResp errorResponse
	if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error != nil {
		errResp.Error.StatusCode = statusCode
		return errResp.Error
	}
	return &APIError{
		StatusCode: statusCode,
		Message:    strings.TrimSpace(string(body)),
	}
}

func isRetryable(statusCode int) bool {
	return statusCode == 429 || statusCode == 500 || statusCode == 502 || statusCode == 503 || statusCode == 504
}

// --- Convenience path builders ---

func AppsPath(pkg string) string {
	return "/applications/" + pkg
}

func EditsPath(pkg, editID string) string {
	return "/applications/" + pkg + "/edits/" + editID
}

func NewEditPath(pkg string) string {
	return "/applications/" + pkg + "/edits"
}

func TracksPath(pkg, editID string) string {
	return EditsPath(pkg, editID) + "/tracks"
}

func TrackPath(pkg, editID, track string) string {
	return TracksPath(pkg, editID) + "/" + track
}

func ListingsPath(pkg, editID string) string {
	return EditsPath(pkg, editID) + "/listings"
}

func ListingPath(pkg, editID, lang string) string {
	return ListingsPath(pkg, editID) + "/" + lang
}

func ImagesPath(pkg, editID, lang, imageType string) string {
	return ListingPath(pkg, editID, lang) + "/" + imageType
}

func ImagePath(pkg, editID, lang, imageType, imageID string) string {
	return ImagesPath(pkg, editID, lang, imageType) + "/" + imageID
}

func DetailsPath(pkg, editID string) string {
	return EditsPath(pkg, editID) + "/details"
}

func TestersPath(pkg, editID, track string) string {
	return EditsPath(pkg, editID) + "/testers/" + track
}

func APKsPath(pkg, editID string) string {
	return EditsPath(pkg, editID) + "/apks"
}

func BundlesPath(pkg, editID string) string {
	return EditsPath(pkg, editID) + "/bundles"
}

func DeobfuscationFilesPath(pkg, editID string, versionCode int, fileType string) string {
	return EditsPath(pkg, editID) + "/apks/" + strconv.Itoa(versionCode) + "/deobfuscationFiles/" + fileType
}

func ExpansionFilesPath(pkg, editID string, versionCode int, fileType string) string {
	return EditsPath(pkg, editID) + "/apks/" + strconv.Itoa(versionCode) + "/expansionFiles/" + fileType
}

func CountryAvailabilityPath(pkg, editID, track string) string {
	return EditsPath(pkg, editID) + "/countryAvailability/" + track
}

func InAppProductsPath(pkg string) string {
	return "/applications/" + pkg + "/inappproducts"
}

func InAppProductPath(pkg, sku string) string {
	return InAppProductsPath(pkg) + "/" + sku
}

func SubscriptionsPath(pkg string) string {
	return "/applications/" + pkg + "/subscriptions"
}

func SubscriptionPath(pkg, productID string) string {
	return SubscriptionsPath(pkg) + "/" + productID
}

func BasePlansPath(pkg, productID string) string {
	return SubscriptionPath(pkg, productID) + "/basePlans"
}

func BasePlanPath(pkg, productID, basePlanID string) string {
	return BasePlansPath(pkg, productID) + "/" + basePlanID
}

func OffersPath(pkg, productID, basePlanID string) string {
	return BasePlanPath(pkg, productID, basePlanID) + "/offers"
}

func OfferPath(pkg, productID, basePlanID, offerID string) string {
	return OffersPath(pkg, productID, basePlanID) + "/" + offerID
}

func OneTimeProductsPath(pkg string) string {
	return "/applications/" + pkg + "/oneTimeProducts"
}

func OneTimeProductPath(pkg, productID string) string {
	return OneTimeProductsPath(pkg) + "/" + productID
}

func PurchaseOptionsPath(pkg, productID string) string {
	return OneTimeProductPath(pkg, productID) + "/purchaseOptions"
}

func OTPOffersPath(pkg, productID, purchaseOptionID string) string {
	return PurchaseOptionsPath(pkg, productID) + "/" + purchaseOptionID + "/offers"
}

func OTPOfferPath(pkg, productID, purchaseOptionID, offerID string) string {
	return OTPOffersPath(pkg, productID, purchaseOptionID) + "/" + offerID
}

func ReviewsPath(pkg string) string {
	return "/applications/" + pkg + "/reviews"
}

func ReviewPath(pkg, reviewID string) string {
	return ReviewsPath(pkg) + "/" + reviewID
}

func OrdersPath(pkg string) string {
	return "/applications/" + pkg + "/orders"
}

func OrderPath(pkg, orderID string) string {
	return OrdersPath(pkg) + "/" + orderID
}

func PurchaseProductPath(pkg, productID, token string) string {
	return "/applications/" + pkg + "/purchases/products/" + productID + "/tokens/" + token
}

func PurchaseProductV2Path(pkg, token string) string {
	return "/applications/" + pkg + "/purchases/productsv2/tokens/" + token
}

func PurchaseSubscriptionPath(pkg, subscriptionID, token string) string {
	return "/applications/" + pkg + "/purchases/subscriptions/" + subscriptionID + "/tokens/" + token
}

func PurchaseSubscriptionV2Path(pkg, token string) string {
	return "/applications/" + pkg + "/purchases/subscriptionsv2/tokens/" + token
}

func VoidedPurchasesPath(pkg string) string {
	return "/applications/" + pkg + "/purchases/voidedpurchases"
}

func DeviceTierConfigsPath(pkg string) string {
	return "/applications/" + pkg + "/deviceTierConfigs"
}

func DeviceTierConfigPath(pkg, configID string) string {
	return DeviceTierConfigsPath(pkg) + "/" + configID
}

func AppRecoveriesPath(pkg string) string {
	return "/applications/" + pkg + "/appRecoveries"
}

func AppRecoveryPath(pkg, recoveryID string) string {
	return AppRecoveriesPath(pkg) + "/" + recoveryID
}

func ExternalTransactionsPath(pkg string) string {
	return fmt.Sprintf("/applications/%s/externalTransactions", pkg)
}

func ExternalTransactionPath(pkg, txnID string) string {
	return fmt.Sprintf("/applications/%s/externalTransactions/%s", pkg, txnID)
}

func GeneratedAPKsPath(pkg string, versionCode int) string {
	return fmt.Sprintf("/applications/%s/generatedApks/%d", pkg, versionCode)
}

func GeneratedAPKDownloadPath(pkg string, versionCode int, downloadID string) string {
	return fmt.Sprintf("/applications/%s/generatedApks/%d/downloads/%s:download", pkg, versionCode, downloadID)
}

func SystemAPKVariantsPath(pkg string, versionCode int) string {
	return fmt.Sprintf("/applications/%s/systemApks/%d/variants", pkg, versionCode)
}

func SystemAPKVariantPath(pkg string, versionCode int, variantID string) string {
	return fmt.Sprintf("/applications/%s/systemApks/%d/variants/%s", pkg, versionCode, variantID)
}

func InternalSharingAPKPath(pkg string) string {
	return fmt.Sprintf("/applications/internalappsharing/%s/artifacts/apk", pkg)
}

func InternalSharingBundlePath(pkg string) string {
	return fmt.Sprintf("/applications/internalappsharing/%s/artifacts/bundle", pkg)
}

func DataSafetyPath(pkg string) string {
	return fmt.Sprintf("/applications/%s/dataSafety", pkg)
}

func PricingConvertPath(pkg string) string {
	return fmt.Sprintf("/applications/%s/pricing:convertRegionPrices", pkg)
}

func UsersPath(developerID string) string {
	return fmt.Sprintf("/developers/%s/users", developerID)
}

func UserPath(developerID, userID string) string {
	return fmt.Sprintf("/developers/%s/users/%s", developerID, userID)
}

func GrantsPath(developerID, userID string) string {
	return fmt.Sprintf("/developers/%s/users/%s/grants", developerID, userID)
}

func GrantPath(developerID, userID, grantID string) string {
	return fmt.Sprintf("/developers/%s/users/%s/grants/%s", developerID, userID, grantID)
}

func ReleasesListPath(pkg, track string) string {
	return fmt.Sprintf("/applications/%s/tracks/%s/releases", pkg, track)
}
