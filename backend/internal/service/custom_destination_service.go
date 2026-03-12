package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/cruisebooking/backend/internal/domain"
)

// CustomDestinationRepo 定义自定义目的地仓储需要实现的接口。
type CustomDestinationRepo interface {
	Create(ctx context.Context, dest *domain.CustomDestination) error
	Update(ctx context.Context, dest *domain.CustomDestination) error
	GetByID(ctx context.Context, id int64) (*domain.CustomDestination, error)
	List(ctx context.Context) ([]domain.CustomDestination, error)
	SearchByKeyword(ctx context.Context, keyword string) ([]domain.CustomDestination, error)
	GetByLabel(ctx context.Context, name, country string) (*domain.CustomDestination, error)
	UpsertByNameCountry(ctx context.Context, dest *domain.CustomDestination) error
	Delete(ctx context.Context, id int64) error
}

type CustomDestinationImportSummary struct {
	Imported int `json:"imported"`
}

// CustomDestinationService 自定义目的地业务层。
type CustomDestinationService struct {
	repo CustomDestinationRepo
}

// NewCustomDestinationService 创建自定义目的地服务实例。
func NewCustomDestinationService(repo CustomDestinationRepo) *CustomDestinationService {
	return &CustomDestinationService{repo: repo}
}

func (s *CustomDestinationService) Create(ctx context.Context, dest *domain.CustomDestination) error {
	return s.repo.Create(ctx, dest)
}

func (s *CustomDestinationService) Update(ctx context.Context, dest *domain.CustomDestination) error {
	return s.repo.Update(ctx, dest)
}

func (s *CustomDestinationService) GetByID(ctx context.Context, id int64) (*domain.CustomDestination, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CustomDestinationService) List(ctx context.Context) ([]domain.CustomDestination, error) {
	return s.repo.List(ctx)
}

func (s *CustomDestinationService) ExportCSV(ctx context.Context) ([]byte, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	buf := &strings.Builder{}
	writer := csv.NewWriter(buf)
	if err := writer.Write([]string{"name", "country", "latitude", "longitude", "keywords", "sort_order", "status", "description"}); err != nil {
		return nil, err
	}
	for _, item := range items {
		row := []string{
			strings.TrimSpace(item.Name),
			strings.TrimSpace(item.Country),
			formatFloatPointer(item.Latitude),
			formatFloatPointer(item.Longitude),
			item.Keywords,
			strconv.Itoa(item.SortOrder),
			strconv.Itoa(int(item.Status)),
			item.Description,
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	return []byte(buf.String()), nil
}

func (s *CustomDestinationService) ImportCSV(ctx context.Context, reader io.Reader) (*CustomDestinationImportSummary, error) {
	parsed := csv.NewReader(reader)
	records, err := parsed.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) <= 1 {
		return &CustomDestinationImportSummary{Imported: 0}, nil
	}
	summary := &CustomDestinationImportSummary{}
	for idx, row := range records[1:] {
		if len(row) < 8 {
			return nil, fmt.Errorf("row %d has insufficient columns", idx+2)
		}
		name := strings.TrimSpace(row[0])
		country := strings.TrimSpace(row[1])
		if name == "" || country == "" {
			return nil, fmt.Errorf("row %d requires name and country", idx+2)
		}
		latitude, err := parseRequiredFloat(row[2], idx+2, "latitude")
		if err != nil {
			return nil, err
		}
		longitude, err := parseRequiredFloat(row[3], idx+2, "longitude")
		if err != nil {
			return nil, err
		}
		sortOrder, err := parseOptionalInt(row[5], 0)
		if err != nil {
			return nil, fmt.Errorf("row %d invalid sort_order: %w", idx+2, err)
		}
		statusValue, err := parseOptionalInt(row[6], 1)
		if err != nil {
			return nil, fmt.Errorf("row %d invalid status: %w", idx+2, err)
		}
		dest := &domain.CustomDestination{
			Name:        name,
			Country:     country,
			Latitude:    &latitude,
			Longitude:   &longitude,
			Keywords:    strings.TrimSpace(row[4]),
			SortOrder:   sortOrder,
			Status:      int16(statusValue),
			Description: strings.TrimSpace(row[7]),
		}
		if err := s.repo.UpsertByNameCountry(ctx, dest); err != nil {
			return nil, err
		}
		summary.Imported++
	}
	return summary, nil
}

func (s *CustomDestinationService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func formatFloatPointer(value *float64) string {
	if value == nil {
		return ""
	}
	return strconv.FormatFloat(*value, 'f', -1, 64)
}

func parseRequiredFloat(raw string, rowNo int, field string) (float64, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return 0, fmt.Errorf("row %d missing %s", rowNo, field)
	}
	value, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return 0, fmt.Errorf("row %d invalid %s: %w", rowNo, field, err)
	}
	return value, nil
}

func parseOptionalInt(raw string, fallback int) (int, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return fallback, nil
	}
	value, err := strconv.Atoi(trimmed)
	if err != nil {
		return 0, err
	}
	return value, nil
}
