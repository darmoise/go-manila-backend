package nominatim

import (
	"context"
	"encoding/json"
	"fmt"
	"go-manila-backend/internal/consts"
	"go-manila-backend/internal/geocoding"
	. "go-manila-backend/internal/util"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	queryCity           = "city"
	queryCountryCodes   = "countrycodes"
	queryFormat         = "format"
	queryLimit          = "limit"
	queryAddressDetails = "addressdetails"
)

type searchResult struct {
	Lat     string  `json:"lat"`
	Lon     string  `json:"lon"`
	Address address `json:"address"`
}

type address struct {
	City    string `json:"city"`
	Town    string `json:"town"`
	Village string `json:"village"`
	Country string `json:"country"`
}

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: consts.HTTPClientTimeout,
		},
		baseURL: consts.NominatimBaseURL,
	}
}

func NewClientWithHTTP(httpClient *http.Client, baseURL string) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

// SearchCity searches for a city by name and country code.
// Parameters:
//   - ctx: context for cancellation and timeouts
//   - city: city name to search for
//   - countryCode: ISO 3166-1 alpha-2 country code
//
// Returns the matching geocoding location or an error if the request fails
// or no suitable location is found.
func (c *Client) SearchCity(
	ctx context.Context,
	city string,
	countryCode string,
) (geocoding.Location, error) {
	if IsEmpty(city) {
		return geocoding.Location{}, fmt.Errorf("city is required")
	}

	if IsEmpty(countryCode) {
		return geocoding.Location{}, fmt.Errorf("country code is required")
	}

	params := url.Values{}
	params.Set(queryCity, city)
	params.Set(queryCountryCodes, strings.ToLower(countryCode))
	params.Set(queryFormat, "jsonv2")
	params.Set(queryLimit, "1")
	params.Set(queryAddressDetails, "1")

	reqURL := fmt.Sprintf("%s/search?%s", c.baseURL, params.Encode())

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		reqURL,
		nil,
	)

	if err != nil {
		return geocoding.Location{}, fmt.Errorf(
			"failed to create Nominatum request: %w",
			err,
		)
	}

	req.Header.Set("User-Agent", consts.UserAgent)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return geocoding.Location{}, fmt.Errorf(
			"failed to search location: %w",
			err,
		)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return geocoding.Location{}, fmt.Errorf(
			"unexpected Nominatum status code: %d",
			resp.StatusCode,
		)
	}

	// Parse response
	var results []searchResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return geocoding.Location{}, fmt.Errorf("failed to decode Nominatum response: %w", err)
	}

	if len(results) == 0 {
		return geocoding.Location{}, fmt.Errorf(
			"location not found: city=%q, countryCode=%q",
			city,
			countryCode,
		)
	}

	result := results[0]

	lat, err := strconv.ParseFloat(result.Lat, 64)
	if err != nil {
		return geocoding.Location{}, fmt.Errorf(
			"failed to parse latitude %q: %w",
			result.Lat,
			err,
		)
	}

	lon, err := strconv.ParseFloat(result.Lon, 64)
	if err != nil {
		return geocoding.Location{}, fmt.Errorf(
			"failed to parse longitude %q: %w",
			result.Lon,
			err,
		)
	}

	return geocoding.Location{
		Name:      FirstNonEmpty(result.Address.City, result.Address.Town, result.Address.Village, city),
		Country:   result.Address.Country,
		Latitude:  lat,
		Longitude: lon,
	}, nil
}
