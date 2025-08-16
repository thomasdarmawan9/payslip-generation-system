package model

import "time"

type Reimbursement struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	UserID      uint      `gorm:"index;not null"`
	Date        time.Time `gorm:"type:date;not null"`
	Amount      float64   `gorm:"type:numeric(12,2);not null"`
	Description string    `gorm:"type:varchar(255)"`
	CreatedAt   time.Time `gorm:"type:timestamp;default:now()"`
	UpdatedAt   time.Time `gorm:"type:timestamp;default:now()"`
}

func (Reimbursement) TableName() string { return "reimbursements" }
