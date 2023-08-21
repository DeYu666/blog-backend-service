package repository

import (
	"context"
	"fmt"
	"strings"

	blog "github.com/DeYu666/blog-backend-service/lib/log"
	"github.com/DeYu666/blog-backend-service/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Post interface {
	CountPosts(ctx context.Context, tx *gorm.DB, condition FindPostArg) (int64, error)
	FindPosts(ctx context.Context, tx *gorm.DB, condition FindPostArg) ([]model.BlogPost, error)
	GetPost(ctx context.Context, tx *gorm.DB, postId uint) (model.BlogPost, error)
	CreatePost(ctx context.Context, tx *gorm.DB, post model.BlogPost) error
	UpdatePost(ctx context.Context, tx *gorm.DB, post model.BlogPost) error
	DeletePost(ctx context.Context, tx *gorm.DB, postId uint) error

	FindPostPwd(ctx context.Context, tx *gorm.DB) ([]model.BlogPostPs, error)
	GetPostPwd(ctx context.Context, tx *gorm.DB, postPsId uint) (model.BlogPostPs, error)
	CreatePwd(ctx context.Context, tx *gorm.DB, post model.BlogPostPs) error
	UpdatePostPwd(ctx context.Context, tx *gorm.DB, post model.BlogPostPs) error
	DeletePostPwd(ctx context.Context, tx *gorm.DB, postPsId uint) error
}

type FindPostArg struct {
	IDs         []uint
	Titles      []string
	CategoryIDs []uint
	Offset      int32
	Limit       int32
	NoLimit     bool
}

type post struct{}

func NewPost() Post {
	return &post{}
}

func (p *post) CountPosts(ctx context.Context, tx *gorm.DB, condition FindPostArg) (int64, error) {
	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString(" 1 = 1")
	if condition.IDs != nil {
		query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(condition.IDs), ",")))
		for _, id := range condition.IDs {
			args = append(args, id)
		}
	}

	if condition.Titles != nil {
		query.WriteString(fmt.Sprintf(" AND `title` IN (%s)", RepeatWithSep("?", len(condition.Titles), ",")))
		for _, title := range condition.Titles {
			args = append(args, title)
		}
	}

	if condition.CategoryIDs != nil {
		query.WriteString(fmt.Sprintf(" AND `category_id` IN (%s)", RepeatWithSep("?", len(condition.CategoryIDs), ",")))
		for _, categoryID := range condition.CategoryIDs {
			args = append(args, categoryID)
		}
	}

	var count int64
	if err := tx.Model(&model.BlogPost{}).Where(query.String(), args...).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (p *post) FindPosts(ctx context.Context, tx *gorm.DB, condition FindPostArg) ([]model.BlogPost, error) {

	log := blog.Extract(ctx)

	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString("1 = 1")
	if condition.IDs != nil {
		query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(condition.IDs), ",")))
		for _, id := range condition.IDs {
			args = append(args, id)
		}
	}

	if condition.Titles != nil {
		query.WriteString(fmt.Sprintf(" AND `title` IN (%s)", RepeatWithSep("?", len(condition.Titles), ",")))
		for _, title := range condition.Titles {
			args = append(args, title)
		}
	}

	if condition.CategoryIDs != nil {
		query.WriteString(fmt.Sprintf(" AND `category_id` IN (%s)", RepeatWithSep("?", len(condition.CategoryIDs), ",")))
		for _, categoryID := range condition.CategoryIDs {
			args = append(args, categoryID)
		}
	}

	if !condition.NoLimit {
		query.WriteString(" ORDER BY `created_time` DESC")
		query.WriteString(" LIMIT ? OFFSET ?")
		args = append(args, condition.Limit, condition.Offset)
	}

	log.Debug("query: " + query.String())
	log.Debug("args: ", zap.Any("args", args))

	var posts []model.BlogPost
	if err := tx.Model(&model.BlogPost{}).Where(query.String(), args...).Find(&posts).Error; err != nil {
		return nil, err
	}

	for i, post := range posts {
		var cate model.BlogCategories
		var tags []model.BlogTag

		if err := tx.Model(&post).Association("Category").Find(&cate); err != nil {
			return nil, err
		}

		if err := tx.Model(&post).Association("Tag").Find(&tags); err != nil {
			return nil, err
		}

		// 这个 tag 中的 ID 是 blog_post_tags 中的 ID ，并不是 blog_tag 中的 ID
		for i, tag := range tags {
			var temp model.BlogTag
			tx.Where("name = ?", tag.Name).First(&temp)
			tags[i] = temp
		}

		posts[i].Category = cate
		posts[i].Tag = tags
	}

	return posts, nil
}

func (p *post) GetPost(ctx context.Context, tx *gorm.DB, postId uint) (model.BlogPost, error) {
	var post model.BlogPost
	if err := tx.Model(&model.BlogPost{}).Where("`id` = ?", postId).First(&post).Error; err != nil {
		return model.BlogPost{}, err
	}

	var cate model.BlogCategories
	var tags []model.BlogTag

	if err := tx.Model(&post).Association("Category").Find(&cate); err != nil {
		return model.BlogPost{}, err
	}

	if err := tx.Model(&post).Association("Tag").Find(&tags); err != nil {
		return model.BlogPost{}, err
	}

	// 这个 tag 中的 ID 是 blog_post_tags 中的 ID ，并不是 blog_tag 中的 ID
	for i, tag := range tags {
		var temp model.BlogTag
		tx.Where("name = ?", tag.Name).First(&temp)
		tags[i] = temp
	}

	post.Category = cate
	post.Tag = tags

	return post, nil
}

func (p *post) CreatePost(ctx context.Context, tx *gorm.DB, post model.BlogPost) error {
	if err := tx.Model(&model.BlogPost{}).Create(&post).Error; err != nil {
		return err
	}
	return nil
}

func (p *post) UpdatePost(ctx context.Context, tx *gorm.DB, post model.BlogPost) error {
	if err := tx.Save(&post).Error; err != nil {
		return err
	}
	return nil
}

func (p *post) DeletePost(ctx context.Context, tx *gorm.DB, postId uint) error {
	if err := tx.Model(&model.BlogPost{}).Where("`id` = ?", postId).Delete(&model.BlogPost{}).Error; err != nil {
		return err
	}
	return nil
}

func (p *post) FindPostPwd(ctx context.Context, tx *gorm.DB) ([]model.BlogPostPs, error) {
	var posts []model.BlogPostPs
	if err := tx.Model(&model.BlogPostPs{}).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (p *post) GetPostPwd(ctx context.Context, tx *gorm.DB, postPsId uint) (model.BlogPostPs, error) {
	var post model.BlogPostPs
	if err := tx.Model(&model.BlogPostPs{}).Where("`id` = ?", postPsId).First(&post).Error; err != nil {
		return model.BlogPostPs{}, err
	}
	return post, nil
}

func (p *post) CreatePwd(ctx context.Context, tx *gorm.DB, post model.BlogPostPs) error {
	if err := tx.Model(&model.BlogPostPs{}).Create(&post).Error; err != nil {
		return err
	}
	return nil
}

func (p *post) UpdatePostPwd(ctx context.Context, tx *gorm.DB, post model.BlogPostPs) error {
	if err := tx.Model(&model.BlogPostPs{}).Where("`id` = ?", post.ID.ID).Updates(&post).Error; err != nil {
		return err
	}
	return nil
}

func (p *post) DeletePostPwd(ctx context.Context, tx *gorm.DB, postPsId uint) error {
	if err := tx.Model(&model.BlogPostPs{}).Where("`id` = ?", postPsId).Delete(&model.BlogPostPs{}).Error; err != nil {
		return err
	}
	return nil
}
