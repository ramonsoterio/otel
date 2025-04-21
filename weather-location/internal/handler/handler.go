package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ramonsoterio/weather-location/internal/domain/weather"
	errors2 "github.com/ramonsoterio/weather-location/pkg/errors"
	"log/slog"
	"net/http"
	"regexp"
)

type WeatherService interface {
	FetchWeather(cep string) (weather.TemperatureOutput, error)
}

var (
	errInvalidCep     = errors.New("invalid zipcode")
	errCompilingRegex = errors.New("error compiling regex")
)

type Handler struct {
	weatherService WeatherService
}

func NewHandler(weatherService WeatherService) Handler {
	return Handler{weatherService: weatherService}
}

func (h *Handler) FetchWeather(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("request received: %v\n", r.URL.Path)
	cep := r.PathValue("cep")
	if err := validateCep(cep); err != nil {
		slog.Error(err.Error(), slog.Int("status", http.StatusUnprocessableEntity))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	currentWeather, err := h.weatherService.FetchWeather(cep)
	if err != nil {
		h.handleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(currentWeather)
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	var svErr errors2.Error
	if errors.As(err, &svErr) {
		http.Error(w, svErr.Error(), svErr.StatusCode)
		slog.Error(svErr.Error(), "status", svErr.StatusCode)
		return
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
	slog.Error(err.Error(), "status", http.StatusInternalServerError)
}

func validateCep(cep string) error {
	reg, err := regexp.Compile("^[0-9]{5}-?([0-9]{3}$)")
	if err != nil {
		return errCompilingRegex
	}
	if !reg.Match([]byte(cep)) {
		return errInvalidCep
	}
	return nil
}
