package domain

import (
	"context"
	"errors"
)

// ErrInsufficientInventory 当库存调整会导致总量为负时返回此错误。
// 定义在领域层以便仓储层和服务层都能引用，避免循环依赖。
var ErrInsufficientInventory = errors.New("insufficient inventory")

// CompanyRepository 定义邮轮公司的数据持久化接口。
// 遵循 DDD 仓储模式，由基础设施层（repository 包）提供具体实现。
type CompanyRepository interface {
	Create(ctx context.Context, company *CruiseCompany) error                                     // 创建邮轮公司
	Update(ctx context.Context, company *CruiseCompany) error                                     // 更新邮轮公司信息
	GetByID(ctx context.Context, id int64) (*CruiseCompany, error)                                // 根据 ID 查询邮轮公司
	List(ctx context.Context, keyword string, page, pageSize int) ([]CruiseCompany, int64, error) // 分页查询公司列表，支持关键词搜索
	Delete(ctx context.Context, id int64) error                                                   // 删除邮轮公司（软删除）
}

// CruiseRepository 定义邮轮的数据持久化接口。
type CruiseRepository interface {
	Create(ctx context.Context, cruise *Cruise) error                                       // 创建邮轮
	Update(ctx context.Context, cruise *Cruise) error                                       // 更新邮轮信息
	GetByID(ctx context.Context, id int64) (*Cruise, error)                                 // 根据 ID 查询邮轮
	List(ctx context.Context, companyID int64, page, pageSize int) ([]Cruise, int64, error) // 分页查询邮轮列表，可按公司筛选
	Delete(ctx context.Context, id int64) error                                             // 删除邮轮（软删除）
}

// CabinTypeRepository 定义舱房类型的数据持久化接口。
type CabinTypeRepository interface {
	Create(ctx context.Context, cabinType *CabinType) error                                           // 创建舱房类型
	Update(ctx context.Context, cabinType *CabinType) error                                           // 更新舱房类型信息
	GetByID(ctx context.Context, id int64) (*CabinType, error)                                        // 根据 ID 查询舱房类型
	ListByCruise(ctx context.Context, cruiseID int64, page, pageSize int) ([]CabinType, int64, error) // 分页查询某邮轮下的舱房类型
	Delete(ctx context.Context, id int64) error                                                       // 删除舱房类型（软删除）
}

// FacilityCategoryRepository 定义设施分类的数据持久化接口。
type FacilityCategoryRepository interface {
	Create(ctx context.Context, category *FacilityCategory) error // 创建设施分类
	List(ctx context.Context) ([]FacilityCategory, error)         // 查询所有设施分类
	Delete(ctx context.Context, id int64) error                   // 删除设施分类
}

// FacilityRepository 定义设施的数据持久化接口。
type FacilityRepository interface {
	Create(ctx context.Context, facility *Facility) error                 // 创建设施
	ListByCruise(ctx context.Context, cruiseID int64) ([]Facility, error) // 查询某邮轮下的所有设施
	Delete(ctx context.Context, id int64) error                           // 删除设施（软删除）
}

// StaffRepository 定义员工账号的数据持久化接口。
type StaffRepository interface {
	Create(ctx context.Context, staff *Staff) error                     // 创建员工账号
	GetByUsername(ctx context.Context, username string) (*Staff, error) // 根据用户名查询员工
	GetByID(ctx context.Context, id int64) (*Staff, error)              // 根据 ID 查询员工
	Update(ctx context.Context, staff *Staff) error                     // 更新员工信息
	Delete(ctx context.Context, id int64) error                         // 删除员工账号
}

// Sprint 2 仓储端口 —— 按照 DDD 规范定义在领域层。

// RouteRepository 定义航线的数据持久化接口。
type RouteRepository interface {
	Create(ctx context.Context, r *Route) error            // 创建航线
	Update(ctx context.Context, r *Route) error            // 更新航线信息
	GetByID(ctx context.Context, id int64) (*Route, error) // 根据 ID 查询航线
	List(ctx context.Context) ([]Route, error)             // 查询所有航线
	Delete(ctx context.Context, id int64) error            // 删除航线
}

// VoyageRepository 定义航次的数据持久化接口。
type VoyageRepository interface {
	Create(ctx context.Context, v *Voyage) error                      // 创建航次
	Update(ctx context.Context, v *Voyage) error                      // 更新航次信息
	GetByID(ctx context.Context, id int64) (*Voyage, error)           // 根据 ID 查询航次
	ListByRoute(ctx context.Context, routeID int64) ([]Voyage, error) // 查询某航线下的所有航次
	Delete(ctx context.Context, id int64) error                       // 删除航次
}

// CabinSKURepository 聚合了舱房产品的所有存储操作，
// 包括 SKU 的 CRUD、原子化库存调整和价格日历管理。
type CabinSKURepository interface {
	CreateSKU(ctx context.Context, s *CabinSKU) error                        // 创建舱房 SKU
	UpdateSKU(ctx context.Context, s *CabinSKU) error                        // 更新舱房 SKU
	GetSKUByID(ctx context.Context, id int64) (*CabinSKU, error)             // 根据 ID 查询舱房 SKU
	ListSKUByVoyage(ctx context.Context, voyageID int64) ([]CabinSKU, error) // 查询某航次下的所有舱房 SKU
	DeleteSKU(ctx context.Context, id int64) error                           // 删除舱房 SKU
	// AdjustInventoryAtomic 使用单条 SQL 语句原子化更新库存总量，
	// 防止并发请求导致超卖。当 total+delta < 0 时返回 ErrInsufficientInventory。
	AdjustInventoryAtomic(ctx context.Context, skuID int64, delta int) error
	GetInventoryBySKU(ctx context.Context, skuID int64) (CabinInventory, error) // 查询舱房库存
	AppendInventoryLog(ctx context.Context, log *InventoryLog) error            // 追加库存变动日志
	ListPricesBySKU(ctx context.Context, skuID int64) ([]CabinPrice, error)     // 查询某 SKU 的价格列表
	UpsertPrice(ctx context.Context, p *CabinPrice) error                       // 新增或更新价格记录
}

// --- Sprint 4: 支付 / 退款 / 通知 / 分析仓储 ---

// PaymentRepository 定义支付持久化操作。
type PaymentRepository interface {
	Create(ctx context.Context, p *Payment) error
	FindByTradeNo(ctx context.Context, tradeNo string) (*Payment, error)
	FindByID(ctx context.Context, id int64) (*Payment, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
}

// RefundRepository 定义退款持久化操作。
type RefundRepository interface {
	Create(ctx context.Context, r *Refund) error
	// SumByPaymentID 返回支付的已批准/待定退款总额，
	// 用于强制执行累计退款限制。
	SumByPaymentID(ctx context.Context, paymentID int64) (int64, error)
}

// NotificationRepository 定义发件箱模式的发件箱持久化操作。
type NotificationRepository interface {
	CreateOutbox(ctx context.Context, n *Notification) error
	ListPending(ctx context.Context, limit int) ([]Notification, error)
	MarkSent(ctx context.Context, id int64) error
	MarkFailed(ctx context.Context, id int64) error
}

// AnalyticsRepository 定义只读分析查询操作。
type AnalyticsRepository interface {
	TodaySales(ctx context.Context) (int64, error)
	WeeklyTrend(ctx context.Context) ([]int64, error)
	TodayOrderCount(ctx context.Context) (int64, error)
}

// BookingStatusRepository 提供订单状态更新功能，
// 供支付回调服务确认订单支付。
type BookingStatusRepository interface {
	UpdateStatus(ctx context.Context, id int64, status string) error
}
