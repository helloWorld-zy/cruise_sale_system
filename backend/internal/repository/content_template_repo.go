package repository

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

type ContentTemplateRepository struct{ db *gorm.DB }

func NewContentTemplateRepository(db *gorm.DB) *ContentTemplateRepository {
	return &ContentTemplateRepository{db: db}
}

func (r *ContentTemplateRepository) List(ctx context.Context, kind domain.ContentTemplateKind) ([]domain.ContentTemplate, error) {
	var list []domain.ContentTemplate
	query := r.db.WithContext(ctx).Order("id desc")
	if kind != "" {
		query = query.Where("kind = ?", kind)
	}
	err := query.Find(&list).Error
	return list, err
}

func (r *ContentTemplateRepository) GetByID(ctx context.Context, id int64) (*domain.ContentTemplate, error) {
	var tpl domain.ContentTemplate
	err := r.db.WithContext(ctx).First(&tpl, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &tpl, nil
}

func (r *ContentTemplateRepository) Create(ctx context.Context, tpl *domain.ContentTemplate) error {
	return r.db.WithContext(ctx).Create(tpl).Error
}

func (r *ContentTemplateRepository) Update(ctx context.Context, tpl *domain.ContentTemplate) error {
	return r.db.WithContext(ctx).Save(tpl).Error
}

func (r *ContentTemplateRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.ContentTemplate{}, id).Error
}

var _ domain.ContentTemplateRepository = (*ContentTemplateRepository)(nil)
