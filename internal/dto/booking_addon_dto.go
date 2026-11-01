package dto

import "time"

type (
	BookingAddonResponse struct {
		ID        uint          `json:"id"`
		BookingID uint          `json:"booking_id"`
		AddonID   uint          `json:"addon_id"`
		Addon     AddonResponse `json:"addon"`
		Quantity  int           `json:"quantity"`
		CreatedAt time.Time     `json:"created_at"`
		UpdatedAt time.Time     `json:"updated_at"`
	}
)
