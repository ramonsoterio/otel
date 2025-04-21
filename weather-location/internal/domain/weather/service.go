package weather

import (
	"errors"
	"fmt"
	"github.com/ramonsoterio/weather-location/internal/domain"
	"github.com/ramonsoterio/weather-location/internal/infra/clients/viacep"
	"github.com/ramonsoterio/weather-location/internal/infra/clients/weatherapi"
	errors2 "github.com/ramonsoterio/weather-location/pkg/errors"
	"net/http"
)

type LocationClient interface {
	FetchLocation(zipCode string) (viacep.APIResponse, error)
}

type WeatherClient interface {
	FetchCurrentWeather(cityName string) (weatherapi.Response, error)
}
type Service struct {
	weatherClient  WeatherClient
	locationClient LocationClient
}

func NewService(we WeatherClient, loc LocationClient) Service {
	return Service{
		weatherClient:  we,
		locationClient: loc,
	}
}

func (s *Service) FetchWeather(cep string) (out TemperatureOutput, err error) {
	location, err := s.locationClient.FetchLocation(cep)
	if err != nil {
		if errors.Is(err, viacep.ErrCepNotFound) {
			err = domain.ErrLocationNotFound
			return
		}
		err = errors2.New(err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("location", location.Localidade)
	currentWeather, err := s.weatherClient.FetchCurrentWeather(location.Localidade)
	if err != nil {
		return
	}
	out = mapTemperatureLocation(currentWeather)
	return
}

func mapTemperatureLocation(weather weatherapi.Response) TemperatureOutput {
	return NewTemperatureOutput(weather.Current.TempC).
		WithCity(weather.Location.Name)
}
