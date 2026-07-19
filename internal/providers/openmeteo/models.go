package openmeteo

type OpenMeteoResponse struct {
	Latitude  float64        `json:"latitude"`
	Longitude float64        `json:"longitude"`
	Timezone  string         `json:"timezone"`
	Current   CurrentWeather `json:"current"`
	Daily     DailyWeather   `json:"daily"`
}

type CurrentWeather struct {
	Time                string  `json:"time"`
	Temperature2m       float64 `json:"temperature_2m"`
	ApparentTemperature float64 `json:"apparent_temperature"`
	WeatherCode         int     `json:"weather_code"`
	Rain                float64 `json:"rain"`
	Snowfall            float64 `json:"snowfall"`
	CloudCover          int     `json:"cloud_cover"`
	WindSpeed10m        float64 `json:"wind_speed_10m"`
	RelativeHumidity2m  int     `json:"relative_humidity_2m"`
}

type DailyWeather struct {
	Time                     []string  `json:"time"`
	Temperature2mMax         []float64 `json:"temperature_2m_max"`
	Temperature2mMin         []float64 `json:"temperature_2m_min"`
	WeatherCode              []int     `json:"weather_code"`
	UVIndexMax               []float64 `json:"uv_index_max"`
	RainSum                  []float64 `json:"rain_sum"`
	SnowfallSum              []float64 `json:"snowfall_sum"`
	WindSpeed10mMax          []float64 `json:"wind_speed_10m_max"`
	WindDirection10mDominant []int     `json:"wind_direction_10m_dominant"`
}
