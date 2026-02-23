package repository

import (
	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

type CabinRepository struct{ db *gorm.DB }

func NewCabinRepository(db *gorm.DB) *CabinRepository { return &CabinRepository{db: db} }

func (r *CabinRepository) CreateSKU(v *domain.CabinSKU) error { return r.db.Create(v).Error }
func (r *CabinRepository) ListSKUByVoyage(voyageID int64) ([]domain.CabinSKU, error) {
	var out []domain.CabinSKU
	return out, r.db.Where("voyage_id = ?", voyageID).Order("id desc").Find(&out).Error
}
