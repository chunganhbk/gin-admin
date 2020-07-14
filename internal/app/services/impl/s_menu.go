package impl

import (
	"context"
	"github.com/chunganhbk/gin-go/internal/app/repositories"
	"github.com/chunganhbk/gin-go/pkg/app"
	"os"
	"github.com/chunganhbk/gin-go/internal/app/iutil"
	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/pkg/util"
)



// Menu
type MenuService struct {
	TransRp              repositories.ITrans
	MenuRp               repositories.IMenu
	MenuActionRp         repositories.IMenuAction
	MenuActionResourceRp repositories.IMenuActionResource
}

func NewMenuService(transRp repositories.ITrans, menuRp repositories.IMenu,
	menuActionRp repositories.IMenuAction,
	menuActionResourceRp repositories.IMenuActionResource) *MenuService{
	return &MenuService{transRp, menuRp, menuActionRp, menuActionResourceRp}
}
// InitData
func (m *MenuService) InitData(ctx context.Context, dataFile string) error {
	result, err := m.MenuRp.Query(ctx, schema.MenuQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return nil
	}

	data, err := m.readData(dataFile)
	if err != nil {
		return err
	}

	return m.createMenus(ctx, "", data)
}

func (m *MenuService) readData(name string) (schema.MenuTrees, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data schema.MenuTrees
	d := util.YAMLNewDecoder(file)
	d.SetStrict(true)
	err = d.Decode(&data)
	return data, err
}

func (m *MenuService) createMenus(ctx context.Context, parentID string, list schema.MenuTrees) error {
	return ExecTrans(ctx, m.TransRp, func(ctx context.Context) error {
		for _, item := range list {
			sitem := schema.Menu{
				Name:       item.Name,
				Sequence:   item.Sequence,
				Icon:       item.Icon,
				Router:     item.Router,
				ParentID:   parentID,
				Status:     1,
				ShowStatus: 1,
				Actions:    item.Actions,
			}
			if v := item.ShowStatus; v > 0 {
				sitem.ShowStatus = v
			}

			nsitem, err := m.Create(ctx, sitem)
			if err != nil {
				return err
			}

			if item.Children != nil && len(*item.Children) > 0 {
				err := m.createMenus(ctx, nsitem.ID, *item.Children)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// Query
func (m *MenuService) Query(ctx context.Context, params schema.MenuQueryParam, opts ...schema.MenuQueryOptions) (*schema.MenuQueryResult, error) {
	menuActionResult, err := m.MenuActionRp.Query(ctx, schema.MenuActionQueryParam{})
	if err != nil {
		return nil, err
	}

	result, err := m.MenuRp.Query(ctx, params, opts...)
	if err != nil {
		return nil, err
	}
	result.Data.FillMenuAction(menuActionResult.Data.ToMenuIDMap())
	return result, nil
}

// Get
func (m *MenuService) Get(ctx context.Context, id string, opts ...schema.MenuQueryOptions) (*schema.Menu, error) {
	item, err := m.MenuRp.Get(ctx, id, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, app.ResponseNotFound()
	}

	actions, err := m.QueryActions(ctx, id)
	if err != nil {
		return nil, err
	}
	item.Actions = actions

	return item, nil
}

// Query menu Actions
func (m *MenuService) QueryActions(ctx context.Context, id string) (schema.MenuActions, error) {
	result, err := m.MenuActionRp.Query(ctx, schema.MenuActionQueryParam{
		MenuID: id,
	})
	if err != nil {
		return nil, err
	} else if len(result.Data) == 0 {
		return nil, nil
	}

	resourceResult, err := m.MenuActionResourceRp.Query(ctx, schema.MenuActionResourceQueryParam{
		MenuID: id,
	})
	if err != nil {
		return nil, err
	}

	result.Data.FillResources(resourceResult.Data.ToActionIDMap())

	return result.Data, nil
}

func (m *MenuService) checkName(ctx context.Context, item schema.Menu) error {
	result, err := m.MenuRp.Query(ctx, schema.MenuQueryParam{
		PaginationParam: schema.PaginationParam{
			OnlyCount: true,
		},
		ParentID: &item.ParentID,
		Name:     item.Name,
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return app.New400Response(app.ERROR_EXIST_MENU_NAME)
	}
	return nil
}

// Create menu
func (m *MenuService) Create(ctx context.Context, item schema.Menu) (*schema.IDResult, error) {
	if err := m.checkName(ctx, item); err != nil {
		return nil, err
	}

	parentPath, err := m.getParentPath(ctx, item.ParentID)
	if err != nil {
		return nil, err
	}
	item.ParentPath = parentPath
	item.ID = iutil.NewID()

	err = ExecTrans(ctx, m.TransRp, func(ctx context.Context) error {
		err := m.createActions(ctx, item.ID, item.Actions)
		if err != nil {
			return err
		}

		return m.MenuRp.Create(ctx, item)
	})
	if err != nil {
		return nil, err
	}

	return schema.NewIDResult(item.ID), nil
}

// create menu action
func (m *MenuService) createActions(ctx context.Context, menuID string, items schema.MenuActions) error {
	for _, item := range items {
		item.ID = iutil.NewID()
		item.MenuID = menuID
		err := m.MenuActionRp.Create(ctx, *item)
		if err != nil {
			return err
		}

		for _, ritem := range item.Resources {
			ritem.ID = iutil.NewID()
			ritem.ActionID = item.ID
			err := m.MenuActionResourceRp.Create(ctx, *ritem)
			if err != nil {
				return err
			}
		}

	}
	return nil
}


func (m *MenuService) getParentPath(ctx context.Context, parentID string) (string, error) {
	if parentID == "" {
		return "", nil
	}

	pitem, err := m.MenuRp.Get(ctx, parentID)
	if err != nil {
		return "", err
	} else if pitem == nil {
		return "", app.New400Response(app.ERROR_INVALID_PARENT)
	}

	return m.joinParentPath(pitem.ParentPath, pitem.ID), nil
}

func (m *MenuService) joinParentPath(parent, id string) string {
	if parent != "" {
		return parent + "/" + id
	}
	return id
}

// Update menu
func (m *MenuService) Update(ctx context.Context, id string, item schema.Menu) error {
	if id == item.ParentID {
		return app.New400Response(app.ERROR_INVALID_PARENT)
	}

	oldItem, err := m.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return app.ResponseNotFound()
	} else if oldItem.Name != item.Name {
		if err := m.checkName(ctx, item); err != nil {
			return err
		}
	}

	item.ID = oldItem.ID
	item.Creator = oldItem.Creator
	item.CreatedAt = oldItem.CreatedAt

	if oldItem.ParentID != item.ParentID {
		parentPath, err := m.getParentPath(ctx, item.ParentID)
		if err != nil {
			return err
		}
		item.ParentPath = parentPath
	} else {
		item.ParentPath = oldItem.ParentPath
	}

	return ExecTrans(ctx, m.TransRp, func(ctx context.Context) error {
		err := m.updateActions(ctx, id, oldItem.Actions, item.Actions)
		if err != nil {
			return err
		}

		err = m.updateChildParentPath(ctx, *oldItem, item)
		if err != nil {
			return err
		}

		return m.MenuModel.Update(ctx, id, item)
	})
}

// update menu action
func (m *MenuService) updateActions(ctx context.Context, menuID string, oldItems, newItems schema.MenuActions) error {
	addActions, delActions, updateActions := m.compareActions(ctx, oldItems, newItems)

	err := m.createActions(ctx, menuID, addActions)
	if err != nil {
		return err
	}

	for _, item := range delActions {
		err := m.MenuActionRp.Delete(ctx, item.ID)
		if err != nil {
			return err
		}

		err = m.MenuActionResourceRp.DeleteByActionID(ctx, item.ID)
		if err != nil {
			return err
		}
	}

	mOldItems := oldItems.ToMap()
	for _, item := range updateActions {
		oitem := mOldItems[item.Code]

		if item.Name != oitem.Name {
			oitem.Name = item.Name
			err := m.MenuActionRp.Update(ctx, item.ID, *oitem)
			if err != nil {
				return err
			}
		}


		addResources, delResources := m.compareResources(ctx, oitem.Resources, item.Resources)
		for _, aritem := range addResources {
			aritem.ID = iutil.NewID()
			aritem.ActionID = oitem.ID
			err := m.MenuActionResourceRp.Create(ctx, *aritem)
			if err != nil {
				return err
			}
		}

		for _, ditem := range delResources {
			err := m.MenuActionResourceRp.Delete(ctx, ditem.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}


func (m *MenuService) compareActions(ctx context.Context, oldActions, newActions schema.MenuActions) (addList, delList, updateList schema.MenuActions) {
	mOldActions := oldActions.ToMap()
	mNewActions := newActions.ToMap()

	for k, item := range mNewActions {
		if _, ok := mOldActions[k]; ok {
			updateList = append(updateList, item)
			delete(mOldActions, k)
			continue
		}
		addList = append(addList, item)
	}

	for _, item := range mOldActions {
		delList = append(delList, item)
	}
	return
}


func (m *MenuService) compareResources(ctx context.Context, oldResources, newResources schema.MenuActionResources) (addList, delList schema.MenuActionResources) {
	mOldResources := oldResources.ToMap()
	mNewResources := newResources.ToMap()

	for k, item := range mNewResources {
		if _, ok := mOldResources[k]; ok {
			delete(mOldResources, k)
			continue
		}
		addList = append(addList, item)
	}

	for _, item := range mOldResources {
		delList = append(delList, item)
	}
	return
}


func (m *MenuService) updateChildParentPath(ctx context.Context, oldItem, newItem schema.Menu) error {
	if oldItem.ParentID == newItem.ParentID {
		return nil
	}

	opath := m.joinParentPath(oldItem.ParentPath, oldItem.ID)
	result, err := m.MenuRp.Query(NewNoTrans(ctx), schema.MenuQueryParam{
		PrefixParentPath: opath,
	})
	if err != nil {
		return err
	}

	npath := m.joinParentPath(newItem.ParentPath, newItem.ID)
	for _, menu := range result.Data {
		err = m.MenuRp.UpdateParentPath(ctx, menu.ID, npath+menu.ParentPath[len(opath):])
		if err != nil {
			return err
		}
	}
	return nil
}

// Delete menu
func (m *MenuService) Delete(ctx context.Context, id string) error {
	oldItem, err := m.MenuRp.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return app.ResponseNotFound()
	}

	result, err := m.MenuRp.Query(ctx, schema.MenuQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
		ParentID:        &id,
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return app.New400Response(app.ERROR_ALLOW_DELETE_WITH_CHILD)
	}

	return ExecTrans(ctx, m.TransRp, func(ctx context.Context) error {
		err = m.MenuActionResourceRp.DeleteByMenuID(ctx, id)
		if err != nil {
			return err
		}

		err := m.MenuActionRp.DeleteByMenuID(ctx, id)
		if err != nil {
			return err
		}

		return m.MenuRp.Delete(ctx, id)
	})
}

// Update Status menu
func (m *MenuService) UpdateStatus(ctx context.Context, id string, status int) error {
	oldItem, err := m.MenuRp.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return app.ResponseNotFound()
	}

	return m.MenuRp.UpdateStatus(ctx, id, status)
}
