package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

// Client represents a client for the OpenLabs API.
type Client struct {
	BaseURL     string
	AuthToken   string
	EncKey      string
	HTTPClient  *http.Client
	CookieJar   http.CookieJar
	LastCookies []*http.Cookie
}

// NewClient creates a new OpenLabs API client.
func NewClient() *Client {
	// Create a cookie jar to store cookies between requests
	jar, _ := cookiejar.New(nil)

	// Create HTTP client with cookie jar
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Jar:     jar,
	}

	// Load current config to ensure we have the latest tokens
	config, err := loadConfig()
	currentAuthToken := AuthToken
	currentEncKey := EncKey

	if err == nil {
		// Prefer config values over global variables
		if config.AuthToken != "" {
			currentAuthToken = config.AuthToken
		}
		if config.EncKey != "" {
			currentEncKey = config.EncKey
		}
	}

	if Debug {
		fmt.Printf("DEBUG: Creating new client with auth token length: %d\n", len(currentAuthToken))
		fmt.Printf("DEBUG: Creating new client with enc key length: %d\n", len(currentEncKey))
	}

	return &Client{
		BaseURL:    APIURL,
		AuthToken:  currentAuthToken,
		EncKey:     currentEncKey,
		HTTPClient: httpClient,
		CookieJar:  jar,
	}
}

// DoRequest performs an HTTP request to the OpenLabs API.
func (c *Client) DoRequest(method, path string, body interface{}) (*http.Response, error) {
	requestURL := fmt.Sprintf("%s%s", c.BaseURL, path)

	var reqBody io.Reader
	var bodyStr string
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %s", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
		bodyStr = string(jsonData)
	}

	req, err := http.NewRequest(method, requestURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// We still need to manually add cookies because HTTP-only cookies from a response won't be accessible to Go
	parsedURL, _ := url.Parse(requestURL)

	// Add access token cookie if available
	if c.AuthToken != "" {
		// Try multiple cookie names to ensure compatibility
		cookieNames := []string{"access_token_cookie", "jwt", "token", "auth_token", "access_token"}

		for _, name := range cookieNames {
			authCookie := &http.Cookie{
				Name:   name,
				Value:  c.AuthToken,
				Path:   "/",
				Domain: parsedURL.Hostname(),
				// For local testing, we may need to make these false
				HttpOnly: false, // Should be true in production
				Secure:   false, // Should be true for HTTPS
			}
			req.AddCookie(authCookie)
		}

		if Debug {
			fmt.Printf("DEBUG: Added auth cookies with token: %s\n", c.AuthToken)
			fmt.Printf("DEBUG: Token length: %d\n", len(c.AuthToken))
		}
	} else {
		if Debug {
			fmt.Println("DEBUG: No auth token available for cookies")
		}
	}

	// Add encryption key cookie if available
	if c.EncKey != "" {
		encKeyCookie := &http.Cookie{
			Name:     "enc_key",
			Value:    c.EncKey,
			Path:     "/",
			Domain:   parsedURL.Hostname(),
			HttpOnly: false,
			Secure:   false,
		}
		req.AddCookie(encKeyCookie)

		if Debug {
			fmt.Printf("DEBUG: Added enc_key cookie with value: %s\n", c.EncKey)
		}
	}

	// Try a fallback to Bearer token header as well
	if c.AuthToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AuthToken))

		if Debug {
			fmt.Println("DEBUG: Also added fallback Authorization header")
		}
	}

	// Debug cookies
	if Debug {
		parsedURL, _ := url.Parse(requestURL)
		if c.CookieJar != nil {
			cookies := c.CookieJar.Cookies(parsedURL)
			if len(cookies) > 0 {
				fmt.Printf("\n--- DEBUG: COOKIES BEING SENT ---\n")
				for _, cookie := range cookies {
					fmt.Printf("  %s: %s\n", cookie.Name, cookie.Value)
				}
				fmt.Printf("------------------------------\n")
			}
		}
	}

	// Print Debug information if enabled
	if Debug {
		fmt.Printf("\n--- DEBUG: REQUEST ---\n")
		fmt.Printf("URL: %s %s\n", method, requestURL)
		fmt.Printf("Headers:\n")
		for key, values := range req.Header {
			fmt.Printf("  %s: %s\n", key, strings.Join(values, ", "))
		}
		if body != nil {
			fmt.Printf("Body: %s\n", bodyStr)
		}
		fmt.Printf("---------------------\n")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %s", err)
	}

	// Store cookies for later access
	c.LastCookies = resp.Cookies()

	// Print Debug information for response if enabled
	if Debug {
		fmt.Printf("\n--- DEBUG: RESPONSE ---\n")
		fmt.Printf("Status: %s\n", resp.Status)
		fmt.Printf("Headers:\n")
		for key, values := range resp.Header {
			fmt.Printf("  %s: %s\n", key, strings.Join(values, ", "))
		}
		fmt.Printf("Cookies:\n")
		for _, cookie := range resp.Cookies() {
			fmt.Printf("  %s: %s\n", cookie.Name, cookie.Value)
		}

		// Don't read the body here as it will consume the reader
		// Instead, we'll Debug output in the ParseResponse function
		fmt.Printf("----------------------\n")
	}

	return resp, nil
}

// ParseResponse parses the response body into the provided struct.
func ParseResponse(resp *http.Response, result interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %s", err)
	}

	// Debug output for response body
	if Debug && len(body) > 0 {
		fmt.Printf("\n--- DEBUG: RESPONSE BODY ---\n")
		// Try to pretty print JSON if possible
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, body, "", "  "); err == nil {
			fmt.Println(prettyJSON.String())
		} else {
			// If not valid JSON, print as string
			fmt.Println(string(body))
		}
		fmt.Printf("---------------------------\n")
	}

	if resp.StatusCode != http.StatusOK {
		// Try to extract error message from response body if it's JSON
		var errorResponse map[string]interface{}
		if len(body) > 0 && json.Unmarshal(body, &errorResponse) == nil {
			if detail, ok := errorResponse["detail"]; ok {
				return fmt.Errorf("request failed with status: %s - %v", resp.Status, detail)
			}
		}
		return fmt.Errorf("request failed with status: %s", resp.Status)
	}

	if len(body) == 0 {
		return nil
	}

	if result != nil {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %s", err)
		}
	}

	return nil
}

// FormatResponse formats the response as pretty JSON.
func FormatResponse(data interface{}) (string, error) {
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %s", err)
	}
	return string(prettyJSON), nil
}

// GetCookiesForURL retrieves all cookies for a URL.
func (c *Client) GetCookiesForURL(urlStr string) []*http.Cookie {
	if c.CookieJar == nil {
		return nil
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil
	}

	return c.CookieJar.Cookies(parsedURL)
}
