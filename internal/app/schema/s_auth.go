package schema

import "github.com/chunganhbk/gin-go/pkg/util"

// Login Param
type LoginParam struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type RegisterUser struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password" binding:"required`
	Phone     string `json:"phone"`
	Email     string `json:"email" required`
}

func (u RegisterUser) ToMapUser(result RoleQueryResult) *User {
	showItem := new(User)
	_ = util.StructMapToStruct(u, showItem)

	for _, item := range result.Data {
		showItem.UserRoles = append(showItem.UserRoles, &UserRole{RoleID: item.ID})
	}
	return showItem
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
