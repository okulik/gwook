package svixer

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/okulik/gigs-svixer/internal/model"
	"github.com/okulik/gigs-svixer/internal/settings"
	svix "github.com/svix/svix-webhooks/go"
)

type SvixClient struct {
	Settings *settings.Settings
	Client   *svix.Svix
}

func New(settings *settings.Settings) *SvixClient {
	// retryablehttp provides a familiar HTTP client interface with automatic retries
	// and exponential backoff. Check https://github.com/hashicorp/go-retryablehttp for
	// more details.
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = settings.Svix.RetryMax

	return NewWithHttpClient(settings, retryClient.StandardClient())
}

func NewWithHttpClient(settings *settings.Settings, httpClient *http.Client) *SvixClient {
	options := &svix.SvixOptions{Debug: settings.Svix.Debug}
	if len(settings.Svix.ServerUrl) > 0 {
		url, err := url.Parse(settings.Svix.ServerUrl)
		if err != nil {
			log.Fatalf("unable to parse SVIX_SERVER_URL: %v", err)
		}
		options.ServerUrl = url
	}

	options.HTTPClient = httpClient

	return &SvixClient{
		Settings: settings,
		Client:   svix.New(settings.Svix.AuthToken, options),
	}
}

func (sc *SvixClient) SendNotification(ctx context.Context, notification *model.Notification) (*svix.MessageOut, error) {
	options := &svix.PostOptions{
		IdempotencyKey: &notification.ID,
	}

	svixMessage := &svix.MessageIn{
		EventType: notification.Type,
		Payload:   notification.Data,
	}

	start := time.Now()

	out, err := sc.Client.Message.CreateWithOptions(ctx, sc.Settings.Svix.ApplicationID, svixMessage, options)
	if err != nil {
		return nil, err
	}

	took := time.Since(start)

	log.Printf("message sent to svix, took %v, received %v", took, out)

	return out, nil
}
