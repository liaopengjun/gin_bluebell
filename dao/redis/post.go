package redis

import (
	"gin_bluebell/models"
	"go.uber.org/zap"
)
func GetPostIDsInOrder(p *models.ParamPostList) ([]string,error) {
	//1.根据用户请求的order参数确定查询redis-key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore{
		key = getRedisKey(KeyPostScoreZSet)
	}
	//2.查询的索引起点
	start := (p.Page -1)* p.Size
	end := start + p.Size -1
	zap.L().Debug("redis range start -> end ",zap.Int64("start",start),zap.Int64("end",end))
	//0-9 9-18 .....
	return client.ZRevRange(key,start,end).Result()
}
