package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
)

var ErrorUserNotLogin = errors.New("用户未登录")
const CtxUserIDKey ="UserID"

// GetCurrentUser 获取当前登录的用户ID
func GetCurrentUser(c *gin.Context) (userID int64,err error) {
	uid,ok := c.Get(CtxUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID,ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}
