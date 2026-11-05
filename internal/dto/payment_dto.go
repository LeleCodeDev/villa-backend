package dto

type MidtransNotification struct {
	OrderID           string `form:"order_id" json:"order_id"`
	StatusCode        string `form:"status_code" json:"status_code"`
	GrossAmount       string `form:"gross_amount" json:"gross_amount"`
	SignatureKey      string `form:"signature_key" json:"signature_key"`
	TransactionStatus string `form:"transaction_status" json:"transaction_status"`
	FraudStatus       string `form:"fraud_status" json:"fraud_status"`
	PaymentType       string `form:"payment_type" json:"payment_type"`
}
