package rest_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/okulik/gigs-svixer/internal/rest"
	"github.com/okulik/gigs-svixer/internal/settings"
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
	"source": "https://api.gigs.com",
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

	handler := rest.NewNotificationsHandlerWithSvixClient(settings, createSvixTestClient(svixCode))

	testRecorder := httptest.NewRecorder()
	router := chi.NewRouter()
	router.Post("/notifications", handler.Notify)
	router.ServeHTTP(testRecorder, req)

	return testRecorder
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// Creates a mock http client to avoid sending a real request to Svix.
// The client will return a response with the given status code.
func createSvixTestClient(statusCode int) *http.Client {
	return newSvixTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: statusCode,
			Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
			Header:     make(http.Header),
		}
	})
}

func newSvixTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}
