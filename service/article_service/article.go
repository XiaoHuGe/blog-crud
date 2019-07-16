package article_service

import (
	"encoding/json"
	"xhblog/models"
	"xhblog/service/cache_service"
	"xhblog/utils/gredis"
	"xhblog/utils/logging"
)

type Article struct {
	ID       int
	TagID    int
	PageNum  int
	PageSize int
	State    int
	Count    int

	Title      string
	Desc       string
	Content    string
	CreatedBy  string
	ModifiedBy string
}

type AddArticleService struct {
	TagID     int    `json:"tag_id"`
	Title     string `json:"title"`
	Desc      string `json:"desc"`
	Content   string `json:"content"`
	CreatedBy string `json:"created_by"`
	State     int    `json:"state"`
}

func (this *Article) ExistByID() (bool, error) {
	return models.ExistArticleByID(this.ID)

}

func (this *Article) Add() (error) {

	addArticle := models.Article{
		TagID:     this.TagID,
		Title:     this.Title,
		Desc:      this.Desc,
		Content:   this.Content,
		CreatedBy: this.CreatedBy,
		State:     this.State,
	}

	return models.AddArticle(&addArticle)
}

func (this *Article) Get() (*models.Article, error) {
	var cacheArticle *models.Article

	cacheService := cache_service.Article{ID: this.ID}
	key := cacheService.GetArticleKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			err := json.Unmarshal(data, &cacheArticle)
			if err == nil {
				logging.Info("get cache article")
				return cacheArticle, nil
			}
		}
	}

	article, err := models.GetArticle(this.ID)
	if err != nil {
		return nil, err
	}
	logging.Info("get mysql article")
	gredis.Set(key, article, 3600)
	return article, nil
}

var Key string;

func (this *Article) GetAll() ([]*models.Article, error) {
	var cacheArticle []*models.Article

	cacheService := cache_service.Article{
		TagID:    this.TagID,
		PageSize: this.PageSize,
		PageNum:  this.PageNum,
		State:    this.State,
		Count:    this.Count,
	}
	Key = cacheService.GetArticlesKey()
	if gredis.Exists(Key) {
		data, err := gredis.Get(Key)
		if err != nil {
			logging.Info(err)
		} else {
			err := json.Unmarshal(data, &cacheArticle)
			if err == nil {
				logging.Info("get cache articles")
				return cacheArticle, nil
			}
		}
	}

	articles, err := models.GetArticles(this.PageNum, this.PageSize, this.getMaps())
	if err != nil {
		return nil, err
	}
	logging.Info("get mysql articles")
	gredis.Set(Key, articles, 3600)
	return articles, nil
}

// 如何处理在修改数据后，list缓存也更新
func (this *Article) Edit() (error) {
	maps := make(map[string]interface{})
	if this.TagID > 0 {
		maps["tag_id"] = this.TagID
	}
	if this.Title != "" {
		maps["title"] = this.Title
	}
	if this.Desc != "" {
		maps["desc"] = this.Desc
	}
	if this.Content != "" {
		maps["content"] = this.Content
	}
	if this.ModifiedBy != "" {
		maps["modified_by"] = this.ModifiedBy
	}
	err := models.EditArticle(this.ID, maps)
	if err != nil {
		logging.Error(err)
	}

	//删除缓存
	cacheService := cache_service.Article{ID: this.ID}
	key := cacheService.GetArticleKey()
	if gredis.Exists(key) {
		logging.Info("delete cache article")
		gredis.Delete(key)
	}
	if Key != "" {
		logging.Info("delete cache articles")
		gredis.Delete(Key)
	}
	return err
}

func (this *Article) Delete() (error) {
	return models.DeleteArticle(this.ID)
}

func (this *Article) GetCount() (int, error) {
	return models.GetArticleTotal(this.getMaps())
}

func (this *Article) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	//maps["deleted_on"] = 0
	if this.State != -1 {
		maps["state"] = this.State
	}
	if this.TagID != -1 {
		maps["tag_id"] = this.TagID
	}

	return maps
}
