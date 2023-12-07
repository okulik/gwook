package model

import "encoding/json"

type Notification struct {
	Object          string         `json:"object"`
	ID              string         `json:"id"`
	Data            map[string]any `json:"data"`
	DataContentType string         `json:"datacontenttype"`
	Project         string         `json:"project"`
	Source          string         `json:"source"`
	SpecVersion     string         `json:"specversion"`
	Time            string         `json:"time"`
	Type            string         `json:"type"`
	Version         string         `json:"version"`
}

func NewNotificationFromJSON(data []byte) (*Notification, error) {
	var notification Notification
	err := json.Unmarshal(data, &notification)
	return &notification, err
}
