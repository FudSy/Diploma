package dto

import (
	"time"

	"github.com/google/uuid"
)

// CalendarBooking — booking enriched with resource info for ICS / Google Calendar export.
type CalendarBooking struct {
	BookingID    uuid.UUID `json:"booking_id"`
	UserID       uuid.UUID `json:"user_id"`
	ResourceID   uuid.UUID `json:"resource_id"`
	ResourceName string    `json:"resource_name"`
	ResourceType string    `json:"resource_type"`
	Location     string    `json:"location"`
	Description  string    `json:"description"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Status       string    `json:"status"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CalendarFeedResponse returned to authenticated user — provides URLs to subscribe to.
type CalendarFeedResponse struct {
	Token     string `json:"token"`
	FeedURL   string `json:"feed_url"`   // https:// URL — direct download / subscribe
	WebcalURL string `json:"webcal_url"` // webcal:// URL — opens in user's default calendar app
}

// GoogleLinkResponse — "Add to Google Calendar" deep link for a single booking.
type GoogleLinkResponse struct {
	URL string `json:"url"`
}
