package interfaces

type User struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Email    string  `json:"email" gorm:"unique;not null" validate:"required,email"`
	Status   *string `json:"status"`
	Username string  `json:"username" gorm:"unique;not null" validate:"required"`
	Password string  `json:"password" validate:"required"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Name     string  `json:"name" validate:"required"`
	Email    string  `json:"email" validate:"required,email"`
	Status   *string `json:"status"`
	Username string  `json:"username" validate:"required"`
	Password string  `json:"password" validate:"required"`
}
