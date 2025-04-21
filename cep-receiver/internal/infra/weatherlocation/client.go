package weatherlocation

import (
	"encoding/json"
	"fmt"
	"github.com/ramonsoterio/cep-receiver/pkg/errors"
	"log/slog"
	"net/http"
)

type Client struct {
	baseURL string
	cli     *http.Client
}

const (
	fetchLocationEndpoint = "/location/%s"
)

func NewClient(baseURL string, cli *http.Client) Client {
	return Client{baseURL, cli}
}

func (w *Client) FetchWeather(cep string) (loc WeatherLocation, err error) {
	endpoint := fmt.Sprintf(w.baseURL+fetchLocationEndpoint, cep)
	resp, err := w.cli.Get(endpoint)
	if err != nil {
		fmt.Printf("error fetching weather data: %s", err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		slog.Error("error fetching weather data", "status_code", resp.StatusCode)
		return loc, errors.New("error fetching weather data", resp.StatusCode)
	}
	err = json.NewDecoder(resp.Body).Decode(&loc)
	fmt.Println(resp.StatusCode)
	if err != nil {
		fmt.Printf("error decoding received weather data: %s", err.Error())
		return
	}
	return
}
