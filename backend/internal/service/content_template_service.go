package service

import (
	"context"
	"errors"

	"github.com/cruisebooking/backend/internal/domain"
)

type ContentTemplateService struct {
	repo domain.ContentTemplateRepository
}

func NewContentTemplateService(repo domain.ContentTemplateRepository) *ContentTemplateService {
	return &ContentTemplateService{repo: repo}
}

func (s *ContentTemplateService) List(ctx context.Context, kind domain.ContentTemplateKind) ([]domain.ContentTemplate, error) {
	return s.repo.List(ctx, kind)
}

func (s *ContentTemplateService) GetByID(ctx context.Context, id int64) (*domain.ContentTemplate, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *ContentTemplateService) Create(ctx context.Context, tpl *domain.ContentTemplate) error {
	if tpl == nil {
		return errors.New("template is required")
	}
	if tpl.Name == "" || tpl.Kind == "" || tpl.ContentJSON == "" {
		return errors.New("name, kind and content are required")
	}
	if tpl.Status == 0 {
		tpl.Status = 1
	}
	return s.repo.Create(ctx, tpl)
}

func (s *ContentTemplateService) Update(ctx context.Context, tpl *domain.ContentTemplate) error {
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
	if tpl.Name == "" {
		tpl.Name = old.Name
	}
	if tpl.Kind == "" {
		tpl.Kind = old.Kind
	}
	if tpl.ContentJSON == "" {
		tpl.ContentJSON = old.ContentJSON
	}
	if tpl.Status == 0 {
		tpl.Status = old.Status
	}
	return s.repo.Update(ctx, tpl)
}

func (s *ContentTemplateService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid id")
	}
	return s.repo.Delete(ctx, id)
}
