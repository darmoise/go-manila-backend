package models

import (
	"encoding/xml"
	"go-manila-backend/internal/consts"
)

type LocationSearchResponse struct {
	XMLName  xml.Name       `xml:"locations"`
	Location LocationResult `xml:"location"`
}

type LocationResult struct {
	Name    string `xml:"name,attr"`
	Country string `xml:"country,attr"`
	Lat     string `xml:"lat,attr"`
	Lon     string `xml:"lon,attr"`
}

func NewTestLocationSearch() LocationSearchResponse {
	return LocationSearchResponse{
		XMLName: xml.Name{},
		Location: LocationResult{
			Name:    "Moscow",
			Country: "Russia",
			Lat:     "55.7558",
			Lon:     "37.6173",
		},
	}
}

func MarshalLocationSearch(r LocationSearchResponse) ([]byte, error) {
	return xml.MarshalIndent(r, consts.IndentPrefix, consts.IndentValue)
}
