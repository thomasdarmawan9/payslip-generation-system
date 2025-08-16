package auth

import "github.com/lib/pq"

type RegisterUserRequest struct {
	Email             string         `json:"email" binding:"required,email"`
	FirstName         string         `json:"first_name" binding:"required"`
	LastName          string         `json:"last_name"`
	ProfileImageURL   string         `json:"profile_image_url"`
	Password          string         `json:"password" binding:"required,min=6"`
	GoogleID          string         `json:"google_id"`
	Age               int            `json:"age"`
	Bio               string         `json:"bio"`
	Location          string         `json:"location"`
	Interests         pq.StringArray `json:"interests" swaggertype:"array,string" example:"coding,reading,travel"`
	IsProfileComplete bool           `json:"is_profile_complete"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
