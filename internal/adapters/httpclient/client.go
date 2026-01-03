package httpclient

import (
	"context"
	"net/http"
	"time"

	"github.com/ChangChunCheng/webprism/internal/infrastructure/logger"
)

// Client implements the HTTPClient port.
type Client struct {
	httpClient *http.Client
	logger     *logger.Logger
}

// NewClient creates a new HTTP client.
func NewClient(log *logger.Logger) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		logger: log,
	}
}

// Do executes an HTTP request and returns the response.
func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	// Add context to request
	req = req.WithContext(ctx)

	// Log request
	c.logger.Debug("executing HTTP request",
		logger.String("method", req.Method),
		logger.String("url", req.URL.String()),
	)

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("HTTP request failed",
			logger.String("method", req.Method),
			logger.String("url", req.URL.String()),
			logger.Any("error", err),
		)
		return nil, err
	}

	// Log response
	c.logger.Debug("HTTP request completed",
		logger.String("method", req.Method),
		logger.String("url", req.URL.String()),
		logger.Int("status_code", resp.StatusCode),
	)

	return resp, nil
}
