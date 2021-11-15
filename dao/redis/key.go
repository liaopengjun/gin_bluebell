package redis
/*
	Redis Key
*/

const (
	KeyPostInfoHashPrefix = "bluebell:post:" //项目前缀
	KeyPostTimeZSet       = "bluebell:post:time" //zset帖子及发帖时间
	KeyPostScoreZSet      = "bluebell:post:score"//zset帖子及投票分数
	KeyPostVotedZSetPrefix = "bluebell:post:voted:"//zset记录用户及投票类型
)
