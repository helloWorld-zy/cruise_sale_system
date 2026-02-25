package service

import (
	"context"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
)

// PriceRepo 定义价格查询的端口接口。
type PriceRepo interface {
	ListBySKU(ctx context.Context, skuID int64) ([]domain.CabinPrice, error) // 查询某 SKU 的所有价格记录
}

// PricingService 提供舱房定价相关的业务逻辑。
type PricingService struct{ repo PriceRepo }

// NewPricingService 创建定价服务实例。
func NewPricingService(repo PriceRepo) *PricingService { return &PricingService{repo: repo} }

// sameDay 判断两个时间是否在同一天（按 UTC 比较）。
// 两个输入都会先转换为 UTC，以避免数据库存储的 UTC 时间
// 与客户端传入的本地时间（如 +08:00）之间的时区差异。
// HIGH-01 修复项。
func sameDay(a, b time.Time) bool {
	a, b = a.UTC(), b.UTC()
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// FindPrice 查找指定 SKU 在某个日期、某个入住人数下的价格（单位：分）。
// 返回值：价格金额、是否找到、错误信息。
// 通过返回 error 使调用方能够区分"无价格"和"数据库故障"（HIGH-02 修复项）。
func (s *PricingService) FindPrice(ctx context.Context, skuID int64, date time.Time, occupancy int) (int64, bool, error) {
	list, err := s.repo.ListBySKU(ctx, skuID)
	if err != nil {
		return 0, false, err
	}
	for _, v := range list {
		if sameDay(v.Date, date) && v.Occupancy == occupancy {
			return v.PriceCents, true, nil
		}
	}
	return 0, false, nil
}
