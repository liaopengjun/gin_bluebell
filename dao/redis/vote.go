package redis

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"math"
	"time"
)

const (
	oneWeekInSeconds = 7*24*3600 //一星期时间戳
	scorePerVote = 432 //每一票的分数
)

var (
	ErrVoteTimeExpire = errors.New("超出帖子投票时间")
)

/**
1.投票分数:
	每票432分 （86400/200）-> 200张赞成续费帖子一天
2.投票的几种情况
	direction =1(赞成票)
		1.之前没有投过票，现在投赞成 原(0) 差值绝对值 1  +432
		2.之前投反对票，现在投赞成  原(-1) 差值绝对值 2  +432*2
	direction =0(取消投票)
		1.之前投反对票，现在取消投票 原(-1) 差值绝对值 1 +432
		2.之前投赞成票，现在取消投票 原(1) 差值绝对值 1  -432
	direction=-1(反对票)
		1.之前没有投过票，现在投反对票 原(0) 差值绝对值 1 -432
		2.之前投赞成票，现在投反对票  原(1) 差值绝对值 2  -432*2
3.投票限制
	每个帖子投票日起一星期之内能投票，超过一个星期不允许投票
	1.到期之后redis保存赞成票及反对票到mysql中
	2.到期之后删除KeyPostVoredZSetPF
*/

func VoteForPost(userID,postID string,value float64) error {
	//1.判断投票限制
	//取出redis帖子发布时间
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())- postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	//2.更新帖子分数
	//先判断当前帖子是否有投票记录
	ov := client.ZScore(getRedisKey(KeyPostVotedZSetPrefix+postID), userID).Val()
	//判断赋值正负方向
	var op float64
	if value > ov {
		op = 1 //原值小于新值
	} else {
		op = -1 //原值大于新值
	}
	diff := math.Abs(ov - value)
	//事务执行
	pipeline := client.Pipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)//给当前帖子更新分数
	//3.记录用户为该帖子投票数据
	if value == 0 {
		fmt.Println("1111")
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPrefix+postID),userID)
	} else{
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPrefix+postID), redis.Z{
			Score:  value,
			Member: userID,
		})
	}
	_,err := pipeline.Exec()
	return err
}

func CreatePost(postID int64) error {
	//事务执行
	pipeline := client.Pipeline()
	//帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet),redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	//帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet),redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	_,err := pipeline.Exec()
	return err
}
