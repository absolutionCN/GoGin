package models

import "github.com/jinzhu/gorm"

type Tag struct {
	Model
	Name       string `gorm:"column:name" json:"tagName"`
	CreateBy   string `gorm:"column:created_by" json:"createdBy"`
	ModifiedBy string `gorm:"column:modified_by" json:"modifiedby"`
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
