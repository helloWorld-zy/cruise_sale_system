package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/cruisebooking/backend/internal/domain"
)

// StaffRepository 定义员工数据访问接口。
type StaffRepository interface {
	Create(ctx context.Context, s *domain.Staff) error            // 创建员工记录
	GetByID(ctx context.Context, id int64) (*domain.Staff, error) // 根据 ID 查询员工
	Update(ctx context.Context, s *domain.Staff) error            // 更新员工信息
	Delete(ctx context.Context, id int64) error                   // 删除员工记录
	List(ctx context.Context) ([]domain.Staff, error)             // 查询所有员工列表
}

// StaffService 提供员工管理业务逻辑，包括员工创建、角色分配、权限同步等。
type StaffService struct {
	repo        StaffRepository      // 员工数据仓储
	roleSyncer  StaffRoleSyncer      // 角色同步器（用于同步到 Casbin）
	auditLogger StaffRoleAuditLogger // 角色变更审计日志记录器
}

// StaffRoleSyncer 定义角色同步接口，用于将员工角色同步到权限系统（如 Casbin）。
type StaffRoleSyncer interface {
	SyncRoleForStaff(ctx context.Context, staffID int64, role string) error // 同步员工角色
}

// StaffRoleAuditEntry 记录角色变更审计条目。
type StaffRoleAuditEntry struct {
	OperatorID    int64  // 操作人 ID
	TargetStaffID int64  // 被操作员工 ID
	OldRole       string // 变更前角色
	NewRole       string // 变更后角色
}

// StaffRoleAuditLogger 定义角色变更审计日志记录接口。
type StaffRoleAuditLogger interface {
	LogRoleChange(ctx context.Context, entry StaffRoleAuditEntry) error // 记录角色变更日志
}

// NewStaffService 创建员工服务实例（无依赖注入版本）。
func NewStaffService(repo StaffRepository) *StaffService {
	return &StaffService{repo: repo}
}

// NewStaffServiceWithDeps 创建员工服务实例（带完整依赖注入）。
func NewStaffServiceWithDeps(repo StaffRepository, roleSyncer StaffRoleSyncer, auditLogger StaffRoleAuditLogger) *StaffService {
	return &StaffService{repo: repo, roleSyncer: roleSyncer, auditLogger: auditLogger}
}

// Create 创建新员工。验证角色有效性后创建员工记录。
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

// AssignRole 为员工分配新角色。支持权限同步和审计日志记录。
// 如果角色同步或日志记录失败，会自动回滚数据库中的角色变更。
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

// List 查询所有员工列表。
func (s *StaffService) List(ctx context.Context) ([]domain.Staff, error) {
	return s.repo.List(ctx)
}

// GetByID 根据 ID 查询员工详情。
func (s *StaffService) GetByID(ctx context.Context, id int64) (*domain.Staff, error) {
	return s.repo.GetByID(ctx, id)
}

// Update 更新员工信息。
func (s *StaffService) Update(ctx context.Context, staff *domain.Staff) error {
	return s.repo.Update(ctx, staff)
}

// Delete 删除员工（软删除）。
func (s *StaffService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
