package webhook

import "time"

type Subscription struct {
	Id        string    `json:"id"`
	EventType string    `json:"eventType"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
