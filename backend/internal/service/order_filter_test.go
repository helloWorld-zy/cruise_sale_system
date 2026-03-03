package service

import (
	"context"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

type fakeBookingRepoForFilter struct {
	orders []domain.Booking
}

func (r *fakeBookingRepoForFilter) ListWithFilter(ctx context.Context, filter OrderFilter) ([]domain.Booking, error) {
	var results []domain.Booking
	for _, b := range r.orders {
		match := true
		if filter.Status != "" && b.Status != filter.Status {
			match = false
		}
		if filter.VoyageID > 0 && b.VoyageID != filter.VoyageID {
			match = false
		}
		if match {
			results = append(results, b)
		}
	}
	return results, nil
}

func TestOrderListFilter(t *testing.T) {
	testOrders := []domain.Booking{
		{ID: 1, Status: "paid", VoyageID: 100},
		{ID: 2, Status: "pending", VoyageID: 101},
		{ID: 3, Status: "paid", VoyageID: 100},
	}
	repo := &fakeBookingRepoForFilter{orders: testOrders}

	results, err := repo.ListWithFilter(context.Background(), OrderFilter{
		Status:   "paid",
		VoyageID: 100,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 filtered orders, got %d", len(results))
	}
	assert.Equal(t, int64(1), results[0].ID)
	assert.Equal(t, int64(3), results[1].ID)
}

func TestOrderListFilter_ByStatus(t *testing.T) {
	testOrders := []domain.Booking{
		{ID: 1, Status: "paid"},
		{ID: 2, Status: "pending"},
		{ID: 3, Status: "paid"},
	}
	repo := &fakeBookingRepoForFilter{orders: testOrders}

	results, err := repo.ListWithFilter(context.Background(), OrderFilter{
		Status: "paid",
	})

	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestOrderListFilter_NoFilter(t *testing.T) {
	testOrders := []domain.Booking{
		{ID: 1, Status: "paid"},
		{ID: 2, Status: "pending"},
	}
	repo := &fakeBookingRepoForFilter{orders: testOrders}

	results, err := repo.ListWithFilter(context.Background(), OrderFilter{})

	assert.NoError(t, err)
	assert.Len(t, results, 2)
}
