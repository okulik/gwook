package rest

import (
	"io"
	"net/http"

	"github.com/pkg/errors"

	gwook "github.com/okulik/gwook/internal"
	"github.com/okulik/gwook/internal/model"
	"github.com/okulik/gwook/internal/settings"
	"github.com/okulik/gwook/internal/web"
)

type NotificationsHandler struct {
	WebhookImpl gwook.WebhookNotifier
}

func NewNotificationsHandler(settings *settings.Settings, webhook gwook.WebhookNotifier) *NotificationsHandler {
	return &NotificationsHandler{
		WebhookImpl: webhook,
	}
}

func (nh *NotificationsHandler) Notify(w http.ResponseWriter, r *http.Request) {
	buffer, err := io.ReadAll(r.Body)
	if err != nil {
		web.WriteErrorResponse(w, errors.Wrap(err, "failed to read body"), http.StatusBadRequest)
		return
	}

	notification, err := model.NewNotificationFromJSON(buffer)
	if err != nil {
		web.WriteErrorResponse(w, errors.Wrap(err, "failed notification validation"), http.StatusBadRequest)
		return
	}

	if err := nh.WebhookImpl.SendNotification(r.Context(), notification); err != nil {
		web.WriteErrorResponse(w, errors.Wrap(err, "failed to send notification"), nh.WebhookImpl.GetStatus(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}
