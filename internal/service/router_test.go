package service_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/okulik/gwook/internal/service"
	"github.com/okulik/gwook/internal/settings"
)

func TestRouterWithNotificationsEndpoint(t *testing.T) {
	settings, _ := settings.Load()
	router := service.NewRouter(settings, nil)
	router.Post("/notifications", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest("POST", "/notifications", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("invalid status returned: got %v want %v", status, http.StatusOK)
	}
}

func TestRouterWithNotificationsEndpointUnsupportedAction(t *testing.T) {
	settings, _ := settings.Load()
	router := service.NewRouter(settings, nil)
	router.Post("/notifications", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/notifications", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("invalid status returned: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestRouterWithMissingEndpoint(t *testing.T) {
	settings, _ := settings.Load()
	router := service.NewRouter(settings, nil)
	router.Post("/notifications", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest("POST", "/signals", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("invalid status returned: got %v want %v", status, http.StatusNotFound)
	}
}

func TestRouterWithHealthEndpoint(t *testing.T) {
	settings, _ := settings.Load()
	router := service.NewRouter(settings, nil)
	router.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("invalid status returned: got %v want %v", status, http.StatusOK)
	}
}
