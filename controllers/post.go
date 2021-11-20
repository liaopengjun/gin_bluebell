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

func PostListHandler(c *gin.Context) {
	page,size := getPageInfo(c)
	data,err := logic.GetPostList(page,size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed",zap.Error(err))
		ResponseError(c,CodeServerBusy)
		return
	}
	ResponseSuccess(c,data)
}

// PostListHandler2 升级帖子列表接口
func PostListHandler2(c *gin.Context) {
	//按照接口参数条件排序 时间或分数
	//初始化帖子列表参数结构体
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	if err := c.ShouldBindQuery(p);err != nil {
		zap.L().Error("PostListHandler2 with failed err ",zap.Error(err))
		ResponseError(c,CodeInvalidParam)
		return
	}
	//获取列表数据
	data,err := logic.GetPostList2(p)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed",zap.Error(err))
		ResponseError(c,CodeServerBusy)
		return
	}
	ResponseSuccess(c,data)
}