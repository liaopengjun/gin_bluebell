package mysql

import (
	"gin_bluebell/models"
	"github.com/jmoiron/sqlx"
	"strings"
)

func CreatePost(p *models.Post) (err error)  {
	sqlStr := "insert into post (post_id,title,content,author_id,community_id) values (?,?,?,?,?)"
	_,err = db.Exec(sqlStr,p.ID, p.Title,p.Content,p.AuthorId,p.CommunityID)
	return
}
func GetPostById(pid int64) (post *models.Post,err error) {
	post = new(models.Post)
	sqlStr := `select post_id,content,author_id,community_id,create_time from post where post_id=?`
	err = db.Get(post,sqlStr,pid)
	return
}
func GetPostList(page,size int64) (posts []*models.Post,err error) {
	sqlStr := `select post_id,content,author_id,community_id,create_time from post order by create_time desc limit ?,?`
	posts = make([]*models.Post,0,2)
	err = db.Select(&posts,sqlStr,(page-1)*size,size)
	return
}

func GetPostListByIDs(ids []string) (postList []*models.Post,err error) {
	sqlStr := `select post_id,content,author_id,title,community_id,create_time from post where post_id in (?) order by FIND_IN_SET (post_id,?)`
	query,args,err := sqlx.In(sqlStr,ids,strings.Join(ids,","))
	if err != nil{
		return nil ,err
	}
	query = db.Rebind(query)
	err = db.Select(&postList,query,args...)
	return
}