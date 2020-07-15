package mock

import "github.com/gin-gonic/gin"

type Auth struct {
}

// @Tags Auth
// @Summary User login
// @Param body body schema.LoginParam true "Request parameters"
// @Success 200 {object} schema.LoginTokenInfo
// @Failure 400 {object} schema.ErrorResult "{code:401, message:Invalid request parameter}"
// @Failure 500 {object} schema.ErrorResult "{code:500, message: Server Error}"
// @Router /api/v1/auth/login [post]
func (a *Auth) Login(c *gin.Context) {
}

// @Tags Auth
// @Summary User register
// @Param body body schema.LoginParam true "Request parameters"
// @Success 200 {object} schema.LoginTokenInfo
// @Failure 400 {object} schema.ErrorResult "{code:401, message:Invalid request parameter}"
// @Failure 500 {object} schema.ErrorResult "{code:500, message: Server Error}"
// @Router /api/v1/auth/register [post]
func (a *Auth) Register(c *gin.Context) {
}


