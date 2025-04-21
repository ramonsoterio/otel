package weatherapi

import (
	"encoding/json"
	"fmt"
	"github.com/ramonsoterio/weather-location/internal/domain"
	"net/http"
	"net/url"
)

const (
	baseURL                = "http://api.weatherapi.com/v1"
	currentWeatherEndpoint = "/current.json?key=%s&q=%s"
)

type Client struct {
	cli    *http.Client
	apiKey string
}

type Response struct {
	Current  Current  `json:"current"`
	Location Location `json:"location"`
}

type Location struct {
	Name string `json:"name"`
}

type Current struct {
	TempC float64 `json:"temp_C"`
}

func NewClient(client http.Client, apiKey string) *Client {
	return &Client{cli: &client, apiKey: apiKey}
}

func (c *Client) FetchCurrentWeather(cityName string) (resp Response, err error) {
	fullURL := fmt.Sprintf(baseURL+currentWeatherEndpoint, c.apiKey, url.PathEscape(cityName))
	fmt.Println(fullURL)
	res, err := c.cli.Get(fullURL)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return resp, handleError(res.StatusCode)
	}
	if err = json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return
	}
	return resp, err
}

func handleError(statusCode int) error {
	if statusCode == http.StatusBadRequest {
		return domain.ErrLocationNotFound
	}
	return domain.ErrWeatherInformationNotFound
}
