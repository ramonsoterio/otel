package domain

import (
	"github.com/ramonsoterio/weather-location/pkg/errors"
	"net/http"
)

var (
	ErrLocationNotFound           = errors.New("can not find zipcode", http.StatusNotFound)
	ErrWeatherInformationNotFound = errors.New("weather info location not found", http.StatusNotFound)
)
