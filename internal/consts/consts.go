package consts

import "time"

// Query params
const (
	ParamLat     = "lat"
	ParamLon     = "lon"
	ParamLocCode = "locCode"
	ParamMetric  = "metric"
)

// Headers
const (
	ContentTypeXML = "application/xml; charset=utf-8"
)

// Marshaling
const (
	IndentPrefix = ""
	IndentValue  = " "
)

const (
	OpenMeteoBaseURL = "https://api.open-meteo.com/v1/forecast"
	NominatimBaseURL = "https://nominatim.openstreetmap.org"
)

const (
	HTTPClientTimeout  = 10 * time.Second
	ServerReadTimeout  = 5 * time.Second
	ServerWriteTimeout = 10 * time.Second
	ServerIdleTimeout  = 30 * time.Second
)

const (
	DefaultPort         = "8080"
	DefaultMetric       = "celsius"
	DefaultForecastDays = 3
	DefaultCityName     = "Moscow"
	DefaultCountry      = "Russia"
	DefaultCopyright    = "Open-Meteo"
	DefaultUsage        = "HTC Pocket PC Weather Forecast"
)

const (
	UserAgent = "go-manila-backend/1.0 (https://github.com/darmoise/go-manila-backend)"
)
