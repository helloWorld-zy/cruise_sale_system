package service

import (
	"context"
	"errors"

	"github.com/cruisebooking/backend/internal/domain"
)

// ErrCruiseHasCabins 当删除邮轮时仍有关联的舱房类型时返回此错误。
var ErrCruiseHasCabins = errors.New("cruise has cabins")

// CruiseService 实现邮轮相关的业务逻辑。
// 提供邮轮的创建、更新、查询、删除等操作，并包含级联删除保护。
type CruiseService struct {
	cruiseRepo  domain.CruiseRepository    // 邮轮数据仓储
	cabinRepo   domain.CabinTypeRepository // 舱房类型数据仓储（用于级联检查）
	companyRepo domain.CompanyRepository   // 公司数据仓储（用于外键验证）
}

// NewCruiseService 创建邮轮服务实例，通过依赖注入传入所需的仓储。
func NewCruiseService(cruiseRepo domain.CruiseRepository, cabinRepo domain.CabinTypeRepository, companyRepo domain.CompanyRepository) *CruiseService {
	return &CruiseService{cruiseRepo: cruiseRepo, cabinRepo: cabinRepo, companyRepo: companyRepo}
}

// Create 创建邮轮前先验证所属公司是否存在，避免外键引用无效。
func (s *CruiseService) Create(ctx context.Context, cruise *domain.Cruise) error {
	if _, err := s.companyRepo.GetByID(ctx, cruise.CompanyID); err != nil {
		return err
	}
	return s.cruiseRepo.Create(ctx, cruise)
}

// Update 保存对已有邮轮的修改。
func (s *CruiseService) Update(ctx context.Context, cruise *domain.Cruise) error {
	return s.cruiseRepo.Update(ctx, cruise)
}

// GetByID 根据主键查询邮轮详情。
func (s *CruiseService) GetByID(ctx context.Context, id int64) (*domain.Cruise, error) {
	return s.cruiseRepo.GetByID(ctx, id)
}

// List 返回分页的邮轮列表，可选按公司 ID 过滤。
func (s *CruiseService) List(ctx context.Context, companyID int64, page, pageSize int) ([]domain.Cruise, int64, error) {
	return s.cruiseRepo.List(ctx, companyID, page, pageSize)
}

// Delete 删除邮轮前检查是否仍有关联的舱房类型，防止级联数据不一致。
func (s *CruiseService) Delete(ctx context.Context, id int64) error {
	cabins, total, err := s.cabinRepo.ListByCruise(ctx, id, 1, 1)
	if err != nil {
		return err
	}
	if total > 0 || len(cabins) > 0 {
		return ErrCruiseHasCabins
	}
	return s.cruiseRepo.Delete(ctx, id)
}
