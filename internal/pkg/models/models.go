package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Login        string    `gorm:"unique;not null"`
	Email        string    `gorm:"unique;not null"`
	Name         string    `gorm:"not null;default:''"`
	Surname      string    `gorm:"not null;default:''"`
	PasswordHash string    `gorm:"not null"`
	FullName     string    `gorm:"not null"`
	Role         string    `gorm:"type:varchar(20);default:'USER'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Bookings []Booking `gorm:"foreignKey:UserID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

type Resource struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"not null"`
	Description string
	Type        string `gorm:"type:varchar(50);not null"`
	Capacity    int
	IsActive    bool   `gorm:"default:true"`
	Location    string `gorm:"type:varchar(255)"`
	PhotoURL    string `gorm:"type:varchar(512)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Bookings []Booking `gorm:"foreignKey:ResourceID"`
}

func (r *Resource) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return
}

type Booking struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;not null"`
	ResourceID uuid.UUID `gorm:"type:uuid;not null"`

	StartTime time.Time `gorm:"not null"`
	EndTime   time.Time `gorm:"not null"`
	Status    string    `gorm:"type:varchar(20);default:'CONFIRMED'"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b *Booking) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return
}

type UpdateResourceInput struct {
	Name        *string `gorm:"column:name"`
	Description *string
	Type        *string
	Capacity    *int
	IsActive    *bool
	Location    *string
}

type UpdateBookingInput struct {
	StartTime *time.Time
	EndTime   *time.Time
	Status    *string
}

type ResourceType struct {
	ID      uuid.UUID            `gorm:"type:uuid;primaryKey"`
	Name    string               `gorm:"type:varchar(50);unique;not null"`
	Options []ResourceTypeOption `gorm:"foreignKey:ResourceTypeID;constraint:OnDelete:CASCADE"`
}

func (rt *ResourceType) BeforeCreate(tx *gorm.DB) (err error) {
	if rt.ID == uuid.Nil {
		rt.ID = uuid.New()
	}
	return
}

type ResourceTypeOption struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResourceTypeID uuid.UUID `gorm:"type:uuid;not null;index"`
	Name           string    `gorm:"type:varchar(100);not null"`
	OptionType     string    `gorm:"type:varchar(20);not null;default:'text'"` // text, number, boolean
	IsRequired     bool      `gorm:"default:false"`
}

func (o *ResourceTypeOption) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return
}
