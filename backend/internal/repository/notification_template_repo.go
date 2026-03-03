package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

// NotificationTemplateRepository 提供通知模板 CRUD。
type NotificationTemplateRepository struct {
	db *gorm.DB
}

func NewNotificationTemplateRepository(db *gorm.DB) *NotificationTemplateRepository {
	return &NotificationTemplateRepository{db: db}
}

func (r *NotificationTemplateRepository) List(ctx context.Context) ([]domain.NotificationTemplate, error) {
	var list []domain.NotificationTemplate
	err := r.db.WithContext(ctx).Order("id DESC").Find(&list).Error
	return list, err
}

func (r *NotificationTemplateRepository) GetByID(ctx context.Context, id int64) (*domain.NotificationTemplate, error) {
	var tpl domain.NotificationTemplate
	err := r.db.WithContext(ctx).First(&tpl, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &tpl, nil
}

func (r *NotificationTemplateRepository) Create(ctx context.Context, tpl *domain.NotificationTemplate) error {
	return r.db.WithContext(ctx).Create(tpl).Error
}

func (r *NotificationTemplateRepository) Update(ctx context.Context, tpl *domain.NotificationTemplate) error {
	return r.db.WithContext(ctx).Save(tpl).Error
}

func (r *NotificationTemplateRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.NotificationTemplate{}, id).Error
}
