package models

//构建一个用户注册结构体
type User struct {
	UserID int64 `db:"user_id"`
	Username string `db:"username"`
	Password string `db:"password"`
}
