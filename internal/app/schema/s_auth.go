package schema

// Login Param
type LoginParam struct {
	UserName    string `json:"user_name" binding:"required"`
	Password    string `json:"password" binding:"required"`

}

// User LoginInfo
type UserLoginInfo struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Roles    Roles  `json:"roles"`
}

// Change Password Param
type UpdatePasswordParam struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}



// Login Token Info
type LoginTokenInfo struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresAt   int64  `json:"expires_at"`
}
