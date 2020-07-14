package schema

import (
	"time"
	"github.com/chunganhbk/gin-go/pkg/util"
)



// User
type User struct {
	ID        string    `json:"id"`
	FullName  string    `json:"full_name" binding:"required"`
	Password  string    `json:"password"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Status    int       `json:"status" binding:"required,max=2,min=1"`
	Creator   string    `json:"creator"`
	CreatedAt time.Time `json:"created_at"`
	UserRoles UserRoles `json:"user_roles" binding:"required,gt=0"`
}

func (a *User) String() string {
	return util.JSONMarshalToString(a)
}

// CleanSecure
func (a *User) CleanSecure() *User {
	a.Password = ""
	return a
}

// UserQueryParam
type UserQueryParam struct {
	PaginationParam
	Email   string   `form:"email"`
	QueryValue string   `form:"queryValue"`
	Status     int      `form:"status"`
	RoleIDs    []string `form:"-"`
}

// User Query Options
type UserQueryOptions struct {
	OrderFields []*OrderField
}

// User Query Result
type UserQueryResult struct {
	Data       Users
	PageResult *PaginationResult
}

// ToShow Result
func (a UserQueryResult) ToShowResult(mUserRoles map[string]UserRoles, mRoles map[string]*Role) *UserShowQueryResult {
	return &UserShowQueryResult{
		PageResult: a.PageResult,
		Data:       a.Data.ToUserShows(mUserRoles, mRoles),
	}
}

// Users
type Users []*User

// ToIDs
func (a Users) ToIDs() []string {
	idList := make([]string, len(a))
	for i, item := range a {
		idList[i] = item.ID
	}
	return idList
}

// ToUserShows
func (a Users) ToUserShows(mUserRoles map[string]UserRoles, mRoles map[string]*Role) UserShows {
	list := make(UserShows, len(a))
	for i, item := range a {
		showItem := new(UserShow)
		util.StructMapToStruct(item, showItem)
		for _, roleID := range mUserRoles[item.ID].ToRoleIDs() {
			if v, ok := mRoles[roleID]; ok {
				showItem.Roles = append(showItem.Roles, v)
			}
		}
		list[i] = showItem
	}

	return list
}

// ----------------------------------------UserRole--------------------------------------

// UserRole
type UserRole struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	RoleID string `json:"role_id"`
}

// UserRole Query Param
type UserRoleQueryParam struct {
	PaginationParam
	UserID  string
	UserIDs []string
}

// UserRole Query Options
type UserRoleQueryOptions struct {
	OrderFields []*OrderField
}

// UserRole Query Result
type UserRoleQueryResult struct {
	Data       UserRoles
	PageResult *PaginationResult
}

// UserRoles
type UserRoles []*UserRole

// ToMap
func (a UserRoles) ToMap() map[string]*UserRole {
	m := make(map[string]*UserRole)
	for _, item := range a {
		m[item.RoleID] = item
	}
	return m
}

// ToRoleIDs
func (a UserRoles) ToRoleIDs() []string {
	list := make([]string, len(a))
	for i, item := range a {
		list[i] = item.RoleID
	}
	return list
}

// ToUserIDMap
func (a UserRoles) ToUserIDMap() map[string]UserRoles {
	m := make(map[string]UserRoles)
	for _, item := range a {
		m[item.UserID] = append(m[item.UserID], item)
	}
	return m
}

// ----------------------------------------UserShow--------------------------------------

// UserShow
type UserShow struct {
	ID        string    `json:"id"`
	FullName  string    `json:"real_name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Roles     []*Role   `json:"roles"`
}

// UserShows
type UserShows []*UserShow

// UserShow QueryResult
type UserShowQueryResult struct {
	Data       UserShows
	PageResult *PaginationResult
}
