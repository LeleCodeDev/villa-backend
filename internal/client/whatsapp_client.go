// Package client
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	appError "github.com/lelecodedev/villa-backend/internal/apperror"
)

type WhatsappClient struct {
	BaseURL   string
	APIKey    string
	SessionID string
	http      *http.Client
}

func NewWhatsappClient(baseURL, APIKey, sessionID string) *WhatsappClient {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   2 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        50,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
	}

	return &WhatsappClient{
		BaseURL:   baseURL,
		APIKey:    APIKey,
		SessionID: sessionID,
		http: &http.Client{
			Transport: transport,
			Timeout:   15 * time.Second,
		},
	}
}

func (c *WhatsappClient) do(ctx context.Context, path string, payload any, out any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+path, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", c.APIKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("call openwa: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &appError.APIError{StatusCode: resp.StatusCode, Body: string(data)}
	}

	if out != nil {
		if err := json.Unmarshal(data, out); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}
	return nil
}

func (c *WhatsappClient) SendText(ctx context.Context, chatID string, message string) error {
	payload := map[string]any{
		"chatId": chatID,
		"text":   message,
	}
	path := fmt.Sprintf("/api/sessions/%s/messages/send-text", c.SessionID)

	return c.do(ctx, path, payload, nil)
}
