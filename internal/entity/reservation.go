package entity

import (
	"time"
)

// Reservation represents a booking made by a user for a car.
type Reservation struct {
	ID         int64     `json:"id" gorm:"primaryKey"`
	UserID     int64     `json:"user_id" gorm:"index"`
	CarID      int64     `json:"car_id" gorm:"index"`
	StartDate  time.Time `json:"start_date" db:"start_date"`
	EndDate    time.Time `json:"end_date" db:"end_date"`
	TotalPrice float64   `json:"total_price" db:"total_price"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`

	Car Car `json:"car" gorm:"foreignKey:CarID"`
}

func (Reservation) TableName() string {
	return `"rental-car".reservations`
}
