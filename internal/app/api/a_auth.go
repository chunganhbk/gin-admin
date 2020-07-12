package api

import (

	"github.com/chunganhbk/gin-go/internal/app/services"
	"github.com/chunganhbk/gin-go/internal/app/config"

	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/pkg/app"

	"github.com/chunganhbk/gin-go/pkg/logger"
	"github.com/gin-gonic/gin"

)



// Login
type Login struct {
	LoginBll services.ILogin
}


// Login 用户登录
func (a *Login) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.LoginParam
	if err := app.ParseJSON(c, &item); err != nil {
		app.ResError(c, err)
		return
	}


	user, err := a.LoginBll.Verify(ctx, item.UserName, item.Password)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	userID := user.ID
	// 将用户ID放入上下文
	ginplus.SetUserID(c, userID)

	ctx = logger.NewUserIDContext(ctx, userID)
	tokenInfo, err := a.LoginBll.GenerateToken(ctx, userID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	logger.StartSpan(ctx, logger.SetSpanTitle("用户登录"), logger.SetSpanFuncName("Login")).Infof("登入系统")
	ginplus.ResSuccess(c, tokenInfo)
}

// Logout 用户登出
func (a *Login) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	// 检查用户是否处于登录状态，如果是则执行销毁
	userID := ginplus.GetUserID(c)
	if userID != "" {
		err := a.LoginBll.DestroyToken(ctx, ginplus.GetToken(c))
		if err != nil {
			logger.Errorf(ctx, err.Error())
		}
		logger.StartSpan(ctx, logger.SetSpanTitle("用户登出"), logger.SetSpanFuncName("Logout")).Infof("登出系统")
	}
	ginplus.ResOK(c)
}

// RefreshToken 刷新令牌
func (a *Login) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	tokenInfo, err := a.LoginBll.GenerateToken(ctx, ginplus.GetUserID(c))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, tokenInfo)
}

// GetUserInfo 获取当前用户信息
func (a *Login) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()
	info, err := a.LoginBll.GetLoginInfo(ctx, ginplus.GetUserID(c))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, info)
}

// QueryUserMenuTree 查询当前用户菜单树
func (a *Login) QueryUserMenuTree(c *gin.Context) {
	ctx := c.Request.Context()
	menus, err := a.LoginBll.QueryUserMenuTree(ctx, ginplus.GetUserID(c))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResList(c, menus)
}

// UpdatePassword 更新个人密码
func (a *Login) UpdatePassword(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.UpdatePasswordParam
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	err := a.LoginBll.UpdatePassword(ctx, ginplus.GetUserID(c), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}
