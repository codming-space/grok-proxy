package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"grok-proxy/internal/cookie"
	"grok-proxy/internal/models"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	GrokURL        = "https://grok.com/rest/app-chat/conversations/new"
	DefaultTimeout = 240 * time.Second
)

type GrokClient struct {
	httpClient    *http.Client
	cookieManager *cookie.Manager
}

type JSONToken struct {
	Result struct {
		Response struct {
			Token string `json:"token"`
		} `json:"response"`
	} `json:"result"`
}

// NewGrokClient creates a new client for Grok API
func NewGrokClient(cookieManager *cookie.Manager) *GrokClient {
	return &GrokClient{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		cookieManager: cookieManager,
	}
}

// StreamRequest sends a streaming request to Grok API
func (c *GrokClient) StreamRequest(ctx context.Context, message string, model string) (<-chan string, <-chan error) {
	tokenCh := make(chan string)
	errCh := make(chan error, 1)

	go func() {
		defer close(tokenCh)
		defer close(errCh)

		// Create request
		grokReq := models.GrokRequest{
			Message:   message,
			ModelName: model,
		}
		reqBody, err := json.Marshal(grokReq)
		if err != nil {
			errCh <- fmt.Errorf("error marshaling request: %w", err)
			return
		}

		req, err := http.NewRequestWithContext(ctx, "POST", GrokURL, bytes.NewBuffer(reqBody))
		if err != nil {
			errCh <- fmt.Errorf("error creating request: %w", err)
			return
		}

		// Set headers
		req.Header.Set("Authority", "grok.com")
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", c.cookieManager.GetCookie())
		req.Header.Set("Origin", "https://grok.com")
		req.Header.Set("Referer", "https://grok.com/?referrer=website")
		req.Header.Set("User-Agent", c.cookieManager.GetUserAgent())

		// Send request
		resp, err := c.httpClient.Do(req)
		if err != nil {
			errCh <- fmt.Errorf("error sending request: %w", err)
			return
		}
		defer resp.Body.Close()

		// Check response status
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("Error status: %d, body: %s", resp.StatusCode, body)

			// Rotate cookie and user-agent on error
			c.cookieManager.GetCookie()

			errCh <- fmt.Errorf("error status: %d", resp.StatusCode)
			return
		}

		// Process streaming response
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}

			// Parse JSON response
			var jsonToken JSONToken
			if err := json.Unmarshal([]byte(line), &jsonToken); err != nil {
				log.Printf("JSON parsing error: %v for line: %s", err, line)
				continue
			}

			// Extract token
			token := jsonToken.Result.Response.Token
			if token != "" {
				select {
				case tokenCh <- token:
					// Token sent successfully
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				}
			}
		}

		// Check for scanning errors
		if err := scanner.Err(); err != nil {
			errCh <- fmt.Errorf("error reading response: %w", err)
		}
	}()

	return tokenCh, errCh
}

// Execute sends a non-streaming request and collects all tokens
func (c *GrokClient) Execute(ctx context.Context, message string, model string) (string, error) {
	tokenCh, errCh := c.StreamRequest(ctx, message, model)

	var result strings.Builder
	for {
		select {
		case token, ok := <-tokenCh:
			if !ok {
				// Channel closed, we're done
				return result.String(), nil
			}
			result.WriteString(token)
		case err := <-errCh:
			return result.String(), err
		case <-ctx.Done():
			return result.String(), ctx.Err()
		}
	}
}
