package logic

import (
	"gin_bluebell/models"
	"gin_bluebell/dao/redis"
	"go.uber.org/zap"
	"strconv"
)
//基于用户投票

// PostForVote
func PostForVote(userID int64,p *models.VoteData) error {
	zap.L().Debug("PostForVote ",zap.Int64("userID",userID),zap.String("postID",p.PostID),zap.Int8("Direction",p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID,float64(p.Direction))
}
