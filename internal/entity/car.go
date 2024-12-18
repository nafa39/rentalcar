package entity

import "time"

// Car represents a car available for rent.
type Car struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" db:"name"`
	Category    string    `json:"category" db:"category"`
	PricePerDay float64   `json:"price_per_day" db:"price_per_day"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func (Car) TableName() string {
	return `"rental-car".cars`
}
