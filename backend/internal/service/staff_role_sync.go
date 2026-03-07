package service

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/casbin/casbin/v2"
)

// CasbinStaffRoleSync 将 Casbin 分组策略与员工角色变更保持同步。
// 实现 StaffRoleSyncer 接口，负责将员工角色同步到 Casbin 权限系统。
type CasbinStaffRoleSync struct {
	enforcer *casbin.Enforcer // Casbin 权限执行器
	mu       sync.Mutex       // 互斥锁，保证并发安全
}

// NewCasbinStaffRoleSync 创建 Casbin 角色同步器实例。
func NewCasbinStaffRoleSync(enforcer *casbin.Enforcer) *CasbinStaffRoleSync {
	return &CasbinStaffRoleSync{enforcer: enforcer}
}

// SyncRoleForStaff 同步员工角色到 Casbin 权限系统。
// 流程：移除旧分组策略 → 添加新分组策略 → 保存策略到存储。
func (s *CasbinStaffRoleSync) SyncRoleForStaff(ctx context.Context, staffID int64, role string) error {
	_ = ctx
	if s == nil || s.enforcer == nil {
		return nil
	}

	subject := strconv.FormatInt(staffID, 10)

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := s.enforcer.RemoveFilteredGroupingPolicy(0, subject); err != nil {
		return fmt.Errorf("remove old grouping policy: %w", err)
	}
	if _, err := s.enforcer.AddGroupingPolicy(subject, role); err != nil {
		return fmt.Errorf("add grouping policy: %w", err)
	}
	if err := s.enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("save policy: %w", err)
	}
	return nil
}
