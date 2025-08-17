package auth

type LoginUserResponse struct {
	Token string       `json:"accessToken"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	Name   string  `json:"name"`
	Email  string  `json:"email"`
	Salary float64 `json:"salary"`
	Role   string  `json:"role"`
}

type RegisterUserResponse struct {
	Data UserResponse `json:"data"`
}
