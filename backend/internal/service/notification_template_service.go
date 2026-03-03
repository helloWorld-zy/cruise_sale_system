package service

import (
	"context"
	"errors"

	"github.com/cruisebooking/backend/internal/domain"
)

type NotificationTemplateRepository interface {
	List(ctx context.Context) ([]domain.NotificationTemplate, error)
	GetByID(ctx context.Context, id int64) (*domain.NotificationTemplate, error)
	Create(ctx context.Context, tpl *domain.NotificationTemplate) error
	Update(ctx context.Context, tpl *domain.NotificationTemplate) error
	Delete(ctx context.Context, id int64) error
}

type NotificationTemplateService struct {
	repo NotificationTemplateRepository
}

func NewNotificationTemplateService(repo NotificationTemplateRepository) *NotificationTemplateService {
	return &NotificationTemplateService{repo: repo}
}

func (s *NotificationTemplateService) List(ctx context.Context) ([]domain.NotificationTemplate, error) {
	return s.repo.List(ctx)
}

func (s *NotificationTemplateService) Create(ctx context.Context, tpl *domain.NotificationTemplate) error {
	if tpl == nil {
		return errors.New("template is required")
	}
	if tpl.EventType == "" || tpl.Template == "" {
		return errors.New("event_type and template are required")
	}
	return s.repo.Create(ctx, tpl)
}

func (s *NotificationTemplateService) Update(ctx context.Context, tpl *domain.NotificationTemplate) error {
	if tpl == nil || tpl.ID <= 0 {
		return errors.New("id is required")
	}
	old, err := s.repo.GetByID(ctx, tpl.ID)
	if err != nil {
		return err
	}
	if old == nil {
		return errors.New("template not found")
	}
	if tpl.EventType == "" {
		tpl.EventType = old.EventType
	}
	if tpl.Channel == "" {
		tpl.Channel = old.Channel
	}
	if tpl.Template == "" {
		tpl.Template = old.Template
	}
	return s.repo.Update(ctx, tpl)
}

func (s *NotificationTemplateService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid id")
	}
	return s.repo.Delete(ctx, id)
}
