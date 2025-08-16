package model

import (
	"time"

	"github.com/lib/pq"
)

type User struct {
	ID                uint           `gorm:"primaryKey;autoIncrement" db:"id"`
	Email             string         `gorm:"column:email;type:varchar(255);uniqueIndex;not null" db:"email"`
	FirstName         string         `gorm:"column:first_name;type:varchar(255)" db:"first_name"`
	LastName          string         `gorm:"column:last_name;type:varchar(255)" db:"last_name"`
	Role              string         `gorm:"column:role;type:varchar(50);default:'user'" db:"role"`
	ProfileImageURL   string         `gorm:"column:profile_image_url;type:varchar(512)" db:"profile_image_url"`
	PasswordHash      string         `gorm:"column:password_hash;type:varchar(255)" db:"password_hash"`
	GoogleID          string         `gorm:"column:google_id;type:varchar(255)" db:"google_id"`
	CreatedAt         time.Time      `gorm:"column:created_at;type:timestamp;default:now()" db:"created_at"`
	UpdatedAt         time.Time      `gorm:"column:updated_at;type:timestamp;default:now()" db:"updated_at"`
	Age               int            `gorm:"column:age;type:integer" db:"age"`
	Bio               string         `gorm:"column:bio;type:text" db:"bio"`
	Location          string         `gorm:"column:location;type:varchar(255)" db:"location"`
	Interests         pq.StringArray `gorm:"column:interests;type:text[]" db:"interests"`
	Salary            float64        `gorm:"column:salary;type:numeric(12,2)" db:"salary"`
	IsProfileComplete bool           `gorm:"column:is_profile_complete;type:boolean;default:false" db:"is_profile_complete"`
}

// TableName optional (kalau mau pastikan nama tabelnya "users")
func (User) TableName() string { return "users" }
