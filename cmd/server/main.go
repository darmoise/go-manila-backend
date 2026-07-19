package main

import (
	"encoding/json"
	"go-manila-backend/internal/consts"
	"go-manila-backend/internal/geocoding"
	"go-manila-backend/internal/mappers"
	"go-manila-backend/internal/models"
	"go-manila-backend/internal/providers/nominatim"
	"go-manila-backend/internal/providers/openmeteo"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const bodyPreviewLimit = 4 * 1024 // 4 KB

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	nominatimClient := nominatim.NewClient()
	openMeteoClient := openmeteo.NewClient()

	mux := http.NewServeMux()
	//	mux.HandleFunc("/", handleAnyRequest(logger))
	mux.HandleFunc("GET /health", handleHealth(logger))

	// HTC legacy endpoints
	mux.HandleFunc("GET /widget/htc/forecast-data_v3.asp", handleHTCForecast(logger))
	mux.HandleFunc("GET /widget/htc/lat-lon-search.asp", handleHTCLatLonSearch(logger))
	mux.HandleFunc("GET /widget/htc2/weather-data.asp",
		handleAccuWeatherV3(logger, nominatimClient, openMeteoClient))

	// Alternative endpoints
	//mux.HandleFunc("GET /getweather", handleGetWeather(logger))
	mux.HandleFunc("GET /forecast-data_v3.asp", handleHTCForecast(logger))

	server := &http.Server{
		Addr:         ":" + resolvePort(),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	logger.Info("sniff server started", slog.String("addr", server.Addr))

	if err := server.ListenAndServe(); err != nil {
		logger.Error("sniff server failed", slog.Any("error", err))
		os.Exit(1)
	}
}

func handleHealth(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(models.HealthResponse{Status: "ok"}); err != nil {
			logger.Error("failed to encode health response", slog.Any("error", err))
		}
	}
}

func handleHTCForecast(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("HTC forecase request",
			slog.String("path", r.URL.Path),
			slog.String("query", r.URL.RawQuery),
			slog.String("user_agent", r.UserAgent()),
		)

		data := models.NewTestWeather()
		xmlBytes, err := models.MarshalWeather(data)

		if err != nil {
			logger.Error("failed to marshal weather", slog.Any("error", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", consts.ContentTypeXML)
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write(xmlBytes); err != nil {
			logger.Error("failed to write response", slog.Any("error", err))
		}
	}
}

func handleHTCLatLonSearch(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lat := r.URL.Query().Get(consts.ParamLat)
		lon := r.URL.Query().Get(consts.ParamLon)

		logger.Info("HTC lat/lon request",
			slog.String("path", r.URL.Path),
			slog.String("lat", lat),
			slog.String("lon", lon),
			slog.String("user_agent", r.UserAgent()),
		)

		data := models.NewTestLocationSearch()
		xmlBytes, err := models.MarshalLocationSearch(data)

		if err != nil {
			logger.Error("failed to marshal location", slog.Any("error", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", consts.ContentTypeXML)
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write(xmlBytes); err != nil {
			logger.Error("failed to write response", slog.Any("error", err))
		}
	}
}

func handleAccuWeatherV3(
	logger *slog.Logger,
	nominatimClient *nominatim.Client,
	openMeteoClient *openmeteo.Client,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawLocCode := r.URL.Query().Get(consts.ParamLocCode)
		metric := r.URL.Query().Get(consts.ParamMetric)

		logger.Info("AccuWeather V3 request",
			slog.String("path", r.URL.Path),
			slog.String("locCode", rawLocCode),
			slog.String("metric", metric),
			slog.String("user_agent", r.UserAgent()),
		)

		locCode, err := geocoding.ParseLocCode(rawLocCode)
		if err != nil {
			logger.Warn(
				"invalid locCode",
				slog.String("locCode", rawLocCode),
				slog.Any("error", err),
			)
			http.Error(w, "invalid locCode", http.StatusBadRequest)
			return
		}

		location, err := nominatimClient.SearchCity(
			r.Context(),
			locCode.CityName,
			locCode.CountryCode,
		)
		if err != nil {
			logger.Error(
				"failed to resolve location",
				slog.String("city", locCode.CityName),
				slog.String("country_code", locCode.CountryCode),
				slog.Any("error", err),
			)

			http.Error(w, "failed to resolve location", http.StatusBadGateway)
			return
		}

		data, err := openMeteoClient.FetchWeather(
			r.Context(),
			location.Latitude,
			location.Longitude,
			openmeteo.DefaultForecastDays,
		)

		if err != nil {
			logger.Error(
				"failed to fetcg weather V3",
				slog.Float64("lat", location.Latitude),
				slog.Float64("lon", location.Longitude),
				slog.Any("error", err),
			)
			http.Error(w, "internal server error", http.StatusBadGateway)
			return
		}

		response := mappers.ToWeatherV3(
			rawLocCode,
			location.Name,
			location.Country,
			data,
		)

		xmlBytes, err := models.MarshalWeatherV3(response)
		if err != nil {
			logger.Error(
				"failed to marshal weather V3",
				slog.Any("error", err),
			)

			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", consts.ContentTypeXML)
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write(xmlBytes); err != nil {
			logger.Error(
				"failed to write response",
				slog.Any("error", err),
			)
		}
	}
}

func handleAnyRequest(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyPreview, bodyTruncated, err := readBodyPreview(r.Body, bodyPreviewLimit)
		if err != nil {
			logger.Warn("failed to read request body", slog.Any("error", err))
		}

		queryParams, err := json.Marshal(r.URL.Query())
		if err != nil {
			queryParams = []byte("{}")
		}

		const staticWeatherXML = `<?xml version="1.0" encoding="UTF-8"?>
<weather>
  <location name="Moscow" country="Russia" />
  <localTime>2026-07-19 22:15:00</localTime>
  <current>
    <temp c="22" f="72" />
    <condition icon="1">Clear</condition>
  </current>
  <forecast>
    <day number="1">
      <name>Sun</name>
      <date>2026-07-19</date>
      <condition icon="1">Clear</condition>
      <high c="25" f="77" />
      <low c="16" f="61" />
      <winddir direction="NW">315</winddir>
      <windspd kph="12" mph="7" mps="3" />
      <uvi>4</uvi>
    </day>
    <day number="2">
      <name>Mon</name>
      <date>2026-07-20</date>
      <condition icon="2">Partly Cloudy</condition>
      <high c="24" f="75" />
      <low c="15" f="59" />
      <winddir direction="W">270</winddir>
      <windspd kph="10" mph="6" mps="3" />
      <uvi>5</uvi>
    </day>
    <day number="3">
      <name>Tue</name>
      <date>2026-07-21</date>
      <condition icon="12">Rain</condition>
      <high c="20" f="68" />
      <low c="14" f="57" />
      <winddir direction="SW">225</winddir>
      <windspd kph="14" mph="9" mps="4" />
      <uvi>2</uvi>
    </day>
  </forecast>
  <copyright>Open-Meteo</copyright>
  <usage>HTC Pocket PC Weather Forecast</usage>
</weather>
`

		logger.Info(
			"incoming request",
			slog.Time("timestamp", time.Now().UTC()),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("raw_query", r.URL.RawQuery),
			slog.String("query_params", string(queryParams)),
			slog.String("user_agent", r.UserAgent()),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("host", r.Host),
			slog.String("content_type", r.Header.Get("Content-Type")),
			slog.Int64("content_length", r.ContentLength),
			slog.String("body_preview", bodyPreview),
			slog.Bool("body_truncated", bodyTruncated),
		)

		if r.URL.Path == "/widget/htc/forecast-data_v3.asp" {
			w.Header().Set("Content-Type", "application/xml; charset=utf-8")
			w.WriteHeader(http.StatusOK)

			_, _ = w.Write([]byte(staticWeatherXML))
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte("go-manila-backend sniff server is running\n"))
	}
}

func readBodyPreview(body io.ReadCloser, limit int64) (string, bool, error) {
	if body == nil {
		return "", false, nil
	}
	defer body.Close()

	data, err := io.ReadAll(io.LimitReader(body, limit+1))
	if err != nil {
		return "", false, err
	}

	truncated := int64(len(data)) > limit
	if truncated {
		data = data[:limit]
	}

	return string(data), truncated, nil
}

func resolvePort() string {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
