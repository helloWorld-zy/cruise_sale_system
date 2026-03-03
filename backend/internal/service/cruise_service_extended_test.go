package service

import (
	"context"
	"errors"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
)

// mockCruiseRepoFull 用于验证批量状态更新流程。
type mockCruiseRepoFull struct {
	items map[int64]*domain.Cruise
}

func (m *mockCruiseRepoFull) Create(ctx context.Context, cruise *domain.Cruise) error { return nil }
func (m *mockCruiseRepoFull) Update(ctx context.Context, cruise *domain.Cruise) error {
	if m.items == nil {
		m.items = map[int64]*domain.Cruise{}
	}
	m.items[cruise.ID] = cruise
	return nil
}
func (m *mockCruiseRepoFull) GetByID(ctx context.Context, id int64) (*domain.Cruise, error) {
	if m.items == nil {
		m.items = map[int64]*domain.Cruise{}
	}
	if item, ok := m.items[id]; ok {
		return item, nil
	}
	item := &domain.Cruise{ID: id, Status: 1}
	m.items[id] = item
	return item, nil
}
func (m *mockCruiseRepoFull) List(ctx context.Context, companyID int64, keyword string, status *int16, sortBy string, page, pageSize int) ([]domain.Cruise, int64, error) {
	return nil, 0, nil
}
func (m *mockCruiseRepoFull) Delete(ctx context.Context, id int64) error { return nil }

// mockImageRepo 用于验证画廊覆盖写入。
type mockImageRepo struct {
	items []domain.Image
}

func (m *mockImageRepo) Create(ctx context.Context, img *domain.Image) error {
	m.items = append(m.items, *img)
	return nil
}
func (m *mockImageRepo) ListByEntity(ctx context.Context, entityType string, entityID int64) ([]domain.Image, error) {
	var out []domain.Image
	for _, it := range m.items {
		if it.EntityType == entityType && it.EntityID == entityID {
			out = append(out, it)
		}
	}
	return out, nil
}
func (m *mockImageRepo) DeleteByEntity(ctx context.Context, entityType string, entityID int64) error {
	filtered := make([]domain.Image, 0, len(m.items))
	for _, it := range m.items {
		if it.EntityType == entityType && it.EntityID == entityID {
			continue
		}
		filtered = append(filtered, it)
	}
	m.items = filtered
	return nil
}
func (m *mockImageRepo) UpdateSortOrder(ctx context.Context, id int64, sortOrder int) error {
	for i := range m.items {
		if m.items[i].ID == id {
			m.items[i].SortOrder = sortOrder
			return nil
		}
	}
	return errors.New("not found")
}
func (m *mockImageRepo) ReplaceImages(ctx context.Context, entityType string, entityID int64, images []*domain.Image) error {
	_ = m.DeleteByEntity(ctx, entityType, entityID)
	for _, img := range images {
		m.items = append(m.items, *img)
	}
	return nil
}

func TestCruiseService_BatchUpdateStatus(t *testing.T) {
	repo := &mockCruiseRepoFull{}
	svc := NewCruiseService(repo, &mockCabinRepo{}, &mockCompanyRepo{})
	err := svc.BatchUpdateStatus(context.Background(), []int64{1, 2}, 0)
	if err != nil {
		t.Fatal("expected success")
	}
	if repo.items[1].Status != 0 || repo.items[2].Status != 0 {
		t.Fatal("expected all status updated")
	}
}

func TestImageService_ManageGallery(t *testing.T) {
	repo := &mockImageRepo{}
	svc := NewImageService(repo)
	err := svc.SetImages(context.Background(), "cruise", 1, []ImageInput{{URL: "http://a.jpg", IsPrimary: true}})
	if err != nil {
		t.Fatal(err)
	}
	items, err := svc.ListImages(context.Background(), "cruise", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || !items[0].IsPrimary {
		t.Fatal("expected primary image")
	}
}
