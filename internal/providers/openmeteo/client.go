package openmeteo

import (
	"context"
	"encoding/json"
	"fmt"
	"go-manila-backend/internal/consts"
	"log"
	"net/http"
	"net/url"
)

const (
	// MinForecastDays is the minimum number of forecast days allowed
	MinForecastDays = 1
	// MaxForecastDays is the maximum number of forecast days allowed (Open-Meteo API limit)
	MaxForecastDays = 16
	// DefaultForecastDays is used when days parameter is not specified
	DefaultForecastDays = 3
)

// Geographic limits
const (
	MinLatitude  = -90.0
	MaxLatitude  = 90.0
	MinLongitude = -180.0
	MaxLongitude = 180.0
)

// Query parameters for Open-Meteo API
const (
	queryParamLatitude  = "latitude"
	queryParamLongitude = "longitude"
	queryParamCurrent   = "current"
	queryParamDaily     = "daily"
	queryParamTimezone  = "timezone"
	queryParamDays      = "forecast_days"
)

// Fields for current weather
const (
	currentFields = "temperature_2m,weather_code,rain,snowfall,cloud_cover,wind_speed_10m,relative_humidity_2m"
)

// Fields for daily forecast
const (
	dailyFields = "temperature_2m_max,temperature_2m_min,weather_code,uv_index_max,rain_sum,snowfall_sum,wind_speed_10m_max,wind_direction_10m_dominant"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: consts.HTTPClientTimeout,
		},
		baseURL: consts.OpenMeteoBaseURL,
	}
}

func NewClientWithHTTP(httpClient *http.Client, baseURL string) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

// FetchWeather fetches weather forecast for given coordinates
// Parameters:
//   - ctx: context for cancellation/timeouts
//   - lat: latitude (-90 to 90)
//   - lon: longitude (-180 to 180)
//   - days: number of forecast days (1-16)
func (c *Client) FetchWeather(
	ctx context.Context,
	lat, lon float64,
	days int,
) (*OpenMeteoResponse, error) {
	if days < MinForecastDays || days > MaxForecastDays {
		return nil, fmt.Errorf(
			"days must be between %d and %d, got %d",
			MinForecastDays,
			MaxForecastDays,
			days,
		)
	}

	if lat < MinLatitude || lat > MaxLatitude {
		return nil, fmt.Errorf(
			"invalid latitude: %.4f (must be between %.1f and %.1f)",
			lat,
			MinLatitude,
			MaxLatitude,
		)
	}

	if lon < MinLongitude || lon > MaxLongitude {
		return nil, fmt.Errorf(
			"invalid longitude: %.4f  (must be between %.1f and %.1f)",
			lon,
			MinLongitude,
			MaxLongitude,
		)
	}

	// Build query parameters
	params := url.Values{}
	params.Set(queryParamLatitude, fmt.Sprintf("%.4f", lat))
	params.Set(queryParamLongitude, fmt.Sprintf("%.4f", lon))
	params.Set(queryParamCurrent, currentFields)
	params.Set(queryParamDaily, dailyFields)
	params.Set(queryParamTimezone, "auto")
	params.Set(queryParamDays, fmt.Sprintf("%d", days))

	reqURL := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", consts.UserAgent)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather: %w", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var result OpenMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
