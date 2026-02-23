package service

import (
	"context"
	"errors"

	"github.com/cruisebooking/backend/internal/domain"
)

// ErrCompanyHasCruises 当删除公司时仍有关联的邮轮时返回此错误。
var ErrCompanyHasCruises = errors.New("company has cruises")

// CompanyService 实现邮轮公司相关的业务逻辑。
// 提供公司的 CRUD 操作，并在删除时进行级联保护检查。
type CompanyService struct {
	repo       domain.CompanyRepository // 公司数据仓储
	cruiseRepo domain.CruiseRepository  // 邮轮数据仓储（用于级联检查）
}

// NewCompanyService 创建公司服务实例。
func NewCompanyService(repo domain.CompanyRepository, cruiseRepo domain.CruiseRepository) *CompanyService {
	return &CompanyService{repo: repo, cruiseRepo: cruiseRepo}
}

// Create 创建新的邮轮公司。
func (s *CompanyService) Create(ctx context.Context, company *domain.CruiseCompany) error {
	return s.repo.Create(ctx, company)
}

// Update 保存对已有公司信息的修改。
func (s *CompanyService) Update(ctx context.Context, company *domain.CruiseCompany) error {
	return s.repo.Update(ctx, company)
}

// GetByID 根据主键查询公司详情。
func (s *CompanyService) GetByID(ctx context.Context, id int64) (*domain.CruiseCompany, error) {
	return s.repo.GetByID(ctx, id)
}

// List 返回分页的公司列表，支持关键词模糊搜索。
func (s *CompanyService) List(ctx context.Context, keyword string, page, pageSize int) ([]domain.CruiseCompany, int64, error) {
	return s.repo.List(ctx, keyword, page, pageSize)
}

// Delete 删除公司前检查是否仍有关联的邮轮，防止级联数据不一致。
// HI-02 FIX：此前缺少级联检查逻辑，已修复。
func (s *CompanyService) Delete(ctx context.Context, id int64) error {
	_, total, err := s.cruiseRepo.List(ctx, id, 1, 1)
	if err != nil {
		return err
	}
	if total > 0 {
		return ErrCompanyHasCruises
	}
	return s.repo.Delete(ctx, id)
}
