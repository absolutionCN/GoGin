package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Article struct {
	Model

	TagID int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"`

	Title      string `gorm:"column:title" json:"title"`
	Desc       string `gorm:"column:desc" json:"desc"`
	Content    string `gorm:"column:content" json:"content"`
	CreatedBy  string `gorm:"column:created_by" json:"created_by"`
	ModifiedBy string `gorm:"column:modified_by" json:"modified_by"`
	State      int    `gorm:"column:state" json:"state"`
}

func ExistArticleByID(id int) bool {
	var article Article

	err := db.Select("id").Where("id = ?", id).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false
	}
	if article.ID > 0 {
		return true
	} else {
		return false
	}
}

func GetArticleTotal(maps interface{}) (count int) {
	db.Model(&Article{}).Where(maps).Count(&count)
	return
}

func GetArticles(pageNum int, pageSize int, maps interface{}) (articles []Article) {
	db.Preloads("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles)
	return
}

func GetArticle(id int) (article Article) {
	db.Where("id = ?", id).First(&article)
	db.Model(&article).Related(&article.Tag)
	return
}

func EditArticle(id int, data interface{}) bool {
	err := db.Model(&Article{}).Where("id = ?", id).Update(data)
	if err != nil {
		return false
	} else {
		return true
	}
}

func AddArticle(data map[string]interface{}) error {
	article := Article{
		TagID:     data["tag_id"].(int),
		Title:     data["title"].(string),
		Desc:      data["desc"].(string),
		Content:   data["content"].(string),
		CreatedBy: data["created_by"].(string),
		State:     data["state"].(int),
	}
	if err := db.Create(&article).Error; err != nil {
		return err
	}
	return nil
}

func DeleteArticle(id int) bool {
	err := db.Model(&Article{}).Where("id = ？", id).Update("delete_on", 1)
	if err != nil {
		return false
	} else {
		return true
	}
}

func (article *Article) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", time.Now().Unix())
	return nil
}

func (article *Article) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("ModifiedOn", time.Now().Unix())
	return nil
}
