package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"gin_bluebell/models"
)

const secret = "liaopengjun"

// CheckUserExist 检验用户是否已经注册
func CheckUserExist(username string) (err error){
	sqlStr := `select count(user_id) from user where username = ? `
	var count int
	if err := db.Get(&count,sqlStr,username);err != nil{
		return err
	}
	if count > 0 {
		return ErrorUserExit
	}
	return
}

//InsertUser 保存用户数据库
func InsertUser(u models.User)(err error){
	//加密
	u.Password = encryptPassword(u.Password)
	//执行sql
	sqlStr := "insert into user(user_id,username,password) values (?,?,?)"
	_,err = db.Exec(sqlStr,u.UserID,u.Username,u.Password)
	return err
}
//Login 登录查询
func Login(user *models.User) (err error) {
	oPasssword := user.Password
	sqlStr :=`select user_id,username,password from user where username= ? `
	err = db.Get(user,sqlStr,user.Username)
	if err == sql.ErrNoRows{
		return ErrorUserNotExit
	}
	if err == nil {
		return err
	}
	//判断密码是否一致
	password := encryptPassword(oPasssword)
	if password != user.Password {
		return ErrorPasswordWrong
	}
	return
}

//加密密码
func encryptPassword(opassword string) string  {
	h :=md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(opassword)))
}
