package schema

import (
	"strings"
	"time"
	"github.com/chunganhbk/gin-go/pkg/util"
)

// Menu
type Menu struct {
	ID         string      `json:"id"`
	Name       string      `json:"name" binding:"required"`
	Order   int            `json:"order"`
	Icon       string      `json:"icon"`
	Router     string      `json:"router"`
	ParentID   string      `json:"parent_id"`
	ParentPath string      `json:"parent_path"`
	ShowStatus int         `json:"show_status" binding:"required,max=2,min=1"` // status(1:enable 2:disable)
	Status     int         `json:"status" binding:"required,max=2,min=1"`      // status(1:enable 2:disable)
	Memo       string      `json:"memo"`                                       // remarks
	Creator    string      `json:"creator"`                                    // create by
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	Actions    MenuActions `json:"actions"`
}

func (a *Menu) String() string {
	return util.JSONMarshalToString(a)
}

// MenuQueryParam
type MenuQueryParam struct {
	PaginationParam
	IDs              []string `form:"-"`
	Name             string   `form:"-"`
	PrefixParentPath string   `form:"-"`
	QueryValue       string   `form:"queryValue"`
	ParentID         *string  `form:"parentID"`
	ShowStatus       int      `form:"showStatus"`
	Status           int      `form:"status"`
}

// MenuQueryOptions 查询可选参数项
type MenuQueryOptions struct {
	OrderFields []*OrderField // 排序字段
}

// MenuQueryResult 查询结果
type MenuQueryResult struct {
	Data       Menus
	PageResult *PaginationResult
}

// Menus 菜单列表
type Menus []*Menu

func (a Menus) Len() int {
	return len(a)
}

func (a Menus) Less(i, j int) bool {
	return a[i].Order > a[j].Order
}

func (a Menus) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// ToMap
func (a Menus) ToMap() map[string]*Menu {
	m := make(map[string]*Menu)
	for _, item := range a {
		m[item.ID] = item
	}
	return m
}

// SplitParentIDs
func (a Menus) SplitParentIDs() []string {
	idList := make([]string, 0, len(a))
	mIDList := make(map[string]struct{})

	for _, item := range a {
		if _, ok := mIDList[item.ID]; ok || item.ParentPath == "" {
			continue
		}

		for _, pp := range strings.Split(item.ParentPath, "/") {
			if _, ok := mIDList[pp]; ok {
				continue
			}
			idList = append(idList, pp)
			mIDList[pp] = struct{}{}
		}
	}

	return idList
}

// ToTree 转换为菜单树
func (a Menus) ToTree() MenuTrees {
	list := make(MenuTrees, len(a))
	for i, item := range a {
		list[i] = &MenuTree{
			ID:         item.ID,
			Name:       item.Name,
			Icon:       item.Icon,
			Router:     item.Router,
			ParentID:   item.ParentID,
			ParentPath: item.ParentPath,
			Order:   item.Order,
			ShowStatus: item.ShowStatus,
			Status:     item.Status,
			Actions:    item.Actions,
		}
	}
	return list.ToTree()
}

// FillMenuAction 填充菜单动作列表
func (a Menus) FillMenuAction(mActions map[string]MenuActions) Menus {
	for _, item := range a {
		if v, ok := mActions[item.ID]; ok {
			item.Actions = v
		}
	}
	return a
}

// ----------------------------------------MenuTree--------------------------------------

// MenuTree
type MenuTree struct {
	ID         string      `yaml:"-" json:"id"`
	Name       string      `yaml:"name" json:"name"`
	Icon       string      `yaml:"icon" json:"icon"`
	Router     string      `yaml:"router,omitempty" json:"router"`
	ParentID   string      `yaml:"-" json:"parent_id"`
	ParentPath string      `yaml:"-" json:"parent_path"`
	Order      int         `yaml:"order" json:"order"`
	ShowStatus int         `yaml:"-" json:"show_status"`
	Status     int         `yaml:"-" json:"status"`
	Actions    MenuActions `yaml:"actions,omitempty" json:"actions"`
	Children   *MenuTrees  `yaml:"children,omitempty" json:"children,omitempty"`
}

// Menu tree list
type MenuTrees []*MenuTree

// ToTree Convert to tree structure
func (a MenuTrees) ToTree() MenuTrees {
	mi := make(map[string]*MenuTree)
	for _, item := range a {
		mi[item.ID] = item
	}

	var list MenuTrees
	for _, item := range a {
		if item.ParentID == "" {
			list = append(list, item)
			continue
		}
		if pitem, ok := mi[item.ParentID]; ok {
			if pitem.Children == nil {
				children := MenuTrees{item}
				pitem.Children = &children
				continue
			}
			*pitem.Children = append(*pitem.Children, item)
		}
	}
	return list
}

// ----------------------------------------MenuAction--------------------------------------

// MenuAction
type MenuAction struct {
	ID        string              `yaml:"-" json:"id"`
	MenuID    string              `yaml:"-" binding:"required" json:"menu_id"`
	Code      string              `yaml:"code" binding:"required" json:"code"`
	Name      string              `yaml:"name" binding:"required" json:"name"`
	Resources MenuActionResources `yaml:"resources,omitempty" json:"resources"`
}

// MenuActionQueryParam
type MenuActionQueryParam struct {
	PaginationParam
	MenuID string
	IDs    []string
}

// MenuAction Query Options
type MenuActionQueryOptions struct {
	OrderFields []*OrderField
}

// MenuActionQueryResult
type MenuActionQueryResult struct {
	Data       MenuActions
	PageResult *PaginationResult
}

// MenuActions
type MenuActions []*MenuAction

// ToMap map
func (a MenuActions) ToMap() map[string]*MenuAction {
	m := make(map[string]*MenuAction)
	for _, item := range a {
		m[item.Code] = item
	}
	return m
}

// FillResources
func (a MenuActions) FillResources(mResources map[string]MenuActionResources) {
	for i, item := range a {
		a[i].Resources = mResources[item.ID]
	}
}

// ToMenuIDMap
func (a MenuActions) ToMenuIDMap() map[string]MenuActions {
	m := make(map[string]MenuActions)
	for _, item := range a {
		m[item.MenuID] = append(m[item.MenuID], item)
	}
	return m
}

// ----------------------------------------MenuActionResource--------------------------------------

// MenuActionResource
type MenuActionResource struct {
	ID       string `yaml:"-" json:"id"`
	ActionID string `yaml:"-" json:"action_id"`
	Method   string `yaml:"method" binding:"required" json:"method"`
	Path     string `yaml:"path" binding:"required" json:"path"`
}

// Menu Action Resource Query Param
type MenuActionResourceQueryParam struct {
	PaginationParam
	MenuID  string
	MenuIDs []string
}

// Menu Action Resource Query Options
type MenuActionResourceQueryOptions struct {
	OrderFields []*OrderField
}

// MenuAction Resource QueryResult
type MenuActionResourceQueryResult struct {
	Data       MenuActionResources
	PageResult *PaginationResult
}

// MenuAction Resources
type MenuActionResources []*MenuActionResource

// ToMap
func (a MenuActionResources) ToMap() map[string]*MenuActionResource {
	m := make(map[string]*MenuActionResource)
	for _, item := range a {
		m[item.Method+item.Path] = item
	}
	return m
}

// Convert to action ID mapping
func (a MenuActionResources) ToActionIDMap() map[string]MenuActionResources {
	m := make(map[string]MenuActionResources)
	for _, item := range a {
		m[item.ActionID] = append(m[item.ActionID], item)
	}
	return m
}
