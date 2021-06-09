package models

import (
	"github.com/jinzhu/gorm"
)

type Token struct {
	Model
	ProjectName string `gorm:"column:prodect" json:"prodect"`
	SvcName     string `gorm:"column:servername" json:"servername"`
	Sid         int    `gorm:"column:sid" json:"sid"`
	Token       string `gorm:"column:token" json:"token"`
	State       int    `gorm:"column:state" json:"state"`
	Coverage    int    `gorm:"column:coverage" json:"coverage"`
	ApiTotal    int    `gorm:"column:total" json:"total"`
	Owner       string `gorm:"column:owner" json:"owner"`
}

func GetToken(id int) (token []Token) {
	db.Where("id = ?", id).First(&token)
	return
}

func GetTokens(pageNum int, pageSize int, state int) (tokens []Token) {
	//
	//if pageSize == 0 {
	//	pageSize := 10
	//	db.Where(maps).Limit(pageSize).Find(&tokens)
	//	return
	//}

	db.Where("state = ?", state).Offset(pageNum).Limit(pageSize).Find(&tokens)

	return
}

func GetTokenTotal(maps interface{}) (count int) {
	db.Model(&Token{}).Where(maps).Count(&count)

	return
}

func AddToken(data map[string]interface{}) error {
	token := Token{
		ProjectName: data["prodect"].(string),
		SvcName:     data["servername"].(string),
		Sid:         data["sid"].(int),
		Token:       data["token"].(string),
		Owner:       data["owner"].(string),
	}
	if err := db.Create(&token).Error; err != nil {
		return err
	}
	return nil
}

func ExistToken(token string) bool {
	var getToken Token
	err := db.Where("token = ?", token).Find(&getToken).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false
	}
	if getToken.Token == token {
		return true
	}
	return false

}

func EditTokenId(id int) bool {
	var tokenId Token
	err := db.Select("id").Where("id = ?", id).First(&tokenId).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false
	}
	if tokenId.ID > 0 {
		return true
	} else {
		return false
	}
}

func EditToken(id int, data map[string]interface{}) bool {
	editToken := Token{
		ProjectName: data["prodect"].(string),
		SvcName:     data["servername"].(string),
		Sid:         data["sid"].(int),
		Token:       data["token"].(string),
		Owner:       data["owner"].(string),
	}
	err := db.Model(&editToken).Where("id = ?", id).Updates(&editToken).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false
	} else {
		return true
	}
}

func DeleteToken(id int) bool {
	err := db.Where("id = ?", id).Delete(Token{}).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false
	} else {
		return true
	}
}
