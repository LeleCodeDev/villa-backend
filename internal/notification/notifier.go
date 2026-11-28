package notification

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/lelecodedev/villa-backend/internal/client"
	"github.com/lelecodedev/villa-backend/internal/model"
	"github.com/lelecodedev/villa-backend/internal/repository"
)

type Message struct {
	Phone string
	Text  string
}

type Notifier struct {
	client          *client.WhatsappClient
	applicationRepo *repository.ApplicationRepository
}

func NewNotifier(c *client.WhatsappClient, appRepo *repository.ApplicationRepository) *Notifier {
	return &Notifier{client: c, applicationRepo: appRepo}
}

func toChatID(phone string) string {
	return "62" + strings.TrimPrefix(phone, "0") + "@c.us"
}

func (n *Notifier) Send(tag string, build func(ctx context.Context, app *model.Application) []Message) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("notify panic for %s: %v", tag, r)
			}
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		app, err := n.applicationRepo.Get(ctx)
		if err != nil {
			log.Printf("notify %s: get application: %v", tag, err)
			return
		}

		for _, m := range build(ctx, app) {
			if m.Phone == "" {
				continue
			}
			if err := n.client.SendText(ctx, toChatID(m.Phone), m.Text); err != nil {
				log.Printf("notify %s: send to %s: %v", tag, m.Phone, err)
			}
		}
	}()
}
