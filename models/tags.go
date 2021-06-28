package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Tag struct {
	Model
	Name       string `gorm:"column:name" json:"tagName"`
	CreateBy   string `gorm:"column:created_by" json:"createdBy"`
	ModifiedBy string `gorm:"column:modified_by" json:"modifiedby"`
	DeletedOn  int    `gorm:"column:deleted_on" json:"deleted_on"`
	State      int    `gorm:"column:state" json:"state"`
}

func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag) {
	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)
	{
		return
	}
}

func GetTagTotal(maps interface{}) (count int) {
	db.Model(&Tag{}).Where(maps).Count(&count)
	return
}

func ExistTag(name string) bool {
	var tag Tag
	err := db.Select("id").Where("name = ?", name).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false
	}
	if tag.ID > 0 {
		return true
	} else {
		return false
	}
}

func ExisTagById(id int) bool {
	var tag Tag
	err := db.Select("id").Where("id = ?", id).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false
	}
	if tag.ID > 0 {
		return true
	} else {
		return false
	}
}

func AddTag(data map[string]interface{}) error {
	tag := Tag{
		Name:     data["name"].(string),
		State:    data["state"].(int),
		CreateBy: data["createdBy"].(string),
	}
	if err := db.Create(&tag).Error; err != nil {
		return err
	}
	return nil
}

func EditTag(id int, data interface{}) bool {
	err := db.Model(&Tag{}).Where("id = ?", id).Update(data)
	if err != nil {
		return false
	} else {
		return true

	}
}

func DeleteTag(id int) bool {
	err := db.Model(&Tag{}).Where("id = ?", id).Update("deleted_on", 1)
	if err != nil {
		return false
	} else {
		return true
	}
}

func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", time.Now().Unix())
	return nil
}

func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("ModifiedOn", time.Now().Unix())

	return nil
}
