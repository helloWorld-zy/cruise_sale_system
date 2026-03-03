package service

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/casbin/casbin/v2"
)

// CasbinStaffRoleSync keeps Casbin grouping policy aligned with staff role changes.
type CasbinStaffRoleSync struct {
	enforcer *casbin.Enforcer
	mu       sync.Mutex
}

func NewCasbinStaffRoleSync(enforcer *casbin.Enforcer) *CasbinStaffRoleSync {
	return &CasbinStaffRoleSync{enforcer: enforcer}
}

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
