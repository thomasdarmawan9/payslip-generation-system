package model

import "time"

type Overtime struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:"index:user_date_unique,unique;not null"`
	Date      time.Time `gorm:"type:date;index:user_date_unique,unique;not null"`
	Hours     float64   `gorm:"type:numeric(6,2);not null"`
	CreatedAt time.Time `gorm:"type:timestamp;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:now()"`
}

func (Overtime) TableName() string { return "overtimes" }
