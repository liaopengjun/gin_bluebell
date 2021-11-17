package controllers

import (
	"gin_bluebell/logic"
	"gin_bluebell/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

//投票
func PostVoteController(c *gin.Context)  {
	//参数校验
	p := new(models.VoteData)
	if err := c.ShouldBindJSON(p); err !=nil {
		errs ,ok := err.(validator.ValidationErrors) //类型断言
		if !ok {
			ResponseError(c,CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans)) //翻译去除掉错误提示的结构体
		ResponseErrorWithMsg(c,CodeInvalidParam,errData)
		return
	}
	//获取请求用户
	userID,err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c,CodeNeedLogin)
	}
	//业务逻辑
	if err := logic.PostForVote(userID,p);err != nil {
		zap.L().Error("logic.PostForVote failed",zap.Error(err))
		ResponseError(c,CodeServerBusy)
		return
	}
	ResponseSuccess(c,nil)
}
