package cache_service

import (
	"fmt"
	"strconv"
	"strings"
	"xhblog/utils/e"
)

type Article struct {
	ID    int
	TagID int
	State int

	PageNum  int
	PageSize int
	Count    int
}

func (a *Article) GetArticleKey() string {
	return e.CACHE_ARTICLE + "_" + strconv.Itoa(a.ID)
}

func (a *Article) GetArticlesKey() string {
	keys := []string{
		e.CACHE_ARTICLE,
		"LIST",
	}

	if a.ID > 0 {
		keys = append(keys, strconv.Itoa(a.ID))
	}
	if a.TagID > 0 {
		keys = append(keys, strconv.Itoa(a.TagID))
	}
	if a.State >= 0 {
		keys = append(keys, strconv.Itoa(a.State))
	}
	if a.PageNum > 0 {
		keys = append(keys, strconv.Itoa(a.PageNum))
	}
	if a.PageSize > 0 {
		keys = append(keys, strconv.Itoa(a.PageSize))
	}
	if a.Count > 0 {
		keys = append(keys, strconv.Itoa(a.Count))
	}
	fmt.Println("key:", strings.Join(keys, "_"))
	return strings.Join(keys, "_")
}
