package handler

import (
	"avito/pkg/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInitRoutes(t *testing.T) {
	services := &service.Service{}
	handler := NewHandler(services)
	router := handler.InitRoutes()

	testRequest(t, router, "GET", "/api/segments/", http.StatusBadRequest)
	testRequest(t, router, "POST", "/api/segments/", http.StatusBadRequest)
	testRequest(t, router, "PATCH", "/api/segments/", http.StatusBadRequest)
	testRequest(t, router, "DELETE", "/api/segments/", http.StatusBadRequest)
}

func testRequest(t *testing.T, router http.Handler, method, url string, expectedStatusCode int) {
	req, _ := http.NewRequest(method, url, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != expectedStatusCode {
		t.Errorf("expected status code %d for %s request to %s, got %d", expectedStatusCode, method, url, w.Code)
	}
}
