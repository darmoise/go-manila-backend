package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const bodyPreviewLimit = 4 * 1024 // 4 KB

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleAnyRequest(logger))

	server := &http.Server{
		Addr:         ":1917",
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
