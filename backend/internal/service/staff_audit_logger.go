package service

import (
	"context"
	"fmt"

	"github.com/cruisebooking/backend/internal/domain"
)

type OperationLogWriter interface {
	Create(ctx context.Context, log *domain.OperationLog) error
}

// StaffOperationLogger writes role change audit logs into operation_logs table.
type StaffOperationLogger struct {
	repo OperationLogWriter
}

func NewStaffOperationLogger(repo OperationLogWriter) *StaffOperationLogger {
	return &StaffOperationLogger{repo: repo}
}

func (l *StaffOperationLogger) LogRoleChange(ctx context.Context, entry StaffRoleAuditEntry) error {
	if l == nil || l.repo == nil {
		return nil
	}

	details := fmt.Sprintf("target_staff_id=%d old_role=%s new_role=%s", entry.TargetStaffID, entry.OldRole, entry.NewRole)
	return l.repo.Create(ctx, &domain.OperationLog{
		StaffID:    entry.OperatorID,
		Operation:  "assign_role",
		Resource:   "staff",
		ResourceID: entry.TargetStaffID,
		Details:    details,
	})
}
