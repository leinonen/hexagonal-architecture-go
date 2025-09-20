package ports

import (
	"context"

	"github.com/leinonen/hexagonal-architecture-go/internal/domain"
)

type WeatherService interface {
	GetWeather(ctx context.Context, city string) (*domain.Weather, error)
}
