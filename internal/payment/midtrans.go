// Package payment
package payment

import (
	"crypto/sha512"
	"crypto/subtle"
	"encoding/hex"

	"github.com/lelecodedev/villa-backend/internal/config"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransGateway struct {
	client    snap.Client
	serverKey string
}

func NewMidtransGateway() *MidtransGateway {
	env := midtrans.Sandbox
	if config.Env.MidtransIsProduction {
		env = midtrans.Production
	}

	serverKey := config.Env.MidtransServerKey

	gateway := &MidtransGateway{serverKey: serverKey}
	gateway.client.New(serverKey, env)

	return gateway
}

func (g *MidtransGateway) CreateSnapToken(orderID string, amount int64, items []midtrans.ItemDetails, customer *midtrans.CustomerDetails, expiryMinutes int) (string, string, error) {
	res, err := g.client.CreateTransaction(&snap.Request{
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

func (g *MidtransGateway) VerifySignature(orderID, statusCode, grossAmount, signature string) bool {
	sum := sha512.Sum512([]byte(orderID + statusCode + grossAmount + g.serverKey))
	return subtle.ConstantTimeCompare([]byte(hex.EncodeToString(sum[:])), []byte(signature)) == 1
}
