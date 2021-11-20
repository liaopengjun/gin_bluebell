package redis

import (
	"gin_bluebell/models"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

func getIDsFromKey(key string,page,size int64) ([]string,error)  {
	start := (page -1) * size
	end := start + size -1
	//0-9 9-18 .....
	return client.ZRevRange(key,start,end).Result()
}

func GetPostIDsInOrder(p *models.ParamPostList) ([]string,error) {
	//1.根据用户请求的order参数确定查询redis-key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore{
		key = getRedisKey(KeyPostScoreZSet)
	}
	//2.查询的索引起点
	return  getIDsFromKey(key,p.Page,p.Size)
}
// GetPostVoteData 根据ids查询每篇帖子投票赞成数
func GetPostVoteData(ids []string)(data []int64,err error) {
	//获取所有投票的key (使用pipeline完成一次性查询key)
	pipeline := client.Pipeline()
	for _,id := range ids {
		key  := getRedisKey(KeyPostVotedZSetPrefix+id)
		pipeline.ZCount(key,"1","1")
	}
	comders,err := pipeline.Exec()
	if err != nil{
		return nil,err
	}
	data = make([]int64,0,len(comders))
	for _,comder := range comders {
		v := comder.(*redis.IntCmd).Val() //类型转换成int
		data = append(data,v)
	}
	//zap.L().Debug("ids len comders len",zap.Int64("ids", int64(len(comders))),zap.Int64("comder",int64(len(comders))))
	return
}

//GetCommunityPostIDsInOrder 根据社区ids查询帖子ids
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string,error) {

	//1.根据用户请求的order参数确定查询redis-key
	orderkey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderkey = getRedisKey(KeyPostScoreZSet)
	}

	//使用zinterstore 取出分区集合和帖子集合（分数，时间）交集 组成一个集合
	Ckey := getRedisKey(KeyCommunitySetPrefix+strconv.Itoa(int(p.CommunityID))) //社区key
	key := orderkey + strconv.Itoa(int(p.CommunityID))
	if client.Exists(key).Val() < 1 {
		//不存在 就计算 （从社区集合和orderkey取出最大值的并集）
		pipline :=client.Pipeline()
		pipline.ZInterStore(key,redis.ZStore{
			Aggregate: "MAX",
		},Ckey,p.Order)
		pipline.Expire(key,60*time.Second) //设置超时时间
		_,err := pipline.Exec()
		if err != nil{
			return nil,err
		}
	}

	//根据key查询并集集合
	return getIDsFromKey(key,p.Page,p.Size)
}