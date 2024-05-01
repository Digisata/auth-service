package domain

type (
	Profile struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	CreateProfileRequest struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
)
