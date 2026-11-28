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

type WhatsappCheckResponse struct {
	Number     string `json:"number"`
	Exists     bool   `json:"exists"`
	WhatsappID string `json:"whatsappId"`
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

func (c *WhatsappClient) do(ctx context.Context, method string, path string, payload any, out any) error {
	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}

		body = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, body)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

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

	return c.do(ctx, http.MethodPost, path, payload, nil)
}

func (c *WhatsappClient) CheckNumber(ctx context.Context, number string) (WhatsappCheckResponse, error) {
	var res WhatsappCheckResponse
	path := fmt.Sprintf("/api/sessions/%s/contacts/check/%s", c.SessionID, number)

	err := c.do(ctx, http.MethodGet, path, nil, &res)
	return res, err
}
