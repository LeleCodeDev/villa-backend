package dto

import "time"

type (
	RentalOptionQuery struct {
		PaginationQuery
		SortQuery
		Unpage bool `form:"unpage"`
	}

	RentalOptionResponse struct {
		ID          uint      `json:"id"`
		Title       string    `json:"title"`
		Subtitle    string    `json:"subtitle"`
		Description string    `json:"description"`
		RoomCount   int       `json:"room_count"`
		MaxCapacity int       `json:"max_capacity"`
		SortOrder   int       `json:"sort_order"`
		Image       *string   `json:"image"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	RentalOptionRequest struct {
		Title       string  `json:"title" form:"title" binding:"required"`
		Subtitle    string  `json:"subtitle" form:"subtitle" binding:"required"`
		Description string  `json:"description" form:"description" binding:"required"`
		RoomCount   int     `json:"room_count" form:"room_count" binding:"required,gt=0"`
		MaxCapacity int     `json:"max_capacity" form:"max_capacity" binding:"required,gt=0"`
		SortOrder   *int    `json:"sort_order" form:"sort_order" binding:"omitempty,gt=0"`
		Image       *string `json:"image" form:"image" binding:"omitempty"`
	}
)
