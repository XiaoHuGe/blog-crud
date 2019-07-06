package models

import (
	"github.com/jinzhu/gorm"
	"xhblog/utils/logging"
)

type Tag struct {
	Model

	Name       string `gorm:"type:varchar(100)" json:"name"`
	CreatedBy  string `gorm:"type:varchar(100)" json:"created_by"`
	ModifiedBy string `gorm:"type:varchar(100)" json:"modified_by"`
	State      int    `gorm:"type:tinyint(3)" json:"state"`
}

func ExistTagByName(name string) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("name = ?", name).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logging.Error(err)
		return false, err
	}
	if tag.ID > 0 {
		return true, err
	}
	return false, nil
}

func ExistTagById(id int) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("id = ?", id).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if tag.ID > 0 {
		return true, err
	}
	return false, nil
}

func GetTagTotal(maps interface{}) (int, error) {
	var count int
	err := db.Model(&Tag{}).Where(maps).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetTags(pageNum int, pageSize int, maps interface{}) ([]*Tag, error) {
	var tags []*Tag
	err := db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return tags, nil
}

func AddTag(name string, status int, createdBy string) error {
	err := db.Create(&Tag{
		Name:      name,
		State:     status,
		CreatedBy: createdBy,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func EditTag(id int, data interface{}) error {
	return db.Model(&Tag{}).Where("id = ?", id).Updates(data).Error

}

func DeleteTag(id int) error {
	return db.Where("id = ?", id).Delete(&Tag{}).Error
}

//// gorm的Callbacks 会自动添加创建时间
//func (this *Tag) BeforeCreate(scope *gorm.Scope) error {
//	scope.SetColumn("CreatedOn", time.Now().Unix())
//	return nil
//}
//
//func (this *Tag) BeforeUpdate(scope *gorm.Scope) error {
//	scope.SetColumn("ModifiedOn", time.Now().Unix())
//	return nil
//}
