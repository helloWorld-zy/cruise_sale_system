package service

import (
	"context"
	"errors"
	"strings"

	"github.com/cruisebooking/backend/internal/domain"
)

// ErrCruiseHasCabins 当删除邮轮时仍有关联的舱房类型时返回此错误。
var ErrCruiseHasCabins = errors.New("cruise has cabins")

// ErrCruiseHasVoyages 当删除邮轮时仍有关联的航次时返回此错误。
var ErrCruiseHasVoyages = errors.New("cruise has voyages")

// CruiseService 实现邮轮相关的业务逻辑。
// 提供邮轮的创建、更新、查询、删除等操作，并包含级联删除保护。
type CruiseService struct {
	cruiseRepo  domain.CruiseRepository    // 邮轮数据仓储
	cabinRepo   domain.CabinTypeRepository // 舱房类型数据仓储（用于级联检查）
	companyRepo domain.CompanyRepository   // 公司数据仓储（用于外键验证）
}

type cruiseBatchStatusWriter interface {
	BatchUpdateStatus(ctx context.Context, ids []int64, status int16) error
}

type cruiseCabinBindingChecker interface {
	HasCabinTypesByCruise(ctx context.Context, cruiseID int64) (bool, error)
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
	if _, err := s.companyRepo.GetByID(ctx, cruise.CompanyID); err != nil {
		return err
	}
	return s.cruiseRepo.Update(ctx, cruise)
}

// GetByID 根据主键查询邮轮详情。
func (s *CruiseService) GetByID(ctx context.Context, id int64) (*domain.Cruise, error) {
	return s.cruiseRepo.GetByID(ctx, id)
}

// List 返回分页的邮轮列表，支持公司、关键词、状态和排序筛选。
func (s *CruiseService) List(ctx context.Context, companyID int64, keyword string, status *int16, sortBy string, page, pageSize int) ([]domain.Cruise, int64, error) {
	return s.cruiseRepo.List(ctx, companyID, keyword, status, sortBy, page, pageSize)
}

// ListPublic 返回前台可见邮轮列表，仅包含启用中的邮轮。
func (s *CruiseService) ListPublic(ctx context.Context, companyID int64, keyword string, sortBy string, page, pageSize int) ([]domain.Cruise, int64, error) {
	return s.cruiseRepo.ListPublic(ctx, companyID, keyword, sortBy, page, pageSize)
}

// ListWithFilters 为调用方提供显式的筛选列表入口，内部转发到仓储层。
func (s *CruiseService) ListWithFilters(ctx context.Context, companyID int64, keyword string, status *int16, sortBy string, page, pageSize int) ([]domain.Cruise, int64, error) {
	return s.List(ctx, companyID, keyword, status, sortBy, page, pageSize)
}

// BatchUpdateStatus 批量更新邮轮状态。
func (s *CruiseService) BatchUpdateStatus(ctx context.Context, ids []int64, status int16) error {
	if writer, ok := s.cruiseRepo.(cruiseBatchStatusWriter); ok {
		return writer.BatchUpdateStatus(ctx, ids, status)
	}
	for _, id := range ids {
		item, err := s.cruiseRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		item.Status = status
		if err := s.cruiseRepo.Update(ctx, item); err != nil {
			return err
		}
	}
	return nil
}

// Delete 删除邮轮前检查是否仍有关联的舱房类型，防止级联数据不一致。
func (s *CruiseService) Delete(ctx context.Context, id int64) error {
	if checker, ok := s.cabinRepo.(cruiseCabinBindingChecker); ok {
		has, err := checker.HasCabinTypesByCruise(ctx, id)
		if err != nil {
			return err
		}
		if has {
			return ErrCruiseHasCabins
		}
	} else {
		cabins, total, err := s.cabinRepo.ListByCruise(ctx, id, 1, 1)
		if err != nil {
			return err
		}
		if total > 0 || len(cabins) > 0 {
			return ErrCruiseHasCabins
		}
	}
	if err := s.cruiseRepo.Delete(ctx, id); err != nil {
		// PostgreSQL FK: voyages.cruise_id -> cruises.id
		msg := strings.ToLower(err.Error())
		if strings.Contains(msg, "voyages_cruise_id_fkey") ||
			(strings.Contains(msg, "sqlstate 23503") && strings.Contains(msg, "voyages")) {
			return ErrCruiseHasVoyages
		}
		return err
	}
	return nil
}
