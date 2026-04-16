package postgres

import (
	"github.com/FudSy/Diploma/internal/dto"
	"github.com/FudSy/Diploma/internal/pkg/models"
	"gorm.io/gorm"
)

type StatsPostgres struct {
	db *gorm.DB
}

func NewStatsPostgres(db *gorm.DB) *StatsPostgres {
	return &StatsPostgres{db}
}

func (r *StatsPostgres) GetOverview() (dto.StatsOverview, error) {
	var stats dto.StatsOverview

	r.db.Model(&models.Booking{}).Count(&stats.TotalBookings)
	r.db.Model(&models.Booking{}).Where("status = ?", "CONFIRMED").Count(&stats.ActiveBookings)
	r.db.Model(&models.Booking{}).Where("status = ?", "CANCELLED").Count(&stats.CancelledBookings)

	r.db.Model(&models.Resource{}).Count(&stats.TotalResources)
	r.db.Model(&models.Resource{}).Where("is_active = ?", true).Count(&stats.ActiveResources)

	r.db.Raw(`
		SELECT r.type, COUNT(b.id) AS count
		FROM bookings b
		JOIN resources r ON b.resource_id = r.id
		WHERE b.status <> 'CANCELLED'
		GROUP BY r.type
		ORDER BY count DESC
	`).Scan(&stats.BookingsByType)

	r.db.Raw(`
		SELECT b.resource_id, r.name AS resource_name, r.type, COUNT(b.id) AS count
		FROM bookings b
		JOIN resources r ON b.resource_id = r.id
		WHERE b.status <> 'CANCELLED'
		GROUP BY b.resource_id, r.name, r.type
		ORDER BY count DESC
		LIMIT 5
	`).Scan(&stats.TopResources)

	r.db.Raw(`
		SELECT TO_CHAR(DATE(start_time), 'YYYY-MM-DD') AS date, COUNT(id) AS count
		FROM bookings
		WHERE start_time >= NOW() - INTERVAL '30 days'
		  AND status <> 'CANCELLED'
		GROUP BY DATE(start_time)
		ORDER BY date
	`).Scan(&stats.BookingsLast30Days)

	r.db.Raw(`
		SELECT EXTRACT(HOUR FROM start_time)::int AS hour, COUNT(id) AS count
		FROM bookings
		WHERE status <> 'CANCELLED'
		GROUP BY EXTRACT(HOUR FROM start_time)
		ORDER BY hour
	`).Scan(&stats.PeakHours)

	if stats.BookingsByType == nil {
		stats.BookingsByType = []dto.TypeStat{}
	}
	if stats.TopResources == nil {
		stats.TopResources = []dto.ResourceStat{}
	}
	if stats.BookingsLast30Days == nil {
		stats.BookingsLast30Days = []dto.DayStat{}
	}
	if stats.PeakHours == nil {
		stats.PeakHours = []dto.HourStat{}
	}

	return stats, nil
}
