package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"users-service/src/router"

	"github.com/go-playground/assert/v2"
)

func TestSendRequestToNotExistingRoute(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	req, err := http.NewRequest("GET", "/not-existing-route", nil)
	assert.Equal(t, err, nil)

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)

	assert.Equal(t, recorder.Code, http.StatusMethodNotAllowed)
}

func TestSendRequestWithInvalidAuthHeader(t *testing.T) {
	router, err := router.CreateRouter()
	assert.Equal(t, err, nil)

	req, err := http.NewRequest("GET", "/users/search", nil)
	req.Header.Add("Authorization", "Brer 0297c9f7-56d7-488a-bbb3-05c6d865f58f")
	assert.Equal(t, err, nil)

	req.Header.Add("content-type", "application/json")
	recorder := httptest.NewRecorder()
	router.Engine.ServeHTTP(recorder, req)

	assert.Equal(t, recorder.Code, http.StatusUnauthorized)
}