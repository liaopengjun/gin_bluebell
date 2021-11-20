package logic

import (
	"gin_bluebell/dao/mysql"
	"gin_bluebell/dao/redis"
	"gin_bluebell/models"
	"gin_bluebell/pkg/snowflake"
	"go.uber.org/zap"
)

//  CreatePost 创建帖子
func CreatePost(p *models.Post)(err error)  {
	//1.生成post_id
	post_id, _ := snowflake.GetID()
	p.ID = int64(post_id)
	//2.保存数据
	err =  mysql.CreatePost(p)
	if err !=nil {
		return err
	}
	err = redis.CreatePost(p.ID)
	return
}

// GetPostById 返回帖子详情
func GetPostById(pid int64)(data *models.ApiPostDetail,err error)  {
	//组合数据
	post,err := mysql.GetPostById(pid)
	//查询作者信息
	user,err := mysql.GetUserById(int64(post.AuthorId))
	if err != nil{
		zap.L().Error("mysql.GetUserById failed",zap.Int64("author_id", int64(post.AuthorId)),zap.Error(err))
		return
	}
	//根据社区id查询详细信息
	community,err := mysql.GetCommunityDetailById(post.CommunityID)
	if err != nil{
		zap.L().Error("mysql.GetCommunityDetailById failed ",zap.Int64("community_id",post.CommunityID),zap.Error(err))
		return
	}
	data = &models.ApiPostDetail{
		AuthorName: user.Username,
		Post:       post,
		CommunityDetail: community,
	}
	return
}

func GetPostList(page,size int64) (data []*models.ApiPostDetail,err error) {
	posts,err := mysql.GetPostList(page,size)
	if err != nil {
		zap.L().Error("mysql.GetPostList failed",zap.Error(err))
		return nil,err
	}
	data = make([]*models.ApiPostDetail,0,len(posts))
	for _,post := range posts {
		//查询作者信息
		user,err := mysql.GetUserById(int64(post.AuthorId))
		if err != nil {
			zap.L().Error("mysql.GetUserById failed",zap.Int64("author_id", int64(post.AuthorId)),zap.Error(err))
			continue
		}
		//根据社区id查询详细信息
		community,err := mysql.GetCommunityDetailById(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailById failed ",zap.Int64("community_id",post.CommunityID),zap.Error(err))
			continue
		}
		PostDetail := &models.ApiPostDetail{
			AuthorName: user.Username,
			Post:       post,
			CommunityDetail: community,
		}
		data = append(data,PostDetail)
	}
	return
}

func GetPostList2(p *models.ParamPostList)  (data []*models.ApiPostDetail,err error) {
	//去redis查询id列表
	ids,err := redis.GetPostIDsInOrder(p)
	if err != nil{
		return
	}
	if len(ids) == 0 {
		return
	}
	//根据id查询数据库帖子信息
	posts ,err := mysql.GetPostListByIDs(ids)

	//zap.L().Debug("ids data",zap.Any("ids",ids))
	//zap.L().Debug("post data",zap.Any("postdata",posts))
	for _,post := range posts {
		//查询作者信息
		user,err := mysql.GetUserById(int64(post.AuthorId))
		if err != nil {
			zap.L().Error("mysql.GetUserById failed",zap.Int64("author_id", int64(post.AuthorId)),zap.Error(err))
			continue
		}
		//根据社区id查询详细信息
		community,err := mysql.GetCommunityDetailById(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailById failed ",zap.Int64("community_id",post.CommunityID),zap.Error(err))
			continue
		}
		PostDetail := &models.ApiPostDetail{
			AuthorName: user.Username,
			Post:       post,
			CommunityDetail: community,
		}
		data = append(data,PostDetail)
	}
	return
}
