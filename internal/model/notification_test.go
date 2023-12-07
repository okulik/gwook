package model_test

import (
	"regexp"
	"testing"

	"github.com/okulik/gigs-svixer/internal/model"
)

func TestNewNotificationFromJSON(t *testing.T) {
	json := `{
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
	notification, err := model.NewNotificationFromJSON([]byte(json))

	if err != nil {
		t.Errorf("Failed to parse JSON: %v", err)
	}

	if notification.ID != "evt_0TZRAuIV3l4rLP1NlZivWexSK93v" {
		t.Errorf("Unexpected ID value: %v", notification.ID)
	}

	if notification.Type != "subscription.updated" {
		t.Errorf("Unexpected type value: %v", notification.Type)
	}

	if notification.Data["id"] != "sub_0TWhbetD3l4rLP25UfWiO8iouR7B" {
		t.Errorf("Unexpected data['id'] value: %v", notification.Data["id"])
	}
}

func TestNewNotificationFromInvalidJSON(t *testing.T) {
	json := `{
		"object": 123,
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
	_, err := model.NewNotificationFromJSON([]byte(json))

	if !regexp.MustCompile(`json: cannot unmarshal number into Go struct field Notification.object of type string`).Match([]byte(err.Error())) {
		t.Errorf("Expected to return `cannot unmarshal error...` but got: %v", err)
	}
}
