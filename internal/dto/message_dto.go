package dto

import (
	"time"
)

type (
	MessageQuery struct {
		PaginationQuery
		SortQuery
		Unpage bool `form:"unpage"`
	}

	MessageResponse struct {
		ID        uint      `json:"id"`
		Firstname string    `json:"firstname"`
		Lastname  string    `json:"lastname"`
		Email     string    `json:"email"`
		Message   string    `json:"message"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	MessageRequest struct {
		Firstname string `json:"firstname" form:"firstname" binding:"required"`
		Lastname  string `json:"lastname" form:"lastname" binding:"required"`
		Email     string `json:"email" form:"email" binding:"required,email"`
		Message   string `json:"message" form:"message" binding:"required"`
	}
)

func (mq *MessageQuery) SetDefault() {
	mq.setDefaultPagination()
	mq.SetDefaultSort(SortAsc)
}
