package dto

type LoginInput struct {
	Email    string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required"`
}

type LoginOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ResgisterInput struct {
	Username string `json:"username" binding:"required,min=3,max=30"`
	Email    string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordInput struct {
	UserID          uint   `json:"-"`
	CurrentPassword string `json:"cur_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

type RefreshTokenInput struct {
	Token string `json:"token" binding:"required"`
}

type LogoutInput struct {
	UserID       uint   `json:"-"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}
