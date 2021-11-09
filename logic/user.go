package logic

import (
	"gin_bluebell/dao/mysql"
	"gin_bluebell/models"
	"gin_bluebell/pkg/jwt"
	"gin_bluebell/pkg/snowflake"
)

//存放业务逻辑代码

func SignUp(p *models.ParamSingnUp)(err error){
 	//1.参数校验
	if err = mysql.CheckUserExist(p.Username);err != nil {
		return err
	}
	//2.生成UID
	userID, _ :=snowflake.GetID()
	user := models.User{
		UserID: int64(userID),
		Username: p.Username,
		Password: p.Password,
	}
	//3.保存数据
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) (aToken, rToken string,err error) {
	user := &models.User{
		UserID: p.UserID,
		Username: p.Username,
		Password: p.Password,
	}
	if err := mysql.Login(user);err !=nil {
		return "","",err
	}
	//生成jwt token
	return  jwt.GenToken(uint64(user.UserID),user.Username)
}
