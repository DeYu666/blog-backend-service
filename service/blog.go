package service

import (
	"context"
	"time"

	"github.com/DeYu666/blog-backend-service/lib/client"
	blog "github.com/DeYu666/blog-backend-service/lib/log"
	"github.com/DeYu666/blog-backend-service/model"
	"github.com/DeYu666/blog-backend-service/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BlogService interface {
	GetGeneralCategories(ctx context.Context) ([]model.BlogGeneralCategories, error)
	CreateGeneralCate(ctx context.Context, cate model.BlogGeneralCategories) error
	UpdateGeneralCate(ctx context.Context, cate model.BlogGeneralCategories) error
	DeleteGeneralCate(ctx context.Context, generalCateId uint) error

	GetCategories(ctx context.Context) ([]model.BlogCategories, error)
	CreateCate(ctx context.Context, cate model.BlogCategories) error
	UpdateCate(ctx context.Context, cate model.BlogCategories) error
	DeleteCate(ctx context.Context, cateId uint) error

	GetPostLists(ctx context.Context, cateIds []uint, offset, limit int32) (PostArrWithCount, error)
	GetPostByPostID(ctx context.Context, id uint) (model.BlogPost, error)
	CreatePost(ctx context.Context, post model.BlogPost) error
	UpdatePost(ctx context.Context, post model.BlogPost) error
	DeletePost(ctx context.Context, id uint) error
	GetPostPs(ctx context.Context) ([]model.BlogPostPs, error)
	CreatePostPs(ctx context.Context, postPs model.BlogPostPs) error
	UpdatePostPs(ctx context.Context, postPs model.BlogPostPs) error
	DeletePostPs(ctx context.Context, psId uint) error

	GetTags(ctx context.Context) ([]model.BlogTag, error)
	CreateTag(ctx context.Context, tag model.BlogTag) error
	UpdateTag(ctx context.Context, tag model.BlogTag) error
	DeleteTag(ctx context.Context, tagId uint) error

	GetChickenSoups(ctx context.Context) ([]model.ChickenSoup, error)
	CreateChickenSoup(ctx context.Context, chickSoup model.ChickenSoup) error
	UpdateChickenSoup(ctx context.Context, chickSoup model.ChickenSoup) error
	DeleteChickenSoup(ctx context.Context, chickenSoupId uint) error
}

type blogService struct {
	generalCate repository.GeneralCate
	cate        repository.Category
	tag         repository.Tag
	post        repository.Post
	chickenSoup repository.ChickenSoup
}

func NewBlogService() BlogService {
	return &blogService{
		generalCate: repository.NewGeneralCate(),
		cate:        repository.NewCategory(),
		tag:         repository.NewTag(),
		post:        repository.NewPost(),
		chickenSoup: repository.NewChickenSoup(),
	}
}

func (b *blogService) GetGeneralCategories(ctx context.Context) ([]model.BlogGeneralCategories, error) {

	var generalCates []model.BlogGeneralCategories
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		cond := repository.FindGeneralCateArg{
			NoLimit: true,
		}
		generalCates, err = b.generalCate.FindGeneralCate(ctx, tx, cond)

		generalCateIds := make([]uint, 0)

		for _, generalCate := range generalCates {
			generalCateIds = append(generalCateIds, uint(generalCate.ID.ID))
		}

		cateCond := repository.FindCategoriesArg{
			GeneralCateIds: generalCateIds,
			NoLimit:        true,
		}

		cates, err := b.cate.FindCategory(ctx, tx, cateCond)
		if err != nil {
			return err
		}

		for i, generalCate := range generalCates {
			for _, cate := range cates {
				if generalCate.ID.ID == uint(cate.GeneralID) {
					generalCates[i].BlogCategories = append(generalCates[i].BlogCategories, cate)
				}
			}
		}

		return err
	}, nil)

	return generalCates, err
}

func (b *blogService) CreateGeneralCate(ctx context.Context, cate model.BlogGeneralCategories) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = b.generalCate.CreateGeneralCate(ctx, tx, cate)
		return err
	}, nil)

	return err
}

func (b *blogService) UpdateGeneralCate(ctx context.Context, cate model.BlogGeneralCategories) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = b.generalCate.UpdateGeneralCate(ctx, tx, cate)
		return err
	}, nil)

	return err
}

func (b *blogService) DeleteGeneralCate(ctx context.Context, generalCateId uint) error {
	var err error
	log := blog.Extract(ctx)

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {

		cateCond := repository.FindCategoriesArg{
			GeneralCateIds: []uint{generalCateId},
			NoLimit:        true,
		}
		cates, err := b.cate.FindCategory(ctx, tx, cateCond)
		if err != nil || len(cates) > 0 {
			return ErrGenerateHasCate
		}

		log.Debug("cates: ", zap.Any("cates", cates))

		err = b.generalCate.DeleteGeneralCate(ctx, tx, generalCateId)
		return err
	}, nil)

	return err
}

func (b *blogService) GetCategories(ctx context.Context) ([]model.BlogCategories, error) {

	var cates []model.BlogCategories
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		cond := repository.FindCategoriesArg{
			NoLimit: true,
		}
		cates, err = b.cate.FindCategory(ctx, tx, cond)

		generalCateIds := make([]uint, 0)
		for _, cate := range cates {
			generalCateIds = append(generalCateIds, uint(cate.GeneralID))
		}

		generalCond := repository.FindGeneralCateArg{
			GeneralCateIds: generalCateIds,
			NoLimit:        true,
		}

		generalCates, err := b.generalCate.FindGeneralCate(ctx, tx, generalCond)
		if err != nil {
			return err
		}

		for i, cate := range cates {
			for _, generalCate := range generalCates {
				if cate.GeneralID == int(generalCate.ID.ID) {
					cates[i].General = generalCate
				}
			}
		}

		return err
	}, nil)

	return cates, err
}

func (b *blogService) CreateCate(ctx context.Context, cate model.BlogCategories) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = b.cate.CreateCategory(ctx, tx, cate)
		return err
	}, nil)

	return err
}

func (b *blogService) UpdateCate(ctx context.Context, cate model.BlogCategories) error {
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = b.cate.UpdateCategory(ctx, tx, cate)
		return err
	}, nil)

	return err
}

func (b *blogService) DeleteCate(ctx context.Context, cateId uint) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = b.cate.DeleteCategory(ctx, tx, cateId)
		return err
	}, nil)

	return err
}

type PostArrWithCount struct {
	Posts []model.BlogPost `json:"posts"`
	Count int64            `json:"count"`
}

func (b *blogService) GetPostLists(ctx context.Context, cateIds []uint, offset, limit int32) (PostArrWithCount, error) {
	log := blog.Extract(ctx)

	var posts []model.BlogPost
	var count int64
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		cond := repository.FindPostArg{}

		if len(cateIds) > 0 {
			cond.CategoryIDs = cateIds
		}

		cond.NoLimit = true

		count, err = b.post.CountPosts(ctx, tx, cond)
		if err != nil {
			return err
		}
		if count == 0 {
			return nil
		}

		if limit != 0 {
			cond.Limit = limit
			cond.Offset = offset
			cond.NoLimit = false
		}

		log.Debug("cond: ", zap.Any("cond", cond))

		posts, err = b.post.FindPosts(ctx, tx, cond)

		for i := range posts {
			posts[i].Body = ""
		}

		return err
	}, nil)

	log.Debug("posts: ", zap.Any("posts", posts))
	log.Debug("count: ", zap.Any("count", count))

	var postArrWithCount PostArrWithCount
	postArrWithCount.Count = count
	postArrWithCount.Posts = posts

	return postArrWithCount, err
}

func (b *blogService) GetPostByPostID(ctx context.Context, id uint) (model.BlogPost, error) {

	var post model.BlogPost
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		post, err = b.post.GetPost(ctx, tx, id)
		return err
	}, nil)

	return post, err
}

func (b *blogService) CreatePost(ctx context.Context, post model.BlogPost) error {

	var err error

	post.CreatedTime = time.Now()
	post.ModifiedTime = time.Now()

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = b.post.CreatePost(ctx, tx, post)
		return err
	}, nil)

	return err
}

func (b *blogService) UpdatePost(ctx context.Context, post model.BlogPost) error {

	var err error
	post.ModifiedTime = time.Now()

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {

		result := tx.Exec("delete from blog_post_tags where blog_post_id = ?", post.ID.ID)
		if result.Error != nil {
			return result.Error
		}

		err = b.post.UpdatePost(ctx, tx, post)
		return err
	}, nil)

	return err
}

func (b *blogService) DeletePost(ctx context.Context, postId uint) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {

		result := tx.Exec("delete from blog_post_tags where blog_post_id = ?", postId)
		if result.Error != nil {
			return result.Error
		}

		err = b.post.DeletePost(ctx, tx, postId)
		return err
	}, nil)

	return err
}

func (b *blogService) GetPostPs(ctx context.Context) ([]model.BlogPostPs, error) {
	var postPs []model.BlogPostPs
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		postPs, err = b.post.FindPostPwd(ctx, tx)
		return err
	}, nil)

	return postPs, err
}

func (b *blogService) CreatePostPs(ctx context.Context, postPs model.BlogPostPs) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = b.post.CreatePwd(ctx, tx, postPs)
		return err
	}, nil)

	return err
}

func (b *blogService) UpdatePostPs(ctx context.Context, postPs model.BlogPostPs) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = b.post.UpdatePostPwd(ctx, tx, postPs)
		return err
	}, nil)

	return err
}

func (b *blogService) DeletePostPs(ctx context.Context, postPsId uint) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = b.post.DeletePostPwd(ctx, tx, postPsId)
		return err
	}, nil)

	return err
}

func (b *blogService) GetTags(ctx context.Context) ([]model.BlogTag, error) {

	var tags []model.BlogTag
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		cond := repository.FindTagArg{
			NoLimit: true,
		}
		tags, err = b.tag.FindTags(ctx, tx, cond)
		return err
	}, nil)

	return tags, err
}

func (b *blogService) CreateTag(ctx context.Context, tag model.BlogTag) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = b.tag.CreateTag(ctx, tx, tag)
		return err
	}, nil)

	return err
}

func (b *blogService) UpdateTag(ctx context.Context, tag model.BlogTag) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = b.tag.UpdateTag(ctx, tx, tag)
		return err
	}, nil)

	return err
}

func (b *blogService) DeleteTag(ctx context.Context, tagId uint) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = b.tag.DeleteTag(ctx, tx, tagId)
		return err
	}, nil)

	return err
}

func (b *blogService) GetChickenSoups(ctx context.Context) ([]model.ChickenSoup, error) {
	var chickenSoups []model.ChickenSoup
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		cond := repository.FindChickenSoupArg{
			NoLimit: true,
		}
		chickenSoups, err = b.chickenSoup.FindChickenSoups(ctx, tx, cond)
		return err
	}, nil)

	return chickenSoups, err
}

func (b *blogService) CreateChickenSoup(ctx context.Context, chickenSoup model.ChickenSoup) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = b.chickenSoup.CreateChickenSoup(ctx, tx, chickenSoup)
		return err
	}, nil)

	return err
}

func (b *blogService) UpdateChickenSoup(ctx context.Context, chickenSoup model.ChickenSoup) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = b.chickenSoup.UpdateChickenSoup(ctx, tx, chickenSoup)
		return err
	}, nil)

	return err
}

func (b *blogService) DeleteChickenSoup(ctx context.Context, chickenSoupId uint) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = b.chickenSoup.DeleteChickenSoup(ctx, tx, chickenSoupId)
		return err
	}, nil)

	return err
}
