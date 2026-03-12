package service

import (
	"context"
	"errors"
	"fmt"

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
	repo        StaffRepository
	roleSyncer  StaffRoleSyncer
	auditLogger StaffRoleAuditLogger
}

type StaffRoleSyncer interface {
	SyncRoleForStaff(ctx context.Context, staffID int64, role string) error
}

type StaffRoleAuditEntry struct {
	OperatorID    int64
	TargetStaffID int64
	OldRole       string
	NewRole       string
}

type StaffRoleAuditLogger interface {
	LogRoleChange(ctx context.Context, entry StaffRoleAuditEntry) error
}

func NewStaffService(repo StaffRepository) *StaffService {
	return &StaffService{repo: repo}
}

func NewStaffServiceWithDeps(repo StaffRepository, roleSyncer StaffRoleSyncer, auditLogger StaffRoleAuditLogger) *StaffService {
	return &StaffService{repo: repo, roleSyncer: roleSyncer, auditLogger: auditLogger}
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

func (s *StaffService) AssignRole(ctx context.Context, id int64, role string, operatorID int64) error {
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
	oldRole := staff.Role
	staff.Role = role
	if err := s.repo.Update(ctx, staff); err != nil {
		return err
	}

	if s.roleSyncer != nil {
		if err := s.roleSyncer.SyncRoleForStaff(ctx, id, role); err != nil {
			// 回滚数据库角色变更
			staff.Role = oldRole
			_ = s.repo.Update(ctx, staff)
			return fmt.Errorf("sync casbin role: %w", err)
		}
	}

	if s.auditLogger != nil {
		if err := s.auditLogger.LogRoleChange(ctx, StaffRoleAuditEntry{
			OperatorID:    operatorID,
			TargetStaffID: id,
			OldRole:       oldRole,
			NewRole:       role,
		}); err != nil {
			// 回滚 Casbin 与数据库
			staff.Role = oldRole
			_ = s.repo.Update(ctx, staff)
			if s.roleSyncer != nil {
				_ = s.roleSyncer.SyncRoleForStaff(ctx, id, oldRole)
			}
			return fmt.Errorf("write role change audit log: %w", err)
		}
	}

	return nil
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
