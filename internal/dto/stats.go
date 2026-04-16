package dto

import "github.com/google/uuid"

type TypeStat struct {
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

type ResourceStat struct {
	ResourceID   uuid.UUID `json:"resource_id"`
	ResourceName string    `json:"resource_name"`
	Type         string    `json:"type"`
	Count        int64     `json:"count"`
}

type DayStat struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type HourStat struct {
	Hour  int   `json:"hour"`
	Count int64 `json:"count"`
}

type StatsOverview struct {
	TotalBookings      int64          `json:"total_bookings"`
	ActiveBookings     int64          `json:"active_bookings"`
	CancelledBookings  int64          `json:"cancelled_bookings"`
	TotalResources     int64          `json:"total_resources"`
	ActiveResources    int64          `json:"active_resources"`
	BookingsByType     []TypeStat     `json:"bookings_by_type"`
	TopResources       []ResourceStat `json:"top_resources"`
	BookingsLast30Days []DayStat      `json:"bookings_last_30_days"`
	PeakHours          []HourStat     `json:"peak_hours"`
}
