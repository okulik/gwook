package rest_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	gwook "github.com/okulik/gwook/internal"
	"github.com/okulik/gwook/internal/model"
	"github.com/okulik/gwook/internal/rest"
	"github.com/okulik/gwook/internal/settings"
)

var json string = `{
	"object": "object",
	"id": "evt_0TZRAuIV3l4rLP1NlZivWexSK93v",
	"data": {
		"object": "subscription",
		"id": "sub_0TWhbetD3l4rLP25UfWiO8iouR7B"
	},
	"datacontenttype": "application/json",
	"project": "dev",
	"source": "https://api.service.com",
	"specversion": "1.0",
	"time": "2023-03-24T15:50:41Z",
	"type": "subscription.updated",
	"version": "2023-01-30"
}`

func TestNotify(t *testing.T) {
	requestRecorder := doCreateRequest(t, json, http.StatusOK)

	if requestRecorder.Code != http.StatusOK {
		t.Fatalf("unexpected status code: %v", requestRecorder.Code)
	}
}

func TestNotifyWithTooManyRequests(t *testing.T) {
	requestRecorder := doCreateRequest(t, json, http.StatusTooManyRequests)

	if requestRecorder.Code != http.StatusTooManyRequests {
		t.Fatalf("unexpected status code: %v", requestRecorder.Code)
	}
}

func doCreateRequest(t *testing.T, body string, svixCode int) *httptest.ResponseRecorder {
	reader := io.NopCloser(strings.NewReader(body))
	req, err := http.NewRequest("POST", "/notifications", reader)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	settings, err := settings.Load()
	if err != nil {
		t.Fatal("failed to load settings")
	}

	handler := rest.NewNotificationsHandler(settings, NewMockWebhookNotifier(svixCode))

	testRecorder := httptest.NewRecorder()
	router := chi.NewRouter()
	router.Post("/notifications", handler.Notify)
	router.ServeHTTP(testRecorder, req)

	return testRecorder
}

type Error struct {
	status int
	error  string
}

// Error returns non-empty string if there was an error.
func (e Error) Error() string {
	return e.error
}

func (e Error) Status() int {
	return e.status
}

type MockWebhookNotifier struct {
	statusCode int
}

func NewMockWebhookNotifier(statusCode int) gwook.WebhookNotifier {
	return &MockWebhookNotifier{
		statusCode: statusCode,
	}
}

func (mwn *MockWebhookNotifier) SendNotification(ctx context.Context, notification *model.Notification) error {
	return Error{
		status: mwn.statusCode,
	}
}

func (mwn *MockWebhookNotifier) GetStatus(err error) int {
	return err.(Error).Status()
}
