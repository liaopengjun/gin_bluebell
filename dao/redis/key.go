package redis
/*
	Redis Key
*/

const (
	KeyPostInfoHashPrefix = "bluebell" //项目前缀
	KeyPostTimeZSet       = "post:time" //zset帖子及发帖时间
	KeyPostScoreZSet      = "post:score"//zset帖子及投票分数
	KeyPostVotedZSetPrefix = "post:voted:"//zset记录用户及投票类型
)

//给redis key加上前缀
func getRedisKey(key string) string  {
 	return KeyPostInfoHashPrefix + key
}