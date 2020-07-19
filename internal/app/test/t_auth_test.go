package test

import (
	"fmt"
	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/pkg/unique"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestAuth(t *testing.T) {
	const registerRoute = apiPrefix + "v1/auth/register"
	var err error

	w := httptest.NewRecorder()

	// register
	addItem := &schema.RegisterUser{
		Email:     fmt.Sprintf("%s@gmail.com", unique.MustUUID().String()),
		FirstName: unique.MustUUID().String(),
		LastName:  unique.MustUUID().String(),
		Password:  "123456",
	}
	engine.ServeHTTP(w, newPostRequest(registerRoute, addItem))
	assert.Equal(t, 200, w.Code)
	var registerRes schema.ResponseResult
	err = parseReader(w.Body, &registerRes)
	assert.Nil(t, err)
	// todo write check struct body response
	//login

}
