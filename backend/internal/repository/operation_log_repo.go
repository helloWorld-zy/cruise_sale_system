package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

// OperationLogRepository handles persistence for operation_logs.
type OperationLogRepository struct {
	db *gorm.DB
}

func NewOperationLogRepository(db *gorm.DB) *OperationLogRepository {
	return &OperationLogRepository{db: db}
}

func (r *OperationLogRepository) Create(ctx context.Context, log *domain.OperationLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}
