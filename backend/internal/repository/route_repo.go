package repository

import (
	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

type RouteRepository struct{ db *gorm.DB }

func (r *RouteRepository) withDB() *gorm.DB { return r.db }

func NewRouteRepository(db *gorm.DB) *RouteRepository { return &RouteRepository{db: db} }

func (r *RouteRepository) Create(v *domain.Route) error { return r.db.Create(v).Error }
func (r *RouteRepository) List() ([]domain.Route, error) {
	var out []domain.Route
	return out, r.db.Order("id desc").Find(&out).Error
}
