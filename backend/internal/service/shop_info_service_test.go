package service

import (
	"context"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

type fakeShopInfoRepo struct {
	info *domain.ShopInfo
}

func (r *fakeShopInfoRepo) Get(ctx context.Context) (*domain.ShopInfo, error) {
	return r.info, nil
}

func (r *fakeShopInfoRepo) Save(ctx context.Context, info *domain.ShopInfo) error {
	r.info = info
	return nil
}

func TestShopInfoService(t *testing.T) {
	repo := &fakeShopInfoRepo{info: &domain.ShopInfo{Name: "Test Shop"}}
	svc := NewShopInfoService(repo)
	info, err := svc.Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if info.Name == "" {
		t.Fatal("expected shop name")
	}
	assert.Equal(t, "Test Shop", info.Name)
}

func TestShopInfoService_Update(t *testing.T) {
	repo := &fakeShopInfoRepo{info: &domain.ShopInfo{Name: "Old Name"}}
	svc := NewShopInfoService(repo)

	err := svc.Update(context.Background(), &domain.ShopInfo{Name: "New Name"})
	assert.NoError(t, err)
	assert.Equal(t, "New Name", repo.info.Name)
}
