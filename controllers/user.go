package controllers

import (
	"errors"
	"fmt"
	"gin_bluebell/dao/mysql"
	"gin_bluebell/logic"
	"gin_bluebell/models"
	"gin_bluebell/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"net/http"
	"strings"
)
// SignUpHandle 注册
func SignUpHandler(c *gin.Context)  {
	//1.处理参数校验
	 p := new(models.ParamSingnUp)
	 if err := c.ShouldBindJSON(p);err != nil{
		 zap.L().Error("singUp with invalid param",zap.Error(err))
		 //判断err是不是validator类型
		 errs,ok := err.(validator.ValidationErrors)
		 if !ok {
			ResponseError(c,CodeInvalidParam) //参数有误
			return
		 }
		 //自定义错误
		 ResponseErrorWithMsg(c,CodeInvalidParam,removeTopStruct(errs.Translate(trans)))
		 return
	 }
	//2.业务处理
	if err := logic.SignUp(p); err !=nil {
		zap.L().Error("singUp with invalid param",zap.Error(err))
		//如果用户存在
		if errors.Is(err,mysql.ErrorUserExit){
			ResponseError(c,CodeUserExist)
			return
		}
		//系统繁忙
		ResponseError(c,CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c,nil)
}
// LoginHandler 登录
func LoginHandler(c *gin.Context){
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p);err != nil{
		zap.L().Error("login with invalid param",zap.Error(err))
		//判断err是不是validator类型
		errs,ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c,CodeInvalidParam)
			return
		}
		//自定义错误
		ResponseErrorWithMsg(c,CodeInvalidParam,removeTopStruct(errs.Translate(trans)))
		return
	}
	//2.业务处理
	aToken, rToken,err := logic.Login(p);
	if  err !=nil {
		zap.L().Error("login with invalid param",zap.String("username",p.Username),zap.String("password",p.Password),zap.Error(err))
		//用户是否存在
		if errors.Is(err,mysql.ErrorUserNotExit){
			ResponseError(c,CodeUserNotExist)
			return
		}
		//用户密码错误
		ResponseError(c,CodeInvalidPassword)
		return
	}

	//3.返回响应
	ResponseSuccess(c, gin.H{
		"username":p.Username,
		"accessToken":  aToken,
		"refreshToken": rToken,
	})
}

//刷新token
func RefreshTokenHandler(c *gin.Context) {
	rt := c.Query("refresh_token")
	// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
	// 这里假设Token放在Header的Authorization中，并使用Bearer开头
	// 这里的具体实现方式要依据你的实际业务情况决定
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		ResponseErrorWithMsg(c, CodeInvalidToken, "请求头缺少Auth Token")
		c.Abort()
		return
	}
	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		ResponseErrorWithMsg(c, CodeInvalidToken, "Token格式不对")
		c.Abort()
		return
	}
	aToken, rToken, err := jwt.RefreshToken(parts[1], rt)
	fmt.Println(err)
	c.JSON(http.StatusOK, gin.H{
		"access_token":  aToken,
		"refresh_token": rToken,
	})
}
