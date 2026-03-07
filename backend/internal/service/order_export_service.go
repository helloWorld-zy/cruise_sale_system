package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
)

// OrderFilter 定义订单导出时的筛选条件。
type OrderFilter struct {
	Status     string  // 订单状态筛选
	Phone      string  // 手机号筛选
	RouteID    int64   // 航线 ID 筛选
	VoyageID   int64   // 航次 ID 筛选
	VoyageCode string  // 航次编码筛选
	CruiseName string  // 邮轮名称筛选（模糊匹配）
	Keyword    string  // 关键词筛选（订单号/手机号）
	StartDate  *string // 开始日期筛选
	EndDate    *string // 结束日期筛选
	BookingNo  string  // 订单号筛选
}

// OrderRepositoryFilter 定义订单仓储的筛选查询接口。
type OrderRepositoryFilter interface {
	ListWithFilter(ctx context.Context, filter OrderFilter) ([]domain.Booking, error) // 根据筛选条件查询订单
}

type orderExportContextKey string

const (
	orderExportPermissionKey  orderExportContextKey = "orderExportAllowed" // 导出权限上下文键
	defaultOrderExportMaxRows                       = 5000                 // 默认最大导出行数
)

var (
	// ErrOrderExportForbidden 表示当前用户无权导出订单。
	ErrOrderExportForbidden = errors.New("order export forbidden")
	// ErrOrderExportExceededLimit 表示导出行数超出限制。
	ErrOrderExportExceededLimit = errors.New("order export exceeded limit")
)

// WithOrderExportPermission 在上下文中注入导出权限标记。
func WithOrderExportPermission(ctx context.Context, allowed bool) context.Context {
	return context.WithValue(ctx, orderExportPermissionKey, allowed)
}

// OrderExportService 提供订单导出为 CSV/Excel 格式的服务。
type OrderExportService struct {
	repo OrderRepositoryFilter // 订单仓储
}

// NewOrderExportService 创建订单导出服务实例。
func NewOrderExportService(repo OrderRepositoryFilter) *OrderExportService {
	return &OrderExportService{repo: repo}
}

// ExportToExcel 将符合筛选条件的订单导出为 CSV 格式字节数组。
// 包含权限检查和导出数量限制。
func (s *OrderExportService) ExportToExcel(ctx context.Context, filter OrderFilter) ([]byte, error) {
	if !hasOrderExportPermission(ctx) {
		return nil, ErrOrderExportForbidden
	}
	orders, err := s.repo.ListWithFilter(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(orders) > defaultOrderExportMaxRows {
		return nil, ErrOrderExportExceededLimit
	}

	return generateExcelBytes(orders), nil
}

// generateExcelBytes 将订单列表转换为 CSV 格式的字节数组。
// 对每个单元格进行安全过滤，防止 CSV 注入攻击。
func generateExcelBytes(orders []domain.Booking) []byte {
	buf := bytes.NewBuffer(nil)
	writer := csv.NewWriter(buf)
	_ = writer.Write([]string{"booking_no", "phone", "voyage_code", "cruise_name", "status", "user_id", "voyage_id", "cabin_sku_id", "total_cents", "paid_cents", "created_at"})
	for _, order := range orders {
		bookingNo := order.BookingNo
		if bookingNo == "" {
			bookingNo = strconv.FormatInt(order.ID, 10)
		}
		_ = writer.Write([]string{
			sanitizeCSVCell(bookingNo),
			sanitizeCSVCell(order.Phone),
			sanitizeCSVCell(order.VoyageCode),
			sanitizeCSVCell(order.CruiseName),
			sanitizeCSVCell(order.Status),
			strconv.FormatInt(order.UserID, 10),
			strconv.FormatInt(order.VoyageID, 10),
			strconv.FormatInt(order.CabinSKUID, 10),
			strconv.FormatInt(order.TotalCents, 10),
			strconv.FormatInt(order.PaidCents, 10),
			order.CreatedAt.Format(time.RFC3339),
		})
	}
	writer.Flush()
	return buf.Bytes()
}

// hasOrderExportPermission 检查上下文中是否具有订单导出权限。
func hasOrderExportPermission(ctx context.Context) bool {
	value := ctx.Value(orderExportPermissionKey)
	allowed, ok := value.(bool)
	return ok && allowed
}

// sanitizeCSVCell 对 CSV 单元格值进行安全过滤，防止 CSV 注入攻击。
// 过滤规则：
// - 以制表符、回车符、换行符开头的值，前面加单引号
// - 以 =、+、-、@ 开头的值，前面加单引号（防止公式注入）
func sanitizeCSVCell(value string) string {
	if value == "" {
		return value
	}
	// 先检查原始值首字符是否为危险空白字符（TrimSpace 会移除这些字符）
	switch value[0] {
	case '\t', '\r', '\n':
		return fmt.Sprintf("'%s", value)
	}
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return value
	}
	switch trimmed[0] {
	case '=', '+', '-', '@':
		return fmt.Sprintf("'%s", value)
	default:
		return value
	}
}
