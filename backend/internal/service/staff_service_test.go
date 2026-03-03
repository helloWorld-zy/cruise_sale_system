package service

import (
	"context"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

type fakeStaffRepo struct {
	staff  map[int64]*domain.Staff
	nextID int64
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
	err := svc.AssignRole(context.Background(), 1, "finance")
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
