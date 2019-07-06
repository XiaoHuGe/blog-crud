package models

import "github.com/jinzhu/gorm"

type Article struct {
	Model

	TagID int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"`

	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Content    string `json:"content"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

// 查询ID是否存在
func ExistArticleByID(id int) (bool, error) {
	var article Article
	err := db.Select("id").Where("id = ?", id).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if article.ID > 0 {
		return true, err
	}
	return false, err
}

// 获取文章总数
func GetArticleTotal(maps interface{}) (count int, err error) {
	err = db.Model(&Article{}).Where(maps).Count(&count).Error
	return
}

// 获取文章列表
func GetArticles(pageNum, pageSize int, maps interface{}) ([]*Article, error) {
	var articles []*Article
	err := db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return articles, err
}

// 获取单个文章
func GetArticle(id int) (*Article, error) {
	var article Article
	err := db.Where("id = ?", id).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	err = db.Model(&article).Related(&article.Tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &article, nil
}

// 修改文章
func EditArticle(id int, data interface{}) error {
	err := db.Model(&Article{}).Where("id = ?", id).Updates(data).Error
	return err
}

// 删除文章
func DeleteArticle(id int) error {
	err := db.Where("id = ?", id).Delete(&Article{}).Error
	return err
}

// 新建文章
func AddArticle(data map[string]interface{}) error {
	err := db.Create(&Article{
		TagID:     data["tag_id"].(int),
		Title:     data["title"].(string),
		Desc:      data["desc"].(string),
		Content:   data["content"].(string),
		CreatedBy: data["created_by"].(string),
		State:     data["state"].(int),
	}).Error
	return err
}

//func (this *Article) BeforeCreate(scope *gorm.Scope) error {
//	scope.SetColumn("CreatedOn", time.Now().Unix())
//	return nil
//}
//
//func (this *Article) BeforeUpdate(scope *gorm.Scope) error {
//	scope.SetColumn("ModifiedOn", time.Now().Unix())
//	return nil
//}
