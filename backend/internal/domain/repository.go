package domain

import (
	"context"
	"errors"
	"time"
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
	ListPublic(ctx context.Context, page, pageSize int) ([]CruiseCompany, int64, error)           // 查询前台可见公司列表（仅启用）
	Delete(ctx context.Context, id int64) error                                                   // 删除邮轮公司（软删除）
}

// CruiseRepository 定义邮轮的数据持久化接口。
type CruiseRepository interface {
	Create(ctx context.Context, cruise *Cruise) error                                                                                     // 创建邮轮
	Update(ctx context.Context, cruise *Cruise) error                                                                                     // 更新邮轮信息
	GetByID(ctx context.Context, id int64) (*Cruise, error)                                                                               // 根据 ID 查询邮轮
	List(ctx context.Context, companyID int64, keyword string, status *int16, sortBy string, page, pageSize int) ([]Cruise, int64, error) // 分页查询邮轮列表，可按公司/关键词/状态筛选
	ListPublic(ctx context.Context, companyID int64, keyword string, sortBy string, page, pageSize int) ([]Cruise, int64, error)          // 查询前台可见邮轮列表（仅启用）
	Delete(ctx context.Context, id int64) error                                                                                           // 删除邮轮（软删除）
}

// CabinTypeRepository 定义舱房类型的数据持久化接口。
type CabinTypeRepository interface {
	Create(ctx context.Context, cabinType *CabinType) error                                           // 创建舱房类型
	Update(ctx context.Context, cabinType *CabinType) error                                           // 更新舱房类型信息
	GetByID(ctx context.Context, id int64) (*CabinType, error)                                        // 根据 ID 查询舱房类型
	ListByCruise(ctx context.Context, cruiseID int64, page, pageSize int) ([]CabinType, int64, error) // 分页查询某邮轮下的舱房类型
	Delete(ctx context.Context, id int64) error                                                       // 删除舱房类型（软删除）
}

// CabinTypeCategoryRepository 定义舱型大类字典的数据持久化接口。
type CabinTypeCategoryRepository interface {
	Create(ctx context.Context, category *CabinTypeCategory) error     // 创建舱型大类
	Update(ctx context.Context, category *CabinTypeCategory) error     // 更新舱型大类
	GetByID(ctx context.Context, id int64) (*CabinTypeCategory, error) // 根据 ID 查询舱型大类
	List(ctx context.Context) ([]CabinTypeCategory, error)             // 查询舱型大类列表
	Delete(ctx context.Context, id int64) error                        // 删除舱型大类
}

// CabinTypeBindingRepository 定义舱型与邮轮绑定关系的数据持久化接口。
type CabinTypeBindingRepository interface {
	ReplaceCruiseBindings(ctx context.Context, cabinTypeID int64, cruiseIDs []int64) error // 覆盖写入舱型绑定的邮轮集合
	ListCruiseIDsByCabinType(ctx context.Context, cabinTypeID int64) ([]int64, error)      // 查询舱型绑定的邮轮 ID 列表
	ListCabinTypeIDsByCruise(ctx context.Context, cruiseID int64) ([]int64, error)         // 查询邮轮下绑定的舱型 ID 列表
	HasCabinTypesByCruise(ctx context.Context, cruiseID int64) (bool, error)               // 判断邮轮是否仍有关联舱型
}

// CabinTypeMediaRepository 定义舱型媒体资源的数据持久化接口。
type CabinTypeMediaRepository interface {
	Create(ctx context.Context, media *CabinTypeMedia) error                                  // 创建舱型媒体
	Update(ctx context.Context, media *CabinTypeMedia) error                                  // 更新舱型媒体
	GetByID(ctx context.Context, id int64) (*CabinTypeMedia, error)                           // 根据 ID 查询舱型媒体
	ListByCabinType(ctx context.Context, cabinTypeID int64) ([]CabinTypeMedia, error)         // 查询舱型媒体列表
	Delete(ctx context.Context, id int64) error                                               // 删除舱型媒体
	SetPrimary(ctx context.Context, cabinTypeID int64, mediaType string, mediaID int64) error // 设置主媒体并清除同类型其他主媒体
}

// VoyageCabinTypePriceRepository 定义航次舱型价格版本与当前态的数据持久化接口。
type VoyageCabinTypePriceRepository interface {
	CreateVersion(ctx context.Context, version *VoyageCabinTypePriceVersion) error                                                   // 创建价格版本
	UpsertCurrent(ctx context.Context, current *VoyageCabinTypeCurrent) error                                                        // 写入当前生效价格
	GetCurrent(ctx context.Context, voyageID, cabinTypeID int64) (*VoyageCabinTypeCurrent, error)                                    // 查询单个当前价
	GetCurrentAt(ctx context.Context, voyageID, cabinTypeID int64, at time.Time) (*VoyageCabinTypeCurrent, error)                    // 按生效时点查询单个当前价
	ListCurrentByVoyages(ctx context.Context, voyageIDs []int64) ([]VoyageCabinTypeCurrent, error)                                   // 查询多个航次的当前价
	ListCurrentByVoyagesAt(ctx context.Context, voyageIDs []int64, at time.Time) ([]VoyageCabinTypeCurrent, error)                   // 按生效时点查询多个航次当前价
	GetLatestVersionAt(ctx context.Context, voyageID, cabinTypeID int64, at time.Time) (*VoyageCabinTypePriceVersion, error)         // 按生效时点查询最新历史版本
	ListVersions(ctx context.Context, voyageID, cabinTypeID int64, page, pageSize int) ([]VoyageCabinTypePriceVersion, int64, error) // 查询历史版本
}

// FacilityCategoryRepository 定义设施分类的数据持久化接口。
type FacilityCategoryRepository interface {
	Create(ctx context.Context, category *FacilityCategory) error     // 创建设施分类
	Update(ctx context.Context, category *FacilityCategory) error     // 更新设施分类
	GetByID(ctx context.Context, id int64) (*FacilityCategory, error) // 根据 ID 查询设施分类
	List(ctx context.Context) ([]FacilityCategory, error)             // 查询所有设施分类
	Delete(ctx context.Context, id int64) error                       // 删除设施分类
}

// FacilityRepository 定义设施的数据持久化接口。
type FacilityRepository interface {
	Create(ctx context.Context, facility *Facility) error                                        // 创建设施
	Update(ctx context.Context, facility *Facility) error                                        // 更新设施
	GetByID(ctx context.Context, id int64) (*Facility, error)                                    // 根据 ID 查询设施
	ListByCruise(ctx context.Context, cruiseID int64) ([]Facility, error)                        // 查询某邮轮下的所有设施
	ListByCruiseAndCategory(ctx context.Context, cruiseID, categoryID int64) ([]Facility, error) // 查询某邮轮某分类下的设施
	Delete(ctx context.Context, id int64) error                                                  // 删除设施（软删除）
}

// ImageRepository 定义图片资源的多态存储接口。
type ImageRepository interface {
	Create(ctx context.Context, img *Image) error                                                // 创建图片记录
	ListByEntity(ctx context.Context, entityType string, entityID int64) ([]Image, error)        // 查询实体关联的图片列表
	DeleteByEntity(ctx context.Context, entityType string, entityID int64) error                 // 删除实体关联的全部图片（当前为物理删除）
	UpdateSortOrder(ctx context.Context, id int64, sortOrder int) error                          // 更新图片排序
	ReplaceImages(ctx context.Context, entityType string, entityID int64, images []*Image) error // 事务内替换实体图片
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
	Create(ctx context.Context, v *Voyage) error            // 创建航次
	Update(ctx context.Context, v *Voyage) error            // 更新航次信息
	GetByID(ctx context.Context, id int64) (*Voyage, error) // 根据 ID 查询航次
	List(ctx context.Context) ([]Voyage, error)             // 查询航次列表
	Delete(ctx context.Context, id int64) error             // 删除航次
}

// CabinSKUFilter 描述舱位商品的后台筛选条件。
type CabinSKUFilter struct {
	VoyageID    int64  // 航次 ID
	CabinTypeID int64  // 舱型 ID
	Status      *int16 // 状态（nil 表示不按状态筛选）
	Keyword     string // 舱位编号关键字
	Page        int    // 页码
	PageSize    int    // 每页条数
}

// CabinSKURepository 聚合了舱房产品的所有存储操作，
// 包括 SKU 的 CRUD、原子化库存调整和价格日历管理。
type CabinSKURepository interface {
	CreateSKU(ctx context.Context, s *CabinSKU) error                                 // 创建舱房 SKU
	UpdateSKU(ctx context.Context, s *CabinSKU) error                                 // 更新舱房 SKU
	GetSKUByID(ctx context.Context, id int64) (*CabinSKU, error)                      // 根据 ID 查询舱房 SKU
	ListSKUByVoyage(ctx context.Context, voyageID int64) ([]CabinSKU, error)          // 查询某航次下的所有舱房 SKU
	ListSKUFiltered(ctx context.Context, f CabinSKUFilter) ([]CabinSKU, int64, error) // 按条件分页查询舱房 SKU
	BatchUpdateStatus(ctx context.Context, ids []int64, status int16) error           // 批量更新舱房状态
	DeleteSKU(ctx context.Context, id int64) error                                    // 删除舱房 SKU
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

type TrendDataItem struct {
	Date   string
	Sales  int64
	Orders int64
}

type CabinRankingItem struct {
	CabinSKUID int64
	CabinName  string
	SoldCount  int64
	ViewCount  int64
}

type InventoryOverviewData struct {
	TotalCabins     int64
	LowStockCount   int64
	OutOfStockCount int64
}

type PageViewData struct {
	Page  string
	Views int64
}

// AnalyticsRepository 定义只读分析查询操作。
type AnalyticsRepository interface {
	TodaySales(ctx context.Context) (int64, error)
	WeeklyTrend(ctx context.Context) ([]int64, error)
	TodayOrderCount(ctx context.Context) (int64, error)
	Trend(ctx context.Context, days int) ([]TrendDataItem, error)
	CabinHotnessRanking(ctx context.Context, limit int) ([]CabinRankingItem, error)
	InventoryOverview(ctx context.Context) (*InventoryOverviewData, error)
	PageViewStats(ctx context.Context) ([]PageViewData, error)
}

// BookingStatusRepository 提供订单状态更新功能，
// 供支付回调服务确认订单支付。
type BookingStatusRepository interface {
	UpdateStatus(ctx context.Context, id int64, status string) error
}
