package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/leinonen/hexagonal-architecture-go/internal/domain"
	"github.com/leinonen/hexagonal-architecture-go/internal/errors"
	"github.com/leinonen/hexagonal-architecture-go/internal/ports"
)

type WeatherClient struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

func NewWeatherClient(apiKey string) ports.WeatherService {
	return &WeatherClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		apiKey:  apiKey,
		baseURL: "https://api.openweathermap.org/data/2.5",
	}
}

type weatherAPIResponse struct {
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Name string `json:"name"`
}

func (c *WeatherClient) GetWeather(ctx context.Context, city string) (*domain.Weather, error) {
	url := fmt.Sprintf("%s/weather?q=%s&appid=%s&units=metric", c.baseURL, city, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Sprintf("failed to create request: %v", err))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.NewExternalServiceError(fmt.Sprintf("failed to fetch weather: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.NewNotFoundError("city not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewExternalServiceError(fmt.Sprintf("weather API returned status: %d", resp.StatusCode))
	}

	var apiResp weatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, errors.NewInternalError(fmt.Sprintf("failed to decode response: %v", err))
	}

	weather := &domain.Weather{
		City:        apiResp.Name,
		Temperature: apiResp.Main.Temp,
		Humidity:    apiResp.Main.Humidity,
		WindSpeed:   apiResp.Wind.Speed,
	}

	if len(apiResp.Weather) > 0 {
		weather.Description = apiResp.Weather[0].Description
	}

	return weather, nil
}
