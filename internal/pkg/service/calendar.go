package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/FudSy/Diploma/internal/dto"
	"github.com/FudSy/Diploma/internal/pkg/repository"
	"github.com/google/uuid"
)

const (
	icsProdID      = "-//Diploma//Resource Booking//EN"
	icsTimeLayout  = "20060102T150405Z"
	icsLineMaxOcts = 75
)

type CalendarService struct {
	bookingRepo repository.Booking
	authRepo    repository.Authorization
}

func NewCalendarService(b repository.Booking, a repository.Authorization) *CalendarService {
	return &CalendarService{bookingRepo: b, authRepo: a}
}

// EnsureToken returns existing calendar token for user or generates a new one.
func (s *CalendarService) EnsureToken(userID uuid.UUID) (string, error) {
	user, err := s.authRepo.GetUserById(userID)
	if err != nil {
		return "", err
	}
	if user.CalendarToken != "" {
		return user.CalendarToken, nil
	}
	token, err := generateToken()
	if err != nil {
		return "", err
	}
	if err := s.authRepo.SetCalendarToken(userID, token); err != nil {
		return "", err
	}
	return token, nil
}

// RotateToken issues a fresh token, invalidating the previous one.
func (s *CalendarService) RotateToken(userID uuid.UUID) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}
	if err := s.authRepo.SetCalendarToken(userID, token); err != nil {
		return "", err
	}
	return token, nil
}

// BookingsByToken resolves a feed token to its owner's bookings.
func (s *CalendarService) BookingsByToken(token string) ([]dto.CalendarBooking, dto.User, error) {
	if token == "" {
		return nil, dto.User{}, errors.New("пустой токен календаря")
	}
	user, err := s.authRepo.GetUserByCalendarToken(token)
	if err != nil {
		return nil, dto.User{}, err
	}
	bookings, err := s.bookingRepo.GetCalendarBookingsByUser(user.ID)
	if err != nil {
		return nil, dto.User{}, err
	}
	return bookings, user, nil
}

func (s *CalendarService) GetBookingForUser(bookingID, userID uuid.UUID, isAdmin bool) (dto.CalendarBooking, error) {
	booking, err := s.bookingRepo.GetCalendarBookingByID(bookingID)
	if err != nil {
		return dto.CalendarBooking{}, err
	}
	if !isAdmin && booking.UserID != userID {
		return dto.CalendarBooking{}, errors.New("доступ запрещён")
	}
	return booking, nil
}

// BuildICS serializes bookings as a single RFC 5545 VCALENDAR document.
func (s *CalendarService) BuildICS(bookings []dto.CalendarBooking, calName string) string {
	var b strings.Builder
	writeICSLine(&b, "BEGIN:VCALENDAR")
	writeICSLine(&b, "VERSION:2.0")
	writeICSLine(&b, "PRODID:"+icsProdID)
	writeICSLine(&b, "CALSCALE:GREGORIAN")
	writeICSLine(&b, "METHOD:PUBLISH")
	if calName != "" {
		writeICSLine(&b, "X-WR-CALNAME:"+icsEscape(calName))
	}

	now := time.Now().UTC().Format(icsTimeLayout)
	for _, bk := range bookings {
		writeICSLine(&b, "BEGIN:VEVENT")
		writeICSLine(&b, "UID:"+bk.BookingID.String()+"@diploma-booking")
		writeICSLine(&b, "DTSTAMP:"+now)
		writeICSLine(&b, "DTSTART:"+bk.StartTime.UTC().Format(icsTimeLayout))
		writeICSLine(&b, "DTEND:"+bk.EndTime.UTC().Format(icsTimeLayout))
		if !bk.UpdatedAt.IsZero() {
			writeICSLine(&b, "LAST-MODIFIED:"+bk.UpdatedAt.UTC().Format(icsTimeLayout))
		}
		writeICSLine(&b, "SUMMARY:"+icsEscape(eventTitle(bk)))
		if bk.Description != "" {
			writeICSLine(&b, "DESCRIPTION:"+icsEscape(bk.Description))
		}
		if bk.Location != "" {
			writeICSLine(&b, "LOCATION:"+icsEscape(bk.Location))
		}
		writeICSLine(&b, "STATUS:"+statusToICS(bk.Status))
		writeICSLine(&b, "TRANSP:OPAQUE")
		writeICSLine(&b, "END:VEVENT")
	}

	writeICSLine(&b, "END:VCALENDAR")
	return b.String()
}

// BuildGoogleLink returns a "Add to Google Calendar" URL for a single booking.
func (s *CalendarService) BuildGoogleLink(bk dto.CalendarBooking) string {
	q := url.Values{}
	q.Set("action", "TEMPLATE")
	q.Set("text", eventTitle(bk))
	q.Set("dates", bk.StartTime.UTC().Format(icsTimeLayout)+"/"+bk.EndTime.UTC().Format(icsTimeLayout))
	if bk.Description != "" {
		q.Set("details", bk.Description)
	}
	if bk.Location != "" {
		q.Set("location", bk.Location)
	}
	return "https://calendar.google.com/calendar/render?" + q.Encode()
}

func eventTitle(bk dto.CalendarBooking) string {
	if bk.ResourceName == "" {
		return "Booking"
	}
	return fmt.Sprintf("Booking: %s", bk.ResourceName)
}

func statusToICS(s string) string {
	switch strings.ToUpper(s) {
	case "CANCELLED":
		return "CANCELLED"
	case "TENTATIVE":
		return "TENTATIVE"
	default:
		return "CONFIRMED"
	}
}

// icsEscape escapes RFC 5545 special chars in TEXT values.
func icsEscape(s string) string {
	r := strings.NewReplacer(
		`\`, `\\`,
		";", `\;`,
		",", `\,`,
		"\r\n", `\n`,
		"\n", `\n`,
		"\r", `\n`,
	)
	return r.Replace(s)
}

// writeICSLine appends an iCalendar content line with CRLF + 75-octet folding.
func writeICSLine(b *strings.Builder, line string) {
	for i := 0; i < len(line); {
		// First line up to 75 octets, continuation lines up to 74 (1 octet for leading space).
		limit := icsLineMaxOcts
		if i > 0 {
			limit = icsLineMaxOcts - 1
			b.WriteByte(' ')
		}
		end := i + limit
		if end > len(line) {
			end = len(line)
		}
		b.WriteString(line[i:end])
		b.WriteString("\r\n")
		i = end
	}
}

func generateToken() (string, error) {
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
