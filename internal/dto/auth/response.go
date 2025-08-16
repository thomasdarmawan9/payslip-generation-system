package auth

type LoginUserResponse struct {
	Token string       `json:"accessToken"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type RegisterUserResponse struct {
	Data UserResponse `json:"data"`
}
