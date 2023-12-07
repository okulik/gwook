package rest

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
	svix "github.com/svix/svix-webhooks/go"

	svixer "github.com/okulik/gigs-svixer/internal"
	"github.com/okulik/gigs-svixer/internal/model"
	"github.com/okulik/gigs-svixer/internal/settings"
	"github.com/okulik/gigs-svixer/internal/web"
)

type NotificationsHandler struct {
	SvixClient *svixer.SvixClient
}

func NewNotificationsHandler(settings *settings.Settings) *NotificationsHandler {
	return &NotificationsHandler{
		SvixClient: svixer.New(settings),
	}
}

func NewNotificationsHandlerWithSvixClient(settings *settings.Settings, client *http.Client) *NotificationsHandler {
	return &NotificationsHandler{
		SvixClient: svixer.NewWithHttpClient(settings, client),
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

	out, err := nh.SvixClient.SendNotification(r.Context(), notification)
	if svixError, ok := err.(*svix.Error); ok {
		web.WriteErrorResponse(w, errors.Wrap(err, "failed to send notification"), svixError.Status())
		return
	}

	if err := web.WriteJSONResponse(w, out, http.StatusOK); err != nil {
		web.WriteErrorResponse(w, errors.Wrap(err.Err, "failed to write response"), err.Code)
	}
}
