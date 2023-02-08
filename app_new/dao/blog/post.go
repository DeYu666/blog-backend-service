package blog

import (
	"errors"
	"time"

	"github.com/DeYu666/blog-backend-service/app_new/dao/common/dbctl"
	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/DeYu666/blog-backend-service/global"
	"gorm.io/gorm"
)

func GetPosts(options ...func(option *gorm.DB)) ([]*models.BlogPost, error) {
	blogPosts, err := dbctl.GetDBData(&models.BlogPost{}, options...)
	if err != nil {
		return nil, err
	}

	db := global.App.DB

	for index, post := range blogPosts {
		var cate models.BlogCategories
		var tags []models.BlogTag

		err = db.Model(&post).Association("Category").Find(&cate)
		if err != nil {
			return nil, err
		}

		err = db.Model(&post).Association("Tag").Find(&tags)
		if err != nil {
			return nil, err
		}

		// 这个 tag 中的 ID 是 blog_post_tags 中的 ID ，并不是 blog_tag 中的 ID
		for tagIndex, tag := range tags {
			var temp models.BlogTag
			db.Where("name = ?", tag.Name).First(&temp)
			tags[tagIndex] = temp
		}

		blogPosts[index].Category = cate
		blogPosts[index].Tag = tags
		blogPosts[index].Author.Username = "Asa"
		blogPosts[index].Author.ID = 1 // 这里要写不，添加数据的时候，author 的 id 值 不是添加进去了么
	}

	return blogPosts, nil
}

func PostId(id uint) Option {
	return setIdByUint(id)
}

/*
PostPaging 分页逻辑
pageNum 指获取第 pageNum 页的数据
categoryId 指 post 的分类，当 -1 时，获取全部分类
pageSize 指每页中有多少条数据
*/
func PostPaging(pageNum int, categoryId int, pageSize int) Option {
	if pageNum == 0 {
		pageNum = 1
	}
	offset := (pageNum - 1) * pageSize

	if categoryId == -1 {
		return func(db *gorm.DB) {
			db.Limit(pageSize).Offset(offset).Order("created_time desc")
		}
	} else {
		return func(db *gorm.DB) {
			db.Where("`category_id` = ? ", categoryId).Order("created_time desc")
		}
	}
}

func PostByCategoryId(categoryId int) Option {
	if categoryId == -1 {
		return func(db *gorm.DB) {

		}
	}

	return func(db *gorm.DB) {
		db.Where("`category_id` = ?", categoryId)
	}

}

func PostByTagId(tagId uint) Option {
	tags, err := GetTag(TagId(tagId))
	if err != nil || len(tags) == 0 {
		return func(db *gorm.DB) {

		}
	}

	return func(db *gorm.DB) {
		db.Model(tags[0]).Order("created_time desc").Order("blog_post_id desc").Association("Post")
	}
}

func PostTitle(title string) Option {
	return func(db *gorm.DB) {
		db.Where("`title` = ?", title)
	}
}

func PostOrderByDesc() Option {
	return func(db *gorm.DB) {
		db.Order("created_time desc")
	}
}

func AddPost(post *models.BlogPost) error {
	post.AuthorID = 1

	if post.CreatedTime.IsZero() || post.ModifiedTime.IsZero() {
		timeObj := time.Now()
		if post.CreatedTime.IsZero() {
			post.CreatedTime = timeObj
		}
		if post.ModifiedTime.IsZero() {
			post.ModifiedTime = timeObj
		}
	}

	return dbctl.AddDBData(post)
}

func DeletePost(postId uint) error {

	db := global.App.DB

	db = db.Exec("delete from blog_post_tags where blog_post_id = ?", postId)
	if db.Error != nil {
		return db.Error
	}

	db = db.Delete(&models.BlogPost{}, postId)
	if db.Error != nil {
		return db.Error
	}

	return nil
}

func ModifyPost(post *models.BlogPost) error {

	if post.ID.ID == 0 {
		return errors.New("postId not exist, please check id")
	}

	data, err := GetPosts(PostId(post.ID.ID))
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("database is not exist category")
	}

	db := global.App.DB

	result := global.App.DB.Exec("delete from blog_post_tags where blog_post_id = ?", post.ID.ID)
	if result.Error != nil {
		return result.Error
	}

	post.AuthorID = 1
	post.ModifiedTime = time.Now()

	result = db.Save(&post)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// 对博客文章 中设置的密码进行 curd

func GetPostPwd(options ...func(option *gorm.DB)) ([]*models.BlogPostPs, error) {
	return dbctl.GetDBData(&models.BlogPostPs{}, options...)
}

func PostPwdId(id uint) Option {
	return setIdByUint(id)
}

func PostPassWord(pwd string) Option {
	return func(db *gorm.DB) {
		db.Where("`password` = ?", pwd)
	}
}

func AddPostPwd(PostPwd *models.BlogPostPs) error {
	return dbctl.AddDBData(PostPwd)
}

func DeletePostPwd(options ...func(option *gorm.DB)) error {
	return dbctl.DeleteDBData(&models.BlogPostPs{}, options...)
}

func ModifyPostPwd(PostPwd *models.BlogPostPs) error {
	if PostPwd.ID.ID == 0 {
		return errors.New("blog Post id not exist, please check id")
	}

	data, err := GetPostPwd(PostPwdId(PostPwd.ID.ID))
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("database is not exist Blog Post")
	}

	db := global.App.DB
	result := db.Save(&PostPwd)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
