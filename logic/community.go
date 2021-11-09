package logic

import (
	"gin_bluebell/dao/mysql"
	"gin_bluebell/models"
)

func GetCommunityList() ([]*models.Community,error) {
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (*models.CommunityDetail,error) {
	return mysql.GetCommunityDetailById(id)
}
