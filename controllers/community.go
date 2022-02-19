package controllers

import (
	"gin_bluebell/logic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// CommunityHandler 返回社区列表
func CommunityHandler(c *gin.Context) {
	data,err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic GetCommunityList failed",zap.Error(err))
		ResponseError(c,CodeServerBusy)
		return
	}
	ResponseSuccess(c,data)
}
// CommunityHandler 返回社区详情
func CommunityDetailHandler(c *gin.Context)  {
	idstr := c.Param("id")
	id,err := strconv.ParseInt(idstr,10,64)
	if err != nil{
		ResponseError(c,CodeInvalidParam)
		return
	}
	data,err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic GetCommunityDetail failed",zap.Error(err))
		ResponseError(c,CodeServerBusy)
		return
	}
	ResponseSuccess(c,data)
}
