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

type OrderFilter struct {
	Status     string
	Phone      string
	RouteID    int64
	VoyageID   int64
	VoyageCode string
	CruiseName string
	Keyword    string
	StartDate  *string
	EndDate    *string
	BookingNo  string
}

type OrderRepositoryFilter interface {
	ListWithFilter(ctx context.Context, filter OrderFilter) ([]domain.Booking, error)
}

type orderExportContextKey string

const (
	orderExportPermissionKey  orderExportContextKey = "orderExportAllowed"
	defaultOrderExportMaxRows                       = 5000
)

var (
	ErrOrderExportForbidden     = errors.New("order export forbidden")
	ErrOrderExportExceededLimit = errors.New("order export exceeded limit")
)

// WithOrderExportPermission 在上下文中注入导出权限。
func WithOrderExportPermission(ctx context.Context, allowed bool) context.Context {
	return context.WithValue(ctx, orderExportPermissionKey, allowed)
}

type OrderExportService struct {
	repo OrderRepositoryFilter
}

func NewOrderExportService(repo OrderRepositoryFilter) *OrderExportService {
	return &OrderExportService{repo: repo}
}

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

func hasOrderExportPermission(ctx context.Context) bool {
	value := ctx.Value(orderExportPermissionKey)
	allowed, ok := value.(bool)
	return ok && allowed
}

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
