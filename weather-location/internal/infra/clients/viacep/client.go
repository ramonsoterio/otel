package viacep

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	baseURL = "https://viacep.com.br/ws/%s/json"
)

var (
	ErrCepNotFound    = errors.New("not_found_error")
	errResponseDecode = errors.New("response_decode_error")
)

type Client struct {
	cli *http.Client
}

func New(cli http.Client) *Client {
	return &Client{
		cli: &cli,
	}
}

func (c *Client) FetchLocation(cep string) (response APIResponse, err error) {
	resp, err := c.cli.Get(fmt.Sprintf(baseURL, cep))
	if err != nil {
		return response, errors.New("error making viacep request")
	}
	defer resp.Body.Close()
	rawResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		err = errResponseDecode
		return
	}
	return decodeResponse(rawResponse)
}

func decodeResponse(raw []byte) (response APIResponse, err error) {
	var probe map[string]json.RawMessage
	if err = json.Unmarshal(raw, &probe); err != nil {
		err = errResponseDecode
		return
	}
	if _, ok := probe["erro"]; ok {
		err = ErrCepNotFound
		return
	}
	if _, ok := probe["localidade"]; ok {
		if err = json.Unmarshal(raw, &response); err != nil {
			err = errResponseDecode
		}
	}
	return
}
