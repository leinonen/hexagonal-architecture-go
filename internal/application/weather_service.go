package application

import (
	"context"

	"github.com/leinonen/hexagonal-architecture-go/internal/domain"
	"github.com/leinonen/hexagonal-architecture-go/internal/errors"
	"github.com/leinonen/hexagonal-architecture-go/internal/ports"
)

type WeatherService struct {
	weatherClient ports.WeatherService
	userRepo      ports.UserRepository
}

func NewWeatherService(weatherClient ports.WeatherService, userRepo ports.UserRepository) *WeatherService {
	return &WeatherService{
		weatherClient: weatherClient,
		userRepo:      userRepo,
	}
}

func (s *WeatherService) GetWeatherForUser(ctx context.Context, userID, city string) (*domain.Weather, error) {
	if userID == "" {
		return nil, errors.NewValidationError("user ID is required")
	}
	if city == "" {
		return nil, errors.NewValidationError("city is required")
	}

	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NewUnauthorizedError("user not found")
		}
		return nil, err
	}

	weather, err := s.weatherClient.GetWeather(ctx, city)
	if err != nil {
		return nil, err
	}

	return weather, nil
}

func (s *WeatherService) GetWeather(ctx context.Context, city string) (*domain.Weather, error) {
	if city == "" {
		return nil, errors.NewValidationError("city is required")
	}

	return s.weatherClient.GetWeather(ctx, city)
}
