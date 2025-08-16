package model

import "time"

type Attendance struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:"index;not null"`
	Date      time.Time `gorm:"type:date;index:user_date_unique,unique;not null"`
	CreatedAt time.Time `gorm:"type:timestamp;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:now()"`
}

func (Attendance) TableName() string { return "attendances" }
