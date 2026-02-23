package service

import "github.com/cruisebooking/backend/internal/domain"

type InventoryRepo interface {
	GetBySKU(id int64) (domain.CabinInventory, error)
	Update(inv *domain.CabinInventory) error
}

type InventoryService struct{ repo InventoryRepo }

func NewInventoryService(repo InventoryRepo) *InventoryService { return &InventoryService{repo: repo} }

func (s *InventoryService) Adjust(skuID int64, delta int) error {
	inv, err := s.repo.GetBySKU(skuID)
	if err != nil {
		return err
	}
	if inv.Total+delta < 0 {
		return nil
	}
	inv.Total += delta
	return s.repo.Update(&inv)
}

func (s *InventoryService) Available(skuID int64) (int, error) {
	inv, err := s.repo.GetBySKU(skuID)
	if err != nil {
		return 0, err
	}
	return inv.Total - inv.Locked - inv.Sold, nil
}
