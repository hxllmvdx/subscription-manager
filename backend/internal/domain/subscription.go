package domain

import "time"

type Subscription struct {
	ID          string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ServiceName string     `gorm:"type:varchar(255);not null" json:"service_name"`
	Price       int        `gorm:"type:integer;not null" json:"price"`
	UserID      string     `gorm:"type:uuid;not null" json:"user_id"`
	StartDate   time.Time  `gorm:"type:date;not null" json:"start_date"`
	EndDate     *time.Time `gorm:"type:date" json:"end_date,omitempty"`
}
