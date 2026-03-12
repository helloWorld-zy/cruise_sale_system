package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// PortCityOption 表示港口城市下拉选项的展示结构。
type PortCityOption struct {
	Label       string `json:"label"`                  // 显示标签，格式如 "上海（中国）"
	CityName    string `json:"city_name,omitempty"`    // 城市名称
	CountryName string `json:"country_name,omitempty"` // 国家名称
	IsSpecial   bool   `json:"is_special,omitempty"`   // 是否为特殊目的地（如自定义目的地）
}

// ResolvedPortCity 表示解析后的港口城市详细信息，包含地理坐标。
type ResolvedPortCity struct {
	Label       string   // 显示标签
	CityName    string   // 城市名称
	CountryName string   // 国家名称
	Latitude    *float64 // 纬度坐标
	Longitude   *float64 // 经度坐标
	IsSpecial   bool     // 是否为特殊目的地
}

// PortCityLookupService 定义港口城市查询服务的接口。
type PortCityLookupService interface {
	Search(ctx context.Context, keyword string) ([]PortCityOption, error)
	ResolveLabel(ctx context.Context, label string) (*ResolvedPortCity, error)
}

// PortCityServiceConfig 定义港口城市服务的配置参数。
type PortCityServiceConfig struct {
	Endpoint string        // 外部城市搜索服务 API 端点地址
	Timeout  time.Duration // HTTP 请求超时时间
}

// PortCityService 提供港口城市搜索和解析的业务逻辑。
type PortCityService struct {
	endpoint   string
	httpClient *http.Client
	customRepo CustomDestinationRepo // 自定义目的地仓储
}

// localPortCityEntry 定义本地港口城市词典的条目结构。
type localPortCityEntry struct {
	Label       string   // 显示标签
	CityName    string   // 城市名称
	CountryName string   // 国家名称
	Latitude    float64  // 纬度
	Longitude   float64  // 经度
	Keywords    []string // 搜索关键词列表
}

var localPortCityCatalog = []localPortCityEntry{}

// nominatimResult 定义 Nominatim 地理编码服务返回的结果结构。
type nominatimResult struct {
	Lat         string            `json:"lat"`         // 纬度
	Lon         string            `json:"lon"`         // 经度
	Name        string            `json:"name"`        // 地点名称
	NameDetails map[string]string `json:"namedetails"` // 名称详情（多语言）
	Address     struct {
		City         string `json:"city"`
		Town         string `json:"town"`
		Village      string `json:"village"`
		Municipality string `json:"municipality"`
		County       string `json:"county"`
		State        string `json:"state"`
		Province     string `json:"province"`
		Island       string `json:"island"`
		Suburb       string `json:"suburb"`
		Quarter      string `json:"quarter"`
		Country      string `json:"country"`
		CountryCode  string `json:"country_code"`
	} `json:"address"`
}

var countryNameByCode = map[string]string{
	"ar": "阿根廷",
	"cn": "中国",
	"jp": "日本",
	"kr": "韩国",
	"mx": "墨西哥",
	"us": "美国",
}

func NewPortCityService(cfg PortCityServiceConfig) *PortCityService {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 8 * time.Second
	}
	return &PortCityService{
		endpoint:   strings.TrimSpace(cfg.Endpoint),
		httpClient: &http.Client{Timeout: cfg.Timeout},
	}
}

// SetCustomDestinationRepo 注入自定义目的地仓储，以便搜索时合并自定义目的地结果。
func (s *PortCityService) SetCustomDestinationRepo(repo CustomDestinationRepo) *PortCityService {
	s.customRepo = repo
	return s
}

func (s *PortCityService) Search(ctx context.Context, keyword string) ([]PortCityOption, error) {
	trimmed := strings.TrimSpace(keyword)
	if trimmed == "" {
		return []PortCityOption{}, nil
	}
	items := make([]PortCityOption, 0, 8)
	items = append(items, searchLocalPortCities(trimmed)...)
	// 搜索自定义目的地（数据库）
	if s.customRepo != nil {
		customResults, _ := s.customRepo.SearchByKeyword(ctx, trimmed)
		for _, dest := range customResults {
			label := dest.Name
			if dest.Country != "" {
				label = fmt.Sprintf("%s（%s）", dest.Name, dest.Country)
			}
			items = appendUniquePortCityOption(items, PortCityOption{
				Label:       label,
				CityName:    dest.Name,
				CountryName: dest.Country,
			})
		}
	}
	if isSeaCruiseKeyword(trimmed) {
		items = appendUniquePortCityOption(items, PortCityOption{Label: "海上巡游", IsSpecial: true})
	}
	if utf8.RuneCountInString(trimmed) < 2 {
		return items, nil
	}
	results, err := s.searchRemote(ctx, trimmed)
	if err != nil {
		return items, nil
	}
	for _, result := range results {
		option := toPortCityOption(result)
		if option.Label == "" {
			continue
		}
		items = appendUniquePortCityOption(items, option)
	}
	return items, nil
}

func (s *PortCityService) ResolveLabel(ctx context.Context, label string) (*ResolvedPortCity, error) {
	trimmed := strings.TrimSpace(label)
	if trimmed == "" {
		return nil, nil
	}
	if trimmed == "海上巡游" {
		return &ResolvedPortCity{Label: trimmed, IsSpecial: true}, nil
	}
	if resolved := resolveLocalPortCity(trimmed); resolved != nil {
		return resolved, nil
	}
	// 解析 "名称（国家）" 格式并尝试匹配自定义目的地
	if s.customRepo != nil {
		name, country := parseLabelParts(trimmed)
		if name != "" && country != "" {
			if dest, err := s.customRepo.GetByLabel(ctx, name, country); err == nil && dest != nil {
				return &ResolvedPortCity{
					Label:       trimmed,
					CityName:    dest.Name,
					CountryName: dest.Country,
					Latitude:    dest.Latitude,
					Longitude:   dest.Longitude,
				}, nil
			}
		}
	}
	results, err := s.searchRemote(ctx, trimmed)
	if err != nil {
		return nil, err
	}
	for _, result := range results {
		resolved, err := toResolvedPortCity(result)
		if err != nil || resolved == nil {
			continue
		}
		if resolved.Label == trimmed {
			return resolved, nil
		}
	}
	for _, result := range results {
		resolved, err := toResolvedPortCity(result)
		if err == nil && resolved != nil {
			return resolved, nil
		}
	}
	return nil, nil
}

// parseLabelParts 解析 "名称（国家）" 格式的标签。
func parseLabelParts(label string) (string, string) {
	// 支持全角和半角括号
	for _, pair := range [][2]string{{"（", "）"}, {"(", ")"}} {
		open := strings.Index(label, pair[0])
		close := strings.Index(label, pair[1])
		if open > 0 && close > open {
			return strings.TrimSpace(label[:open]), strings.TrimSpace(label[open+len(pair[0]) : close])
		}
	}
	return "", ""
}

func (s *PortCityService) searchRemote(ctx context.Context, keyword string) ([]nominatimResult, error) {
	if s == nil || s.endpoint == "" {
		return []nominatimResult{}, nil
	}
	requestURL, err := url.Parse(s.endpoint)
	if err != nil {
		return nil, err
	}
	query := requestURL.Query()
	query.Set("q", keyword)
	query.Set("format", "jsonv2")
	query.Set("accept-language", "zh-CN")
	query.Set("addressdetails", "1")
	query.Set("namedetails", "1")
	query.Set("limit", "8")
	requestURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "cruise-sale-system/1.0")
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("city search returned status %d", resp.StatusCode)
	}
	var results []nominatimResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}
	return results, nil
}

func toPortCityOption(item nominatimResult) PortCityOption {
	cityName := resolveCityName(item)
	countryName := resolveCountryName(item)
	if cityName == "" || countryName == "" || !containsHan(cityName) {
		return PortCityOption{}
	}
	return PortCityOption{
		Label:       fmt.Sprintf("%s（%s）", cityName, countryName),
		CityName:    cityName,
		CountryName: countryName,
	}
}

func toResolvedPortCity(item nominatimResult) (*ResolvedPortCity, error) {
	option := toPortCityOption(item)
	if option.Label == "" {
		return nil, nil
	}
	lat, err := strconv.ParseFloat(strings.TrimSpace(item.Lat), 64)
	if err != nil {
		return nil, err
	}
	lon, err := strconv.ParseFloat(strings.TrimSpace(item.Lon), 64)
	if err != nil {
		return nil, err
	}
	return &ResolvedPortCity{
		Label:       option.Label,
		CityName:    option.CityName,
		CountryName: option.CountryName,
		Latitude:    &lat,
		Longitude:   &lon,
	}, nil
}

func resolveCityName(item nominatimResult) string {
	for _, key := range []string{"name:zh-Hans", "name:zh-CN", "name:zh_CN", "name:zh"} {
		trimmed := normalizeMultilingualValue(item.NameDetails[key])
		if trimmed != "" {
			return trimmed
		}
	}
	for _, candidate := range []string{item.Address.City, item.Address.Town, item.Address.Village, item.Address.Municipality, item.Address.County, item.Address.Island, item.Address.Suburb, item.Address.Quarter, item.Address.State, item.Address.Province, item.Name} {
		trimmed := normalizeMultilingualValue(candidate)
		if trimmed != "" && containsHan(trimmed) {
			return trimmed
		}
	}
	return ""
}

func resolveCountryName(item nominatimResult) string {
	if mapped := countryNameByCode[strings.ToLower(strings.TrimSpace(item.Address.CountryCode))]; mapped != "" {
		return mapped
	}
	return normalizeMultilingualValue(item.Address.Country)
}

func isSeaCruiseKeyword(keyword string) bool {
	for _, candidate := range []string{"海上", "巡游", "巡航"} {
		if strings.Contains(keyword, candidate) {
			return true
		}
	}
	return false
}

func searchLocalPortCities(keyword string) []PortCityOption {
	trimmed := normalizeSearchKeyword(keyword)
	if trimmed == "" {
		return []PortCityOption{}
	}
	prefixItems := make([]PortCityOption, 0, 6)
	fuzzyItems := make([]PortCityOption, 0, 6)
	for _, entry := range localPortCityCatalog {
		matchedPrefix := false
		matchedContains := false
		for _, candidate := range entry.Keywords {
			normalizedCandidate := normalizeSearchKeyword(candidate)
			if normalizedCandidate == "" {
				continue
			}
			if strings.HasPrefix(normalizedCandidate, trimmed) || strings.HasPrefix(trimmed, normalizedCandidate) {
				matchedPrefix = true
				break
			}
			if strings.Contains(normalizedCandidate, trimmed) || strings.Contains(trimmed, normalizedCandidate) {
				matchedContains = true
			}
		}
		if !matchedPrefix && !matchedContains {
			continue
		}
		option := PortCityOption{
			Label:       entry.Label,
			CityName:    entry.CityName,
			CountryName: entry.CountryName,
		}
		if matchedPrefix {
			prefixItems = appendUniquePortCityOption(prefixItems, option)
		} else {
			fuzzyItems = appendUniquePortCityOption(fuzzyItems, option)
		}
	}
	return append(prefixItems, fuzzyItems...)
}

func normalizeMultilingualValue(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}
	parts := strings.Split(trimmed, ";")
	for _, part := range parts {
		candidate := strings.TrimSpace(part)
		if candidate == "" {
			continue
		}
		if containsHan(candidate) {
			return candidate
		}
	}
	return strings.TrimSpace(parts[0])
}

func normalizeSearchKeyword(raw string) string {
	replacer := strings.NewReplacer("（", "(", "）", ")", " ", "", "-", "", "'", "", "’", "", "·", "")
	return strings.ToLower(strings.TrimSpace(replacer.Replace(raw)))
}

func containsHan(text string) bool {
	for _, r := range text {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

func resolveLocalPortCity(label string) *ResolvedPortCity {
	trimmed := strings.TrimSpace(label)
	for _, entry := range localPortCityCatalog {
		if entry.Label != trimmed {
			continue
		}
		lat := entry.Latitude
		lon := entry.Longitude
		return &ResolvedPortCity{
			Label:       entry.Label,
			CityName:    entry.CityName,
			CountryName: entry.CountryName,
			Latitude:    &lat,
			Longitude:   &lon,
		}
	}
	return nil
}

func appendUniquePortCityOption(items []PortCityOption, option PortCityOption) []PortCityOption {
	if option.Label == "" {
		return items
	}
	for _, item := range items {
		if item.Label == option.Label {
			return items
		}
	}
	return append(items, option)
}
