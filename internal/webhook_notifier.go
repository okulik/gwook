package gwook

import (
	"context"

	"github.com/okulik/gwook/internal/model"
)

type WebhookNotifier interface {
	SendNotification(ctx context.Context, notification *model.Notification) error
	GetStatus(err error) int
}
