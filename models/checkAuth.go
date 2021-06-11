package models

import "github.com/jinzhu/gorm"

type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `gorm:"column:username" json:"username"`
	Password string `gorm:"column:password" json:"password"`
	Role     int    `gorm:"column:role" json:"role"`
	//Avatar   string `gorm:"column:avatar" json:"avatar"`
}

func CheckAuthToken(username, password string) *Auth {
	var auth Auth
	db.Where(Auth{Username: username, Password: password}).First(&auth)
	if auth.ID > 0 {
		return &auth
	}
	return nil
}

func GetUserInfoByName(name string) *Auth {
	var auth Auth
	db.Where("username = ?", name).First(&auth)
	return &auth
}

func CreateUserNumber(username, password string, role int) bool {
	db.Create(&Auth{
		Username: username,
		Password: password,
		Role:     role,
	})
	return true
}

func ExistAuthByName(username string) bool {
	var auth Auth
	err := db.Where("username = ?", username).Find(&auth).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return true
	}
	if auth.Username == username {
		return false
	}
	return true
}
