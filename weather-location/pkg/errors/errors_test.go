package errors_test

import (
	"github.com/ramonsoterio/weather-location/pkg/errors"
	"net/http"
	"testing"
)

func TestNewError(t *testing.T) {
	err := errors.New("msg", http.StatusBadRequest)
	if err.StatusCode != http.StatusBadRequest {
		t.Fail()
	}
	if err.Message != "msg" {
		t.Fail()
	}
}
