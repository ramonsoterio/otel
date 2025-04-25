package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ramonsoterio/cep-receiver/internal/infra/otelsdk"
	"github.com/ramonsoterio/cep-receiver/internal/infra/weatherlocation"
	"github.com/ramonsoterio/cep-receiver/pkg/cep"
	errors2 "github.com/ramonsoterio/cep-receiver/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var tracer trace.Tracer

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() (err error) {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Set up OpenTelemetry.
	otelShutdown, err := otelsdk.Setup(ctx)
	if err != nil {
		return
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	apiPort := os.Getenv("API_PORT")
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", apiPort),
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      newHTTPHandler(),
	}
	srvErr := make(chan error, 1)
	go func() {
		fmt.Println("Starting HTTP server on port ", apiPort)
		srvErr <- srv.ListenAndServe()
	}()

	// Wait for interruption.
	select {
	case err = <-srvErr:
		// Error when starting HTTP server.
		return
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()
	}

	// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
	err = srv.Shutdown(context.Background())
	return
}

func newHTTPHandler() http.Handler {
	mux := http.NewServeMux()
	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		// Configure the "http.route" for the HTTP instrumentation.
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
		mux.Handle(pattern, handler)
	}
	handleFunc("POST /cep", ProcessCEP)
	// Add HTTP instrumentation for the whole server.
	handler := otelhttp.NewHandler(mux, "/")
	return handler
}

func ProcessCEP(w http.ResponseWriter, r *http.Request) {
	type CepRequest struct {
		Cep string `json:"cep"`
	}
	var req CepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if err := cep.IsValid(req.Cep); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	tracer = otel.Tracer(os.Getenv("OTEL_SERVICE_NAME"))
	_, span := tracer.Start(r.Context(), "request-weather-client")
	defer span.End()
	weatherClient := weatherlocation.NewClient("http://weather-location:8080", http.DefaultClient)
	weather, err := weatherClient.FetchWeather(req.Cep)
	if err != nil {
		var clientErr errors2.Error
		if errors.As(err, &clientErr) {
			if clientErr.StatusCode == http.StatusNotFound {
				http.Error(w, "can not find zipcode", clientErr.StatusCode)
				return
			}
			http.Error(w, clientErr.Error(), clientErr.StatusCode)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Printf("Weather data for CEP %s: %+v\n", req.Cep, weather)
	if err = json.NewEncoder(w).Encode(weather); err != nil {
		http.Error(w, "error encoding to json: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
