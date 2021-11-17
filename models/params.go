package models

const (
	OrderTime = "time"
	OrderScore = "score"
)

//定义请求的参数结构体
type ParamSingnUp struct {
	Username string `json:"username" binding:"required"`
	Password string  `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	UserID   int64 `json:"user_id"`
}

//投票
type VoteData struct {
	PostID string `json:"post_id" binding:"required"` //帖子id
	Direction int8 `json:"direction,string" binding:"oneof=1 0 -1"` //赞成=1,反对=-1,取消0
}
