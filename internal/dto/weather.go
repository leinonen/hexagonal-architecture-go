package dto

import "github.com/leinonen/hexagonal-architecture-go/internal/domain"

type WeatherResponseDTO struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
	Humidity    int     `json:"humidity"`
	WindSpeed   float64 `json:"wind_speed"`
}

func ToWeatherResponseDTO(weather *domain.Weather) *WeatherResponseDTO {
	if weather == nil {
		return nil
	}

	return &WeatherResponseDTO{
		City:        weather.City,
		Temperature: weather.Temperature,
		Description: weather.Description,
		Humidity:    weather.Humidity,
		WindSpeed:   weather.WindSpeed,
	}
}
