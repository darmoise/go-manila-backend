package models

import (
	"encoding/xml"
	"go-manila-backend/internal/consts"
)

type WeatherResponse struct {
	XMLName   xml.Name `xml:"weather"`
	Location  Location `xml:"location"`
	LocalTime string   `xml:"localTime"`
	Current   Current  `xml:"current"`
	Forecast  Forecast `xml:"forecast"`
	Copyright string   `xml:"copyright"`
	Usage     string   `xml:"usage"`
}

type Location struct {
	Name    string `xml:"name,attr"`
	Country string `xml:"country,attr"`
}

type Current struct {
	Temp      Temp      `xml:"temp"`
	Condition Condition `xml:"condition"`
}

type Temp struct {
	C int `xml:"c,attr"`
	F int `xml:"f,attr"`
}

type Condition struct {
	Icon int    `xml:"icon,attr"`
	Text string `xml:",chardata"`
}

type Forecast struct {
	Days []Day `xml:"day"`
}

type Day struct {
	Number    int       `xml:"number,attr"`
	Name      string    `xml:"name"`
	Date      string    `xml:"date"`
	Condition Condition `xml:"condition"`
	High      Temp      `xml:"high"`
	Low       Temp      `xml:"low"`
	WindDir   WindDir   `xml:"winddir"`
	WindSpd   WindSpd   `xml:"windspd"`
	// UVI represents the Ultraviolet Index, a measure of the intensity of UV radiation.
	UVI int `xml:"uvi"`
}

type WindDir struct {
	Direction string `xml:"direction,attr"`
	Degrees   int    `xml:",chardata"`
}

type WindSpd struct {
	Kph int `xml:"kph,attr"`
	Mph int `xml:"mph,attr"`
	Mps int `xml:"mps,attr"`
}

func NewTestWeather() WeatherResponse {
	return WeatherResponse{
		XMLName:   xml.Name{},
		Location:  Location{Name: "Moscow", Country: "Russia"},
		LocalTime: "2026-07-19 22:15:00",
		Current: Current{
			Temp:      Temp{C: 22, F: 72},
			Condition: Condition{Icon: 1, Text: "Clear"},
		},
		Forecast: Forecast{
			Days: []Day{
				{
					Number:    1,
					Name:      "Sun",
					Date:      "2026-07-19",
					Condition: Condition{Icon: 1, Text: "Clear"},
					High:      Temp{C: 25, F: 77},
					Low:       Temp{C: 16, F: 61},
					WindDir:   WindDir{Direction: "NW", Degrees: 315},
					WindSpd:   WindSpd{Kph: 12, Mph: 7, Mps: 3},
					UVI:       4,
				},
			},
		},
		Copyright: "Open-Meteo",
		Usage:     "HTC Pocket PC Weather Forecast",
	}
}

func MarshalWeather(r WeatherResponse) ([]byte, error) {
	return xml.MarshalIndent(r, consts.IndentPrefix, consts.IndentValue)
}
