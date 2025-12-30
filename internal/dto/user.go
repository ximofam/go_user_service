package dto

type UserLoginInput struct {
	Email    string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserResgisterInput struct {
	Email    string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required"`
}

type UserChangePasswordInput struct {
	UserID          uint   `json:"-"`
	CurrentPassword string `json:"cur_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

type UserRefreshTokenInput struct {
	Token string `json:"token" binding:"required"`
}

type UserLogoutInput struct {
	UserID       uint   `json:"-"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}
