package service

import (
	"context"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
)

type fakePassengerRepo struct {
	passengers []domain.Passenger
}

func (f *fakePassengerRepo) ListByUser(ctx context.Context, userID int64) ([]domain.Passenger, error) {
	return f.passengers, nil
}

func (f *fakePassengerRepo) UpdateFavorite(ctx context.Context, id int64, isFavorite bool) error {
	for i := range f.passengers {
		if f.passengers[i].ID == id {
			f.passengers[i].IsFavorite = isFavorite
			break
		}
	}
	return nil
}

func TestPassengerServiceListFavorites(t *testing.T) {
	repo := &fakePassengerRepo{passengers: []domain.Passenger{
		{ID: 1, UserID: 1, Name: "张三", IsFavorite: true},
		{ID: 2, UserID: 1, Name: "李四", IsFavorite: false},
	}}
	svc := NewPassengerService(repo)
	favorites, err := svc.ListFavorites(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(favorites) != 1 {
		t.Fatalf("expected 1 favorite, got %d", len(favorites))
	}
}

func TestPassengerServiceToggleFavorite(t *testing.T) {
	repo := &fakePassengerRepo{passengers: []domain.Passenger{
		{ID: 1, UserID: 1, Name: "张三", IsFavorite: false},
	}}
	svc := NewPassengerService(repo)
	err := svc.ToggleFavorite(context.Background(), 1, true)
	if err != nil {
		t.Fatal(err)
	}
	if !repo.passengers[0].IsFavorite {
		t.Fatal("expected favorite to be true")
	}
}
