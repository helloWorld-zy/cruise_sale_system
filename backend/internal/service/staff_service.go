package service

import (
	"context"
	"errors"

	"github.com/cruisebooking/backend/internal/domain"
)

type StaffRepository interface {
	Create(ctx context.Context, s *domain.Staff) error
	GetByID(ctx context.Context, id int64) (*domain.Staff, error)
	Update(ctx context.Context, s *domain.Staff) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]domain.Staff, error)
}

type StaffService struct {
	repo StaffRepository
}

func NewStaffService(repo StaffRepository) *StaffService {
	return &StaffService{repo: repo}
}

func (s *StaffService) Create(ctx context.Context, name, email, role string) (*domain.Staff, error) {
	if !domain.IsValidStaffRole(role) {
		return nil, errors.New("invalid role")
	}
	staff := &domain.Staff{
		RealName: name,
		Email:    email,
		Role:     role,
		Status:   1,
	}
	if err := s.repo.Create(ctx, staff); err != nil {
		return nil, err
	}
	return staff, nil
}

func (s *StaffService) AssignRole(ctx context.Context, id int64, role string) error {
	if !domain.IsValidStaffRole(role) {
		return errors.New("invalid role")
	}
	staff, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if staff == nil {
		return errors.New("staff not found")
	}
	staff.Role = role
	return s.repo.Update(ctx, staff)
}

func (s *StaffService) List(ctx context.Context) ([]domain.Staff, error) {
	return s.repo.List(ctx)
}

func (s *StaffService) GetByID(ctx context.Context, id int64) (*domain.Staff, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *StaffService) Update(ctx context.Context, staff *domain.Staff) error {
	return s.repo.Update(ctx, staff)
}

func (s *StaffService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
