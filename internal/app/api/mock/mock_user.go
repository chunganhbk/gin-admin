package mock

import (
	"github.com/gin-gonic/gin"

)


// User Management
type User struct {
}

// Query data
// @Tags User Management
// @Summary Query data
// @Security ApiKeyAuth
// @Param current query int true "Page index" default(1)
// @Param pageSize query int true "Page Size" default(10)
// @Param queryValue query string false "Query value"
// @Param roleIDs query string false "ID (multiple separated by commas)"
// @Param status query int false "status(1:Enable 2:Disable)"
// @Success 200 {array} schema.UserShow "search result：{list:List data,pagination:{current:页索引,pageSize:页大小,total:总数量}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users [get]
func (a *User) Query(c *gin.Context) {
}

// Get 查询指定数据
// Get 查询指定数据
// @Tags 用户管理
// @Summary 查询指定数据
// @Security ApiKeyAuth
// @Param id path string true "唯一标识"
// @Success 200 {object} schema.User
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 404 {object} schema.ErrorResult "{error:{code:0,message:资源不存在}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users/{id} [get]
func (a *User) Get(c *gin.Context) {
}

// Create 创建数据
// @Tags 用户管理
// @Summary 创建数据
// @Security ApiKeyAuth
// @Param body body schema.User true "创建数据"
// @Success 200 {object} schema.IDResult
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users [post]
func (a *User) Create(c *gin.Context) {
}

// Update 更新数据
// @Tags 用户管理
// @Summary 更新数据
// @Security ApiKeyAuth
// @Param id path string true "唯一标识"
// @Param body body schema.User true "更新数据"
// @Success 200 {object} schema.User
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users/{id} [put]
func (a *User) Update(c *gin.Context) {
}

// Delete 删除数据
// @Tags 用户管理
// @Summary 删除数据
// @Security ApiKeyAuth
// @Param id path string true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users/{id} [delete]
func (a *User) Delete(c *gin.Context) {
}

// Enable 启用数据
// @Tags 用户管理
// @Summary 启用数据
// @Security ApiKeyAuth
// @Param id path string true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users/{id}/enable [patch]
func (a *User) Enable(c *gin.Context) {
}

// Disable 禁用数据
// @Tags 用户管理
// @Summary 禁用数据
// @Security ApiKeyAuth
// @Param id path string true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users/{id}/disable [patch]
func (a *User) Disable(c *gin.Context) {
}
// RefreshToken 刷新令牌
// @Tags 登录管理
// @Summary 刷新令牌
// @Security ApiKeyAuth
// @Success 200 {object} schema.LoginTokenInfo
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/pub/refresh-token [post]
func (a *User) RefreshToken(c *gin.Context) {
}

// GetUserInfo 获取当前用户信息
// @Tags 登录管理
// @Summary 获取当前用户信息
// @Security ApiKeyAuth
// @Success 200 {object} schema.UserLoginInfo
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/pub/current/user [get]
func (a *User) GetUserInfo(c *gin.Context) {
}

// QueryUserMenuTree 查询当前用户菜单树
// @Tags 登录管理
// @Summary 查询当前用户菜单树
// @Security ApiKeyAuth
// @Success 200 {object} schema.Menu "查询结果：{list:菜单树}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/pub/current/menutree [get]
func (a *User) QueryUserMenuTree(c *gin.Context) {
}

// UpdatePassword
// @Tags Current User
// @Summary Change personal password
// @Security ApiKeyAuth
// @Param body body schema.UpdatePasswordParam true "Request parameters"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 400 {object} schema.ErrorResult "{code:400, message:Invalid request parameter}"
// @Failure 401 {object} schema.ErrorResult "{code:401, message:unauthorized}}"
// @Failure 500 {object} schema.ErrorResult "{code:500, message:server error}}"
// @Router /api/v1/pub/current/password [put]
func (a *User) UpdatePassword(c *gin.Context) {
}
