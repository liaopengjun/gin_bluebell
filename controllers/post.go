package controllers

import (
	"gin_bluebell/logic"
	"gin_bluebell/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func CreatePostHandler(c *gin.Context)  {
	//1.获取参数校验
	p := new(models.Post)
	//从请求中获取用户id
	userID,err := GetCurrentUserID(c)
	if err !=nil{
		ResponseError(c,CodeNeedLogin)
	}
	p.AuthorId = uint64(userID)
	if err := c.ShouldBindJSON(p);err != nil {
		zap.L().Error("CreatePostHandler with invalid param",zap.Error(err))
		ResponseError(c,CodeInvalidParam) //参数有误
		return
	}
	//2.创建帖子
	if err := logic.CreatePost(p);err !=nil{
		zap.L().Error("logic.CreatePost failed",zap.Error(err))
		ResponseError(c,CodeServerBusy)
	}
	//3.返回数据
	ResponseSuccess(c,"创建帖子成功")
}

func CreatePostDetailHandler(c *gin.Context)  {
	//1.获取帖子id
	pidStr := c.Param("id")
	pid,err := strconv.ParseInt(pidStr,10,64)
	if err != nil{
		zap.L().Error("get post detail with invalid param",zap.Error(err))
		ResponseError(c,CodeInvalidParam)
		return
	}
	//2.查询数据
	data,err := logic.GetPostById(pid)
	if err != nil{
		zap.L().Error("logic.GetPostById failed err ",zap.Error(err))
	}
	//3.返回数据
	ResponseSuccess(c,data)
}
