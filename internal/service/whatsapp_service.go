package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	appError "github.com/lelecodedev/villa-backend/internal/apperror"
	"github.com/lelecodedev/villa-backend/internal/client"
)

type WhatsappService struct {
	Client *client.WhatsappClient
}

func NewWhatsappService(client *client.WhatsappClient) *WhatsappService {
	return &WhatsappService{Client: client}
}

func (s *WhatsappService) Send(ctx context.Context, phoneNumber string, message string) error {
	phone := strings.Replace(phoneNumber, "0", "62", 1) + "@c.us"

	if err := s.Client.SendText(ctx, phone, message); err != nil {
		if apiErr, ok := errors.AsType[*appError.APIError](err); ok {
			slog.Error("openwa request failed", "status", apiErr.StatusCode, "body", apiErr.Body)
			return appError.BadGateway("whatsapp provider error")
		}

		return err
	}

	return nil
}
