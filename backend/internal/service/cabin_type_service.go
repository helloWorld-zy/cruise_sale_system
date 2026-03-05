package service

import (
	"context"

	"github.com/cruisebooking/backend/internal/domain"
)

// CabinTypeService 实现舱房类型相关的业务逻辑。
// 提供舱房类型的 CRUD 操作。
type CabinTypeService struct {
	repo        domain.CabinTypeRepository        // 舱房类型数据仓储
	bindingRepo domain.CabinTypeBindingRepository // 舱型-邮轮绑定仓储（可选）
}

// NewCabinTypeService 创建舱房类型服务实例，通过依赖注入传入仓储。
func NewCabinTypeService(repo domain.CabinTypeRepository, bindingRepo ...domain.CabinTypeBindingRepository) *CabinTypeService {
	var bind domain.CabinTypeBindingRepository
	if len(bindingRepo) > 0 {
		bind = bindingRepo[0]
	}
	return &CabinTypeService{repo: repo, bindingRepo: bind}
}

// Create 创建新的舱房类型。
func (s *CabinTypeService) Create(ctx context.Context, ct *domain.CabinType) error {
	fillCabinTypeDefaults(ct)
	if err := s.repo.Create(ctx, ct); err != nil {
		return err
	}
	if s.bindingRepo != nil && ct.ID > 0 && ct.CruiseID > 0 {
		return s.bindingRepo.ReplaceCruiseBindings(ctx, ct.ID, []int64{ct.CruiseID})
	}
	return nil
}

// Update 保存对已有舱房类型的修改。
func (s *CabinTypeService) Update(ctx context.Context, ct *domain.CabinType) error {
	fillCabinTypeDefaults(ct)
	return s.repo.Update(ctx, ct)
}

// CreateBatchByCruises 按邮轮批量创建舱型，每个邮轮生成一条独立舱型记录。
func (s *CabinTypeService) CreateBatchByCruises(ctx context.Context, base *domain.CabinType, cruiseIDs []int64) ([]domain.CabinType, error) {
	result := make([]domain.CabinType, 0, len(cruiseIDs))
	for _, cruiseID := range cruiseIDs {
		item := *base
		item.ID = 0
		item.CruiseID = cruiseID
		fillCabinTypeDefaults(&item)
		if err := s.repo.Create(ctx, &item); err != nil {
			return nil, err
		}
		if s.bindingRepo != nil && item.ID > 0 {
			if err := s.bindingRepo.ReplaceCruiseBindings(ctx, item.ID, []int64{cruiseID}); err != nil {
				return nil, err
			}
		}
		result = append(result, item)
	}
	return result, nil
}

// GetByID 根据主键查询舱房类型详情。
func (s *CabinTypeService) GetByID(ctx context.Context, id int64) (*domain.CabinType, error) {
	return s.repo.GetByID(ctx, id)
}

// List 返回指定邮轮下的分页舱房类型列表。
func (s *CabinTypeService) List(ctx context.Context, cruiseID int64, page, pageSize int) ([]domain.CabinType, int64, error) {
	return s.repo.ListByCruise(ctx, cruiseID, page, pageSize)
}

// Delete 删除指定的舱房类型。
func (s *CabinTypeService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func fillCabinTypeDefaults(ct *domain.CabinType) {
	if ct == nil {
		return
	}
	if ct.CategoryID <= 0 {
		ct.CategoryID = 1
	}
	if ct.Occupancy <= 0 {
		if ct.Capacity > 0 {
			ct.Occupancy = ct.Capacity
		} else {
			ct.Occupancy = 2
		}
	}
	if ct.Intro == "" && ct.Description != "" {
		ct.Intro = ct.Description
	}
}
