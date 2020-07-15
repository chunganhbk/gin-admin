package impl

import (
	"context"
	"github.com/chunganhbk/gin-go/internal/app/iutil"
	"github.com/chunganhbk/gin-go/internal/app/repositories"
	"github.com/chunganhbk/gin-go/pkg/app"

	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/casbin/casbin/v2"

)



// Role service
type RoleService struct {
	Enforcer      *casbin.SyncedEnforcer
	TransRp    repositories.ITrans
	RoleRp     repositories.IRole
	RoleMenuRp repositories.IRoleMenu
	UserRp     repositories.IUser
}
func NewRoleService (enforcer *casbin.SyncedEnforcer, transRp repositories.ITrans,
	roleRp repositories.IRole, roleMenuRp repositories.IRoleMenu, userRp repositories.IUser) *RoleService{
	return &RoleService{enforcer, transRp,roleRp, roleMenuRp, userRp}
}
// Query role
func (r *RoleService) Query(ctx context.Context, params schema.RoleQueryParam, opts ...schema.RoleQueryOptions) (*schema.RoleQueryResult, error) {
	return r.RoleRp.Query(ctx, params, opts...)
}

// Get list role service
func (r *RoleService) Get(ctx context.Context, id string, opts ...schema.RoleQueryOptions) (*schema.Role, error) {
	item, err := r.RoleRp.Get(ctx, id, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, app.ResponseNotFound()
	}

	roleMenus, err := r.QueryRoleMenus(ctx, id)
	if err != nil {
		return nil, err
	}
	item.RoleMenus = roleMenus

	return item, nil
}

// Query Role Menus
func (r *RoleService) QueryRoleMenus(ctx context.Context, roleID string) (schema.RoleMenus, error) {
	result, err := r.RoleMenuRp.Query(ctx, schema.RoleMenuQueryParam{
		RoleID: roleID,
	})
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}

// Create role
func (r *RoleService) Create(ctx context.Context, item schema.Role) (*schema.IDResult, error) {
	err := r.checkName(ctx, item)
	if err != nil {
		return nil, err
	}

	item.ID = iutil.NewID()
	err = ExecTrans(ctx, r.TransRp, func(ctx context.Context) error {
		for _, rmItem := range item.RoleMenus {
			rmItem.ID = iutil.NewID()
			rmItem.RoleID = item.ID
			err := r.RoleMenuRp.Create(ctx, *rmItem)
			if err != nil {
				return err
			}
		}
		return r.RoleRp.Create(ctx, item)
	})
	if err != nil {
		return nil, err
	}
	LoadCasbinPolicy(ctx, r.Enforcer)
	return schema.NewIDResult(item.ID), nil
}

func (r *RoleService) checkName(ctx context.Context, item schema.Role) error {
	result, err := r.RoleRp.Query(ctx, schema.RoleQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
		Name:            item.Name,
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return app.New400Response(app.ERROR_EXIST_ROLE)
	}
	return nil
}

// Update role
func (r *RoleService) Update(ctx context.Context, id string, item schema.Role) error {
	oldItem, err := r.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return app.ResponseNotFound()
	} else if oldItem.Name != item.Name {
		err := r.checkName(ctx, item)
		if err != nil {
			return err
		}
	}

	item.ID = oldItem.ID
	item.Creator = oldItem.Creator
	item.CreatedAt = oldItem.CreatedAt
	err = ExecTrans(ctx, r.TransRp, func(ctx context.Context) error {
		addRoleMenus, delRoleMenus := r.compareRoleMenus(ctx, oldItem.RoleMenus, item.RoleMenus)
		for _, rmitem := range addRoleMenus {
			rmitem.ID = iutil.NewID()
			rmitem.RoleID = id
			err := r.RoleMenuRp.Create(ctx, *rmitem)
			if err != nil {
				return err
			}
		}

		for _, rmitem := range delRoleMenus {
			err := r.RoleMenuRp.Delete(ctx, rmitem.ID)
			if err != nil {
				return err
			}
		}

		return r.RoleRp.Update(ctx, id, item)
	})
	if err != nil {
		return err
	}
	LoadCasbinPolicy(ctx, r.Enforcer)
	return nil
}

func (r *RoleService) compareRoleMenus(ctx context.Context, oldRoleMenus, newRoleMenus schema.RoleMenus) (addList, delList schema.RoleMenus) {
	mOldRoleMenus := oldRoleMenus.ToMap()
	mNewRoleMenus := newRoleMenus.ToMap()

	for k, item := range mNewRoleMenus {
		if _, ok := mOldRoleMenus[k]; ok {
			delete(mOldRoleMenus, k)
			continue
		}
		addList = append(addList, item)
	}

	for _, item := range mOldRoleMenus {
		delList = append(delList, item)
	}
	return
}

// Delete role
func (r *RoleService) Delete(ctx context.Context, id string) error {
	oldItem, err := r.RoleRp.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return app.ResponseNotFound()
	}

	userResult, err := r.UserRp.Query(ctx, schema.UserQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
		RoleIDs:         []string{id},
	})
	if err != nil {
		return err
	} else if userResult.PageResult.Total > 0 {
		return app.New400Response(app.ERROR_EXIST_ROLE_USER)
	}

	err = ExecTrans(ctx, r.TransRp, func(ctx context.Context) error {
		err := r.RoleMenuRp.DeleteByRoleID(ctx, id)
		if err != nil {
			return err
		}

		return r.RoleRp.Delete(ctx, id)
	})
	if err != nil {
		return err
	}

	LoadCasbinPolicy(ctx, r.Enforcer)
	return nil
}

// Update Status Role
func (r *RoleService) UpdateStatus(ctx context.Context, id string, status int) error {
	oldItem, err := r.RoleRp.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return app.ResponseNotFound()
	}

	err = r.RoleRp.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}
	LoadCasbinPolicy(ctx, r.Enforcer)
	return nil
}
