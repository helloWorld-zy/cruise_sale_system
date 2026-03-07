package service

import (
	"context"
	"fmt"

	"github.com/cruisebooking/backend/internal/domain"
)

// OperationLogWriter 定义操作日志写入接口。
type OperationLogWriter interface {
	Create(ctx context.Context, log *domain.OperationLog) error // 写入操作日志
}

// StaffOperationLogger 将角色变更审计日志写入 operation_logs 表。
// 实现 StaffRoleAuditLogger 接口。
type StaffOperationLogger struct {
	repo OperationLogWriter // 操作日志仓储
}

// NewStaffOperationLogger 创建员工操作日志记录器。
func NewStaffOperationLogger(repo OperationLogWriter) *StaffOperationLogger {
	return &StaffOperationLogger{repo: repo}
}

// LogRoleChange 记录角色变更操作日志。
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
