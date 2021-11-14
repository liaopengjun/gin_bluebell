package models

import "time"

//内存对齐 int->string->time->.......
type Post struct {
	ID int64 `json:"id" db:"post_id"`
	Status int32 `json:"status" db:"status"`
	CommunityID int64 `json:"community_id" db:"community_id" binding:"required"`
	AuthorId   uint64  `json:"author_id" db:"author_id" binding:"required"`
	Title string `json:"title" db:"title" binding:"required"`
	Content string `json:"content" db:"content" binding:"required"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
}
//帖子详情返回结构体
type ApiPostDetail struct {
	AuthorName string `json:"author_name"`
	*Post 								//帖子结构体
	*CommunityDetail `json:"community"` //社区信息
}