package hints

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var ErrMLHintServiceUnavailable = errors.New("ml hint service unavailable")

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) GetHint(ctx context.Context, req MLHintRequest) (*MLHintResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := c.baseURL + "/api/v1/hints"

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewReader(payload),
	)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMLHintServiceUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf(
			"%w: status code %d",
			ErrMLHintServiceUnavailable,
			resp.StatusCode,
		)
	}

	var result MLHintResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	result.Hint = strings.TrimSpace(result.Hint)
	if result.Hint == "" {
		return nil, errors.New("ml hint service returned empty hint")
	}

	if strings.TrimSpace(result.Source) == "" {
		result.Source = "ml"
	}

	return &result, nil
}
