package server

import (
	"fmt"
	"github.com/ramonsoterio/weather-location/internal/domain/weather"
	"github.com/ramonsoterio/weather-location/internal/handler"
	"github.com/ramonsoterio/weather-location/internal/infra/clients/viacep"
	"github.com/ramonsoterio/weather-location/internal/infra/clients/weatherapi"
	"log"
	"net/http"
	"os"
)

func Run() {
	port := os.Getenv("API_PORT")
	log.Printf("Starting server on port :%v", port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), server())
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func server() *http.ServeMux {
	zipcodeClient := viacep.New(*http.DefaultClient)
	apiKey := os.Getenv("WEATHER_API_KEY")
	weatherClient := weatherapi.NewClient(*http.DefaultClient, apiKey)
	sv := weather.NewService(weatherClient, zipcodeClient)
	h := handler.NewHandler(&sv)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /location/{cep}", h.FetchWeather)
	return mux
}
