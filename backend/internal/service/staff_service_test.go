package service

import (
	"context"
	"errors"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

type fakeStaffRepo struct {
	staff  map[int64]*domain.Staff
	nextID int64
}

type fakeStaffRoleSync struct {
	lastStaffID int64
	lastRole    string
	err         error
}

func (s *fakeStaffRoleSync) SyncRoleForStaff(ctx context.Context, staffID int64, role string) error {
	_ = ctx
	s.lastStaffID = staffID
	s.lastRole = role
	return s.err
}

type fakeStaffRoleAuditLogger struct {
	entries []StaffRoleAuditEntry
	err     error
}

func (l *fakeStaffRoleAuditLogger) LogRoleChange(ctx context.Context, entry StaffRoleAuditEntry) error {
	_ = ctx
	if l.err != nil {
		return l.err
	}
	l.entries = append(l.entries, entry)
	return nil
}

func newFakeStaffRepo() *fakeStaffRepo {
	return &fakeStaffRepo{staff: make(map[int64]*domain.Staff), nextID: 1}
}

func (r *fakeStaffRepo) Create(ctx context.Context, s *domain.Staff) error {
	s.ID = r.nextID
	r.staff[r.nextID] = s
	r.nextID++
	return nil
}

func (r *fakeStaffRepo) GetByID(ctx context.Context, id int64) (*domain.Staff, error) {
	if s, ok := r.staff[id]; ok {
		return s, nil
	}
	return nil, nil
}

func (r *fakeStaffRepo) Update(ctx context.Context, s *domain.Staff) error {
	r.staff[s.ID] = s
	return nil
}

func (r *fakeStaffRepo) Delete(ctx context.Context, id int64) error {
	delete(r.staff, id)
	return nil
}

func (r *fakeStaffRepo) List(ctx context.Context) ([]domain.Staff, error) {
	var result []domain.Staff
	for _, s := range r.staff {
		result = append(result, *s)
	}
	return result, nil
}

func TestStaffServiceCreate(t *testing.T) {
	repo := newFakeStaffRepo()
	svc := NewStaffService(repo)
	staff, err := svc.Create(context.Background(), "张三", "admin@example.com", "operator")
	if err != nil {
		t.Fatal(err)
	}
	if staff.Role != "operator" {
		t.Fatal("expected operator role")
	}
	assert.Equal(t, "张三", staff.RealName)
}

func TestStaffServiceAssignRole(t *testing.T) {
	repo := newFakeStaffRepo()
	repo.staff[1] = &domain.Staff{ID: 1, RealName: "张三", Role: "operator"}

	svc := NewStaffService(repo)
	err := svc.AssignRole(context.Background(), 1, "finance", 0)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "finance", repo.staff[1].Role)
}

func TestStaffServiceList(t *testing.T) {
	repo := newFakeStaffRepo()
	repo.staff[1] = &domain.Staff{ID: 1, RealName: "张三", Role: "operator"}
	repo.staff[2] = &domain.Staff{ID: 2, RealName: "李四", Role: "admin"}

	svc := NewStaffService(repo)
	list, err := svc.List(context.Background())

	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestStaffServiceAssignRoleSyncsCasbinAndWritesAudit(t *testing.T) {
	repo := newFakeStaffRepo()
	repo.staff[1] = &domain.Staff{ID: 1, RealName: "张三", Role: "operator"}

	syncer := &fakeStaffRoleSync{}
	audit := &fakeStaffRoleAuditLogger{}
	svc := NewStaffServiceWithDeps(repo, syncer, audit)

	err := svc.AssignRole(context.Background(), 1, "finance", 99)
	assert.NoError(t, err)
	assert.Equal(t, "finance", repo.staff[1].Role)
	assert.Equal(t, int64(1), syncer.lastStaffID)
	assert.Equal(t, "finance", syncer.lastRole)
	if assert.Len(t, audit.entries, 1) {
		assert.Equal(t, int64(99), audit.entries[0].OperatorID)
		assert.Equal(t, int64(1), audit.entries[0].TargetStaffID)
		assert.Equal(t, "operator", audit.entries[0].OldRole)
		assert.Equal(t, "finance", audit.entries[0].NewRole)
	}
}

func TestStaffServiceAssignRoleReturnsErrorWhenCasbinSyncFails(t *testing.T) {
	repo := newFakeStaffRepo()
	repo.staff[1] = &domain.Staff{ID: 1, RealName: "张三", Role: "operator"}

	syncer := &fakeStaffRoleSync{err: errors.New("casbin unavailable")}
	audit := &fakeStaffRoleAuditLogger{}
	svc := NewStaffServiceWithDeps(repo, syncer, audit)

	err := svc.AssignRole(context.Background(), 1, "finance", 99)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sync casbin role")
	assert.Len(t, audit.entries, 0)
	// F-9: 验证 Casbin 同步失败后数据库角色已回滚
	assert.Equal(t, "operator", repo.staff[1].Role, "role should be rolled back after casbin sync failure")
}
