package service

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	gwook "github.com/okulik/gwook/internal"
	"github.com/okulik/gwook/internal/rest"
	"github.com/okulik/gwook/internal/settings"
)

const (
	healthPath        string = "/health"
	notificationsPath string = "/notifications"
)

func NewRouter(settings *settings.Settings, webhookImpl gwook.WebhookNotifier) *chi.Mux {
	r := chi.NewRouter()

	r.Use(loggingMiddleware)

	r.Get(healthPath, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	r.Mount(notificationsPath, createNotificationsRoutes(settings, webhookImpl))

	return r
}

func createNotificationsRoutes(settings *settings.Settings, webhookImpl gwook.WebhookNotifier) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.BasicAuth(settings.Auth.Realm, map[string]string{settings.Auth.Username: settings.Auth.Password}))

	notificationsHandler := rest.NewNotificationsHandler(settings, webhookImpl)
	router.Post("/", notificationsHandler.Notify)

	return router
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
