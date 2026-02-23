package repository

import (
	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

type VoyageRepository struct{ db *gorm.DB }

func NewVoyageRepository(db *gorm.DB) *VoyageRepository { return &VoyageRepository{db: db} }

func (r *VoyageRepository) Create(v *domain.Voyage) error { return r.db.Create(v).Error }
func (r *VoyageRepository) ListByRoute(routeID int64) ([]domain.Voyage, error) {
	var out []domain.Voyage
	return out, r.db.Where("route_id = ?", routeID).Order("depart_date asc").Find(&out).Error
}
