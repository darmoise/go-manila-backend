package mappers

import (
	"go-manila-backend/internal/models"
	"go-manila-backend/internal/providers/openmeteo"
	"go-manila-backend/internal/util"
)

func ToWeatherV3(
	locCode string,
	city string,
	country string,
	data *openmeteo.OpenMeteoResponse,
) models.WeatherV3Response {
	currentText, currentIcon := weatherCondition(data.Current.WeatherCode)
	dayCount := forecastDayCount(data.Daily)
	days := make([]models.DayV3, 0, dayCount)

	for i := range dayCount {
		text, icon := weatherCondition(data.Daily.WeatherCode[i])
		date := data.Daily.Time[i]

		days = append(days, models.DayV3{
			Number: i + 1,
			Name:   util.DayName(date),
			Date:   date,
			High: models.TemperatureV3{
				Value: util.RoundToInt(data.Daily.Temperature2mMax[i]),
				Unit:  util.TemperatureUnit,
			},
			Low: models.TemperatureV3{
				Value: util.RoundToInt(data.Daily.Temperature2mMin[i]),
				Unit:  util.TemperatureUnit,
			},
			WeatherText: models.WeatherTextV3{
				Value: text,
			},
			WeatherIcon: models.WeatherIconV3{
				Value: icon,
			},
		})
	}

	return models.WeatherV3Response{
		Location: models.LocationV3{
			City:    city,
			Country: country,
			Code:    locCode,
		},
		Current: models.CurrentV3{
			Temperature: models.TemperatureV3{
				Value: util.RoundToInt(data.Current.Temperature2m),
				Unit:  util.TemperatureUnit,
			},
			FeelsLike: models.TemperatureV3{
				Value: util.RoundToInt(data.Current.ApparentTemperature),
				Unit:  util.TemperatureUnit,
			},
			WeatherText: models.WeatherTextV3{Value: currentText},
			WeatherIcon: models.WeatherIconV3{Value: currentIcon},
			WindSpeed: models.WindSpeedV3{
				Value: util.RoundToInt(data.Current.WindSpeed10m),
				Unit:  util.KmPh,
			},
			Humidity: models.HumidityV3{Value: data.Current.RelativeHumidity2m},
			UVIndex:  models.UVIndexV3{Value: util.RoundToInt(util.FirstOrZero(data.Daily.UVIndexMax))},
		},
		Forecast: models.ForecastV3{Days: days},
	}
}

func forecastDayCount(daily openmeteo.DailyWeather) int {
	c := len(daily.Time)

	c = min(c, len(daily.Temperature2mMax))
	c = min(c, len(daily.Temperature2mMin))
	c = min(c, len(daily.WeatherCode))

	return c
}

func weatherCondition(code int) (text string, icon int) {
	switch code {
	case 0:
		return "Sunny", 1
	case 1:
		return "Mostly Sunny", 2
	case 2:
		return "Partly Cloudy", 3
	case 3:
		return "Cloudy", 7
	case 45, 48:
		return "Fog", 11
	case 51, 53, 55:
		return "Showers", 12
	case 56, 57:
		return "Freezing Drizzle", 25
	case 61, 63, 65:
		return "Rain", 18
	case 66, 67:
		return "Freezing Rain", 29
	case 71, 73, 75, 77:
		return "Snow", 22
	case 80, 81, 82:
		return "Showers", 12
	case 85, 86:
		return "Snow Showers", 22
	case 95, 96, 99:
		return "Thunderstorms", 15
	default:
		return "Unknown", 7
	}
}
