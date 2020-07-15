package impl

import (
	"context"
	"fmt"
	"github.com/chunganhbk/gin-go/internal/app/repositories"
	"github.com/chunganhbk/gin-go/internal/app/iutil"
	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/pkg/app"
	"github.com/chunganhbk/gin-go/pkg/util"
	"github.com/casbin/casbin/v2"
	"sort"
)

// UserService
type UserService struct {
	Enforcer       *casbin.SyncedEnforcer
	TransRp        repositories.ITrans
	UserRp         repositories.IUser
	UserRoleRp     repositories.IUserRole
	RoleRp         repositories.IRole
	RoleMenuRp     repositories.IRoleMenu
	MenuRp         repositories.IMenu
	MenuActionRp         repositories.IMenuAction
}
func NewUserService (
	enforcer      *casbin.SyncedEnforcer,
	transRp    repositories.ITrans,
	userRp     repositories.IUser,
	userRoleRp repositories.IUserRole,
	roleMenuRp     repositories.IRoleMenu,
	menuRp         repositories.IMenu,
	menuActionRp         repositories.IMenuAction,
	roleRp     repositories.IRole) *UserService {
	return &UserService{
		Enforcer: enforcer,
		TransRp: transRp,
		RoleRp: roleRp,
		UserRp: userRp,
		UserRoleRp: userRoleRp,
		RoleMenuRp: roleMenuRp,
		MenuRp: menuRp,
		MenuActionRp: menuActionRp,
	}
}
// Query
func (u *UserService) Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error) {
	return u.UserRp.Query(ctx, params, opts...)
}

// QueryShow
func (u *UserService) QueryShow(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserShowQueryResult, error) {
	result, err := u.UserRp.Query(ctx, params, opts...)
	if err != nil {
		return nil, err
	} else if result == nil {
		return nil, nil
	}

	userRoleResult, err := u.UserRoleRp.Query(ctx, schema.UserRoleQueryParam{
		UserIDs: result.Data.ToIDs(),
	})
	if err != nil {
		return nil, err
	}

	roleResult, err := u.RoleRp.Query(ctx, schema.RoleQueryParam{
		IDs: userRoleResult.Data.ToRoleIDs(),
	})
	if err != nil {
		return nil, err
	}

	return result.ToShowResult(userRoleResult.Data.ToUserIDMap(), roleResult.Data.ToMap()), nil
}

// Get
func (u *UserService) Get(ctx context.Context, id string, opts ...schema.UserQueryOptions) (*schema.User, error) {
	item, err := u.UserRp.Get(ctx, id, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, app.ResponseNotFound()
	}

	userRoleResult, err := u.UserRoleRp.Query(ctx, schema.UserRoleQueryParam{
		UserID: id,
	})
	if err != nil {
		return nil, err
	}
	item.UserRoles = userRoleResult.Data

	return item, nil
}

// Create
func (u *UserService) Create(ctx context.Context, item schema.User) (*schema.IDResult, error) {
	err := u.checkEmail(ctx, item)
	if err != nil {
		return nil, err
	}

	item.Password = util.BcryptPwd(item.Password)
	item.FullName = fmt.Sprintf("%s %s", item.FirstName, item.LastName);
	item.ID = iutil.NewID()
	err = ExecTrans(ctx, u.TransRp, func(ctx context.Context) error {
		for _, urItem := range item.UserRoles {
			urItem.ID = iutil.NewID()
			urItem.UserID = item.ID
			err := u.UserRoleRp.Create(ctx, *urItem)
			if err != nil {
				return err
			}
		}

		return u.UserRp.Create(ctx, item)
	})
	if err != nil {
		return nil, err
	}

	LoadCasbinPolicy(ctx, u.Enforcer)
	return schema.NewIDResult(item.ID), nil
}

func (u *UserService) checkEmail(ctx context.Context, item schema.User) error {


	result, err := u.UserRp.Query(ctx, schema.UserQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
		Email:        item.Email,
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return app.New400Response(app.ERROR_EXIST_EMAIL, nil)
	}
	return nil
}

// Update user
func (u *UserService) Update(ctx context.Context, id string, item schema.User) error {
	oldItem, err := u.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return app.ResponseNotFound()
	} else if oldItem.Email != item.Email {
		err := u.checkEmail(ctx, item)
		if err != nil {
			return err
		}
	}

	if item.Password != "" {
		item.Password = util.BcryptPwd(item.Password)
	} else {
		item.Password = oldItem.Password
	}

	item.ID = oldItem.ID
	item.Creator = oldItem.Creator
	item.CreatedAt = oldItem.CreatedAt
	err = ExecTrans(ctx, u.TransRp, func(ctx context.Context) error {
		addUserRoles, delUserRoles := u.compareUserRoles(ctx, oldItem.UserRoles, item.UserRoles)
		for _, rmitem := range addUserRoles {
			rmitem.ID = iutil.NewID()
			rmitem.UserID = id
			err := u.UserRoleRp.Create(ctx, *rmitem)
			if err != nil {
				return err
			}
		}

		for _, rmitem := range delUserRoles {
			err := u.UserRoleRp.Delete(ctx, rmitem.ID)
			if err != nil {
				return err
			}
		}

		return u.UserRp.Update(ctx, id, item)
	})
	if err != nil {
		return err
	}

	LoadCasbinPolicy(ctx, u.Enforcer)
	return nil
}

func (u *UserService) compareUserRoles(ctx context.Context, oldUserRoles, newUserRoles schema.UserRoles) (addList, delList schema.UserRoles) {
	mOldUserRoles := oldUserRoles.ToMap()
	mNewUserRoles := newUserRoles.ToMap()

	for k, item := range mNewUserRoles {
		if _, ok := mOldUserRoles[k]; ok {
			delete(mOldUserRoles, k)
			continue
		}
		addList = append(addList, item)
	}

	for _, item := range mOldUserRoles {
		delList = append(delList, item)
	}
	return
}

// Delete user
func (u *UserService) Delete(ctx context.Context, id string) error {
	oldItem, err := u.UserRp.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return app.ResponseNotFound()
	}

	err = ExecTrans(ctx, u.TransRp, func(ctx context.Context) error {
		err := u.UserRoleRp.DeleteByUserID(ctx, id)
		if err != nil {
			return err
		}

		return u.UserRp.Delete(ctx, id)
	})
	if err != nil {
		return err
	}

	LoadCasbinPolicy(ctx, u.Enforcer)
	return nil
}

// Update Status
func (u *UserService) UpdateStatus(ctx context.Context, id string, status int) error {
	oldItem, err := u.UserRp.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return app.ResponseNotFound()
	}
	oldItem.Status = status

	err = u.UserRp.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}

	LoadCasbinPolicy(ctx, u.Enforcer)
	return nil
}
func (u *UserService) checkAndGetUser(ctx context.Context, userID string) (*schema.User, error) {
	user, err := u.Get(ctx, userID)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, app.New400Response(app.ERROR_NOT_EXIST_USER, nil)
	} else if user.Status != 1 {
		return nil, app.New400Response(app.ERROR_USER_DISABLED, nil)
	}
	return user, nil
}

// Get Login Info
func (u *UserService) GetLoginInfo(ctx context.Context, userID string) (*schema.UserLoginInfo, error) {

	user, err := u.checkAndGetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	info := &schema.UserLoginInfo{
		UserID:   user.ID,
		Email: user.Email,
		FullName: user.FullName,
	}

	userRoleResult, err := u.UserRoleRp.Query(ctx, schema.UserRoleQueryParam{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	if roleIDs := userRoleResult.Data.ToRoleIDs(); len(roleIDs) > 0 {
		roleResult, err := u.RoleRp.Query(ctx, schema.RoleQueryParam{
			IDs:    roleIDs,
			Status: 1,
		})
		if err != nil {
			return nil, err
		}
		info.Roles = roleResult.Data
	}

	return info, nil
}

// Query UserMenu Tree
func (u *UserService) QueryUserMenuTree(ctx context.Context, userID string) (schema.MenuTrees, error) {


	userRoleResult, err := u.UserRoleRp.Query(ctx, schema.UserRoleQueryParam{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	} else if len(userRoleResult.Data) == 0 {
		return nil, app.NoPermissionResponse()
	}

	roleMenuResult, err := u.RoleMenuRp.Query(ctx, schema.RoleMenuQueryParam{
		RoleIDs: userRoleResult.Data.ToRoleIDs(),
	})
	if err != nil {
		return nil, err
	} else if len(roleMenuResult.Data) == 0 {
		return nil, app.NoPermissionResponse()
	}

	menuResult, err := u.MenuRp.Query(ctx, schema.MenuQueryParam{
		IDs:    roleMenuResult.Data.ToMenuIDs(),
		Status: 1,
	})
	if err != nil {
		return nil, err
	} else if len(menuResult.Data) == 0 {
		return nil, app.NoPermissionResponse()
	}

	mData := menuResult.Data.ToMap()
	var qIDs []string
	for _, pid := range menuResult.Data.SplitParentIDs() {
		if _, ok := mData[pid]; !ok {
			qIDs = append(qIDs, pid)
		}
	}

	if len(qIDs) > 0 {
		pmenuResult, err := u.MenuRp.Query(ctx, schema.MenuQueryParam{
			IDs: menuResult.Data.SplitParentIDs(),
		})
		if err != nil {
			return nil, err
		}
		menuResult.Data = append(menuResult.Data, pmenuResult.Data...)
	}

	sort.Sort(menuResult.Data)
	menuActionResult, err := u.MenuActionRp.Query(ctx, schema.MenuActionQueryParam{
		IDs: roleMenuResult.Data.ToActionIDs(),
	})
	if err != nil {
		return nil, err
	}
	return menuResult.Data.FillMenuAction(menuActionResult.Data.ToMenuIDMap()).ToTree(), nil
}

// Update Password
func (u *UserService) ChangePassword(ctx context.Context, userID string, params schema.UpdatePasswordParam) error {


	user, err := u.checkAndGetUser(ctx, userID)
	if err != nil {
		return err
	} else if util.ComparePasswords(user.Password, params.OldPassword) {
		return app.New400Response(app.ERROR_INVALID_OLD_PASS, nil)
	}

	params.NewPassword = util.BcryptPwd(params.NewPassword)
	return u.UserRp.UpdatePassword(ctx, userID, params.NewPassword)
}
