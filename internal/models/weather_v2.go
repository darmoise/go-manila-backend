package models

import (
	"encoding/xml"
	"go-manila-backend/internal/consts"
)

type WeatherV3Response struct {
	XMLName  xml.Name   `xml:"weather"`
	Location LocationV3 `xml:"location"`
	Current  CurrentV3  `xml:"current"`
	Forecast ForecastV3 `xml:"forecast"`
}

type LocationV3 struct {
	City    string `xml:"city,attr"`
	Country string `xml:"country,attr"`
	Code    string `xml:"code,attr"`
}

type CurrentV3 struct {
	Temperature TemperatureV3 `xml:"temperate"`
	FeelsLike   TemperatureV3 `xml:"feelslike"`
	WeatherText WeatherTextV3 `xml:"weathertext"`
	WeatherIcon WeatherIconV3 `xml:"weathericon"`
	WindSpeed   WindSpeedV3   `xml:"windspeed"`
	Humidity    HumidityV3    `xml:"humidity"`
	UVIndex     UVIndexV3     `xml:"uvindex"`
}

type TemperatureV3 struct {
	Value int    `xml:"value,attr"`
	Unit  string `xml:"unit,attr"`
}

type WeatherTextV3 struct {
	Value string `xml:"value,attr"`
}

type WeatherIconV3 struct {
	Value int `xml:"value,attr"`
}

type WindSpeedV3 struct {
	Value int    `xml:"value,attr"`
	Unit  string `xml:"unit,attr"`
}

type HumidityV3 struct {
	Value int `xml:"value,attr"`
}

type UVIndexV3 struct {
	Value int `xml:"value,attr"`
}

type ForecastV3 struct {
	Days []DayV3 `xml:"day"`
}

type DayV3 struct {
	Number      int           `xml:"number,attr"`
	Name        string        `xml:"name,attr"`
	Date        string        `xml:"date,attr"`
	High        TemperatureV3 `xml:"high"`
	Low         TemperatureV3 `xml:"low"`
	WeatherText WeatherTextV3 `xml:"weathertext"`
	WeatherIcon WeatherIconV3 `xml:"weathericon"`
}

func NewTestWeatherV3() WeatherV3Response {
	return WeatherV3Response{
		Location: LocationV3{
			City:    "Moscow",
			Country: "Russia",
			Code:    "EU|RU|MOW",
		},
		Current: CurrentV3{
			Temperature: TemperatureV3{Value: 22, Unit: "c"},
			FeelsLike:   TemperatureV3{Value: 20, Unit: "c"},
			WeatherText: WeatherTextV3{Value: "Clear"},
			WeatherIcon: WeatherIconV3{Value: 1},
			WindSpeed:   WindSpeedV3{Value: 12, Unit: "km/h"},
			Humidity:    HumidityV3{Value: 65},
			UVIndex:     UVIndexV3{Value: 4},
		},
		Forecast: ForecastV3{
			Days: []DayV3{
				{
					Number:      1,
					Name:        "Sun",
					Date:        "2026-07-19",
					High:        TemperatureV3{Value: 25, Unit: "c"},
					Low:         TemperatureV3{Value: 16, Unit: "c"},
					WeatherText: WeatherTextV3{Value: "Clear"},
					WeatherIcon: WeatherIconV3{Value: 1},
				},
				{
					Number:      2,
					Name:        "Mon",
					Date:        "2026-07-20",
					High:        TemperatureV3{Value: 24, Unit: "c"},
					Low:         TemperatureV3{Value: 15, Unit: "c"},
					WeatherText: WeatherTextV3{Value: "Partly Cloudy"},
					WeatherIcon: WeatherIconV3{Value: 2},
				},
				{
					Number:      3,
					Name:        "Tue",
					Date:        "2026-07-21",
					High:        TemperatureV3{Value: 20, Unit: "c"},
					Low:         TemperatureV3{Value: 14, Unit: "c"},
					WeatherText: WeatherTextV3{Value: "Rain"},
					WeatherIcon: WeatherIconV3{Value: 12},
				},
			},
		},
	}
}

func MarshalWeatherV3(r WeatherV3Response) ([]byte, error) {
	return xml.MarshalIndent(r, consts.IndentPrefix, consts.IndentValue)
}
