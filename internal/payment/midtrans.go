// Package payment
package payment

import (
	"context"
	"fmt"

	"github.com/lelecodedev/villa-backend/internal/config"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransGateway struct {
	snapClient snap.Client
	coreClient coreapi.Client
	serverKey  string
}

func NewMidtransGateway() *MidtransGateway {
	env := midtrans.Sandbox
	if config.Env.MidtransIsProduction {
		env = midtrans.Production
	}

	serverKey := config.Env.MidtransServerKey

	gateway := &MidtransGateway{serverKey: serverKey}
	gateway.snapClient.New(serverKey, env)
	gateway.coreClient.New(serverKey, env)

	return gateway
}

func (g *MidtransGateway) CreateSnapToken(orderID string, amount int64, items []midtrans.ItemDetails, customer *midtrans.CustomerDetails, expiryMinutes int) (string, string, error) {
	res, err := g.snapClient.CreateTransaction(&snap.Request{
		TransactionDetails: midtrans.TransactionDetails{OrderID: orderID, GrossAmt: amount},
		Items:              &items,
		CustomerDetail:     customer,
		Expiry: &snap.ExpiryDetails{
			Unit:     "minute",
			Duration: int64(expiryMinutes),
		},
	})
	if err != nil {
		return "", "", err
	}

	return res.Token, res.RedirectURL, nil
}

func (g *MidtransGateway) CheckTransaction(ctx context.Context, orderID string) (*coreapi.TransactionStatusResponse, error) {
	res, err := g.coreClient.CheckTransaction(orderID)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != "200" && res.StatusCode != "201" && res.StatusCode != "407" {
		return nil, fmt.Errorf("midtrans check: status_code=%s %s", res.StatusCode, res.StatusMessage)
	}

	return res, nil
}
