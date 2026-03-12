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

	"github.com/cruisebooking/backend/internal/domain"
)

type SeaRouteClientConfig struct {
	Endpoint     string
	Timeout      time.Duration
	ResolutionKM int
}

type SeaRouteClient struct {
	endpoint     string
	resolutionKM int
	httpClient   *http.Client
}

type seaRouteResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Dist    float64         `json:"dist"`
	Geom    json.RawMessage `json:"geom"`
}

type seaRouteGeometryEnvelope struct {
	Type string `json:"type"`
}

type routeStop struct {
	city      string
	latitude  float64
	longitude float64
}

var knownPortCoordinates = map[string][2]float64{
	"上海":  {31.2304, 121.4737},
	"天津":  {39.0842, 117.2009},
	"大连":  {38.9140, 121.6147},
	"青岛":  {36.0671, 120.3826},
	"厦门":  {24.4798, 118.0894},
	"香港":  {22.3193, 114.1694},
	"福冈":  {33.5902, 130.4017},
	"长崎":  {32.7503, 129.8777},
	"鹿儿岛": {31.5966, 130.5571},
	"济州":  {33.4996, 126.5312},
	"西归浦": {33.2541, 126.5601},
	"釜山":  {35.1796, 129.0756},
	"仁川":  {37.4563, 126.7052},
	"横滨":  {35.4437, 139.6380},
	"神户":  {34.6901, 135.1955},
	"大阪":  {34.6937, 135.5023},
}

var seaCruisePlaceholders = []string{"海上巡游", "海上巡航", "海上观光", "巡游日", "海上"}

func NewSeaRouteClient(cfg SeaRouteClientConfig) *SeaRouteClient {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 8 * time.Second
	}
	if cfg.ResolutionKM <= 0 {
		cfg.ResolutionKM = 20
	}
	return &SeaRouteClient{
		endpoint:     strings.TrimSpace(cfg.Endpoint),
		resolutionKM: cfg.ResolutionKM,
		httpClient:   &http.Client{Timeout: cfg.Timeout},
	}
}

func (c *SeaRouteClient) BuildVoyageRouteMap(ctx context.Context, itineraries []domain.VoyageItinerary) (*domain.VoyageRouteMap, error) {
	if c == nil || c.endpoint == "" {
		return nil, nil
	}
	stops := normalizeVoyageRouteStops(itineraries)
	if len(stops) < 2 {
		return nil, nil
	}

	coordinates := make([][][]float64, 0, len(stops)-1)
	totalDistance := 0.0
	for index := 0; index < len(stops)-1; index += 1 {
		payload, err := c.fetchSeaRoutePayload(ctx, stops[index], stops[index+1])
		if err != nil {
			return nil, err
		}
		if payload == nil {
			return nil, nil
		}
		segmentCoords, err := normalizeSeaRouteGeometry(payload.Geom)
		if err != nil {
			return nil, err
		}
		coordinates = append(coordinates, segmentCoords...)
		totalDistance += payload.Dist
	}
	if len(coordinates) == 0 {
		return nil, nil
	}
	return &domain.VoyageRouteMap{
		Provider:     "searoute",
		GeometryType: "MultiLineString",
		Coordinates:  coordinates,
		DistanceKM:   totalDistance,
		ResolutionKM: c.resolutionKM,
	}, nil
}

func normalizeVoyageRouteStops(itineraries []domain.VoyageItinerary) []routeStop {
	stops := make([]routeStop, 0, len(itineraries))
	for _, item := range itineraries {
		if isSeaCruiseStop(item.City, item.Summary) {
			continue
		}
		lat, lon, ok := resolveItineraryCoordinate(item)
		if !ok {
			return nil
		}
		if len(stops) > 0 {
			previous := stops[len(stops)-1]
			if previous.city == item.City && previous.latitude == lat && previous.longitude == lon {
				continue
			}
		}
		stops = append(stops, routeStop{city: item.City, latitude: lat, longitude: lon})
	}
	return stops
}

func resolveItineraryCoordinate(item domain.VoyageItinerary) (float64, float64, bool) {
	if item.Latitude != nil && item.Longitude != nil {
		return *item.Latitude, *item.Longitude, true
	}
	normalized := normalizePortName(item.City)
	if point, ok := knownPortCoordinates[normalized]; ok {
		return point[0], point[1], true
	}
	return 0, 0, false
}

func normalizePortName(city string) string {
	trimmed := strings.TrimSpace(city)
	if index := strings.Index(trimmed, "（"); index >= 0 {
		trimmed = trimmed[:index]
	}
	if index := strings.Index(trimmed, "("); index >= 0 {
		trimmed = trimmed[:index]
	}
	return strings.TrimSpace(strings.NewReplacer("港口", "", "港", "", "市", "").Replace(trimmed))
}

func isSeaCruiseStop(city string, summary string) bool {
	normalized := normalizePortName(city)
	for _, candidate := range seaCruisePlaceholders {
		if strings.Contains(normalized, candidate) || strings.Contains(summary, candidate) {
			return true
		}
	}
	return false
}

func (c *SeaRouteClient) buildURL(from routeStop, to routeStop) (string, error) {
	base, err := url.Parse(c.endpoint)
	if err != nil {
		return "", err
	}
	query := base.Query()
	query.Set("ser", "rou")
	query.Set("g", "1")
	query.Set("d", "1")
	query.Set("res", strconv.Itoa(c.resolutionKM))
	query.Set("opos", fmt.Sprintf("%.6f,%.6f", from.longitude, from.latitude))
	query.Set("dpos", fmt.Sprintf("%.6f,%.6f", to.longitude, to.latitude))
	base.RawQuery = query.Encode()
	return base.String(), nil
}

func (c *SeaRouteClient) fetchSeaRoutePayload(ctx context.Context, from routeStop, to routeStop) (*seaRouteResponse, error) {
	requestURL, err := c.buildURL(from, to)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var payload seaRouteResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if strings.ToLower(payload.Status) != "ok" {
		if strings.ToLower(payload.Status) == "empty" {
			return nil, nil
		}
		return nil, fmt.Errorf("searoute returned status %s: %s", payload.Status, payload.Message)
	}
	return &payload, nil
}

func normalizeSeaRouteGeometry(raw json.RawMessage) ([][][]float64, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var envelope seaRouteGeometryEnvelope
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, err
	}
	switch envelope.Type {
	case "LineString":
		var line struct {
			Coordinates [][]float64 `json:"coordinates"`
		}
		if err := json.Unmarshal(raw, &line); err != nil {
			return nil, err
		}
		return [][][]float64{line.Coordinates}, nil
	case "MultiLineString":
		var multi struct {
			Coordinates [][][]float64 `json:"coordinates"`
		}
		if err := json.Unmarshal(raw, &multi); err != nil {
			return nil, err
		}
		return multi.Coordinates, nil
	default:
		return nil, fmt.Errorf("unsupported searoute geometry type: %s", envelope.Type)
	}
}
