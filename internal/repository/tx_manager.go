package repository

import (
	"context"

	"gorm.io/gorm"
)

type TxManager struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) *TxManager {
	return &TxManager{db: db}
}

func (tm *TxManager) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return tm.db.WithContext(ctx).Transaction(fn)
}
