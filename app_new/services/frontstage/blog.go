package frontstage

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/DeYu666/blog-backend-service/app_new/common/response"
	"github.com/DeYu666/blog-backend-service/app_new/dao/blog"
	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/gin-gonic/gin"
)

func GetGeneralCate(c *gin.Context) {

	generalCate, err := blog.GetGeneralCategories()
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, generalCate)
}

// GetCategory 是获取所有分类列表
func GetCategory(c *gin.Context) {

	category, err := blog.GetCategories()
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, category)
}

type postNav struct {
	models.BlogGeneralCategories
	Categories []*models.BlogCategories `json:"cub_cate"`
}

// GetPostCategories 是获取所有文章的分类: 大分类->小分类
func GetPostCategories(c *gin.Context) {

	generalCates, err := blog.GetGeneralCategories()
	if err != nil {
		response.SqlFail(c, err.Error())
	}

	nav := make([]postNav, len(generalCates))
	for index, cate := range generalCates {
		nav[index].ID = cate.ID
		nav[index].Name = cate.Name
	}

	for index, tempNav := range nav {
		nav[index].Categories, err = blog.GetCategories(blog.CategoryByGeneralID(tempNav.ID.ID))
	}

	response.Success(c, nav)
}

func GetTags(c *gin.Context) {
	blogTags, err := blog.GetTag()
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, blogTags)
}

// GetPostLists 是获取全部文章列表
func GetPostLists(c *gin.Context) {

	categoryID, err := strconv.Atoi(c.Query("categoryId"))
	if err != nil {
		categoryID = -1
	}
	pageNum, err := strconv.Atoi(c.Query("pageNum"))
	if err != nil {
		pageNum = 0
	}

	posts, err := blog.GetPosts(blog.PostPaging(pageNum, categoryID, 12))
	if err != nil {
		response.SqlFail(c, err.Error())
	}

	tempPost, err := blog.GetPosts(blog.PostByCategoryId(categoryID))

	response.Success(c, gin.H{"postContent": posts, "pagesNum": len(tempPost)})
}

// GetPostListsByTagID 通过标签 ID 获取所属全部文章列表
func GetPostListsByTagID(c *gin.Context) {
	tagID, err := strconv.Atoi(c.Param("tagID"))
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	posts, err := blog.GetPosts(blog.PostByTagId(uint(tagID)))
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	tag, _ := blog.GetTag(blog.TagId(uint(tagID)))

	response.Success(c, gin.H{
		"posts":    posts,
		"tagTitle": tag[0].Name,
	})
}

func GetPostByPostID(c *gin.Context) {
	param := c.Param("postID")

	postID, err := strconv.Atoi(param)
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	blogPosts, err := blog.GetPosts(blog.PostId(uint(postID)))
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	if len(blogPosts) == 0 {
		response.Success(c, nil)
		return
	}

	blogPosts[0].Views += 1
	err = blog.ModifyPost(blogPosts[0])
	if err != nil {
		response.SqlFail(c, err.Error())
	}

	response.Success(c, blogPosts[0])
}

func GetChickenSoup(c *gin.Context) {

	chickenSoups, err := blog.GetChickenSoups()
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	rand.Seed(time.Now().Unix())
	soupId := rand.Intn(len(chickenSoups))
	if soupId > 0 {
		soupId -= 1
	}
	response.Success(c, chickenSoups[soupId])
}

// GetAllPostList 是获取全部文章
func GetAllPostList(c *gin.Context) {
	posts, err := blog.GetPosts(blog.PostOrderByDesc())
	if err != nil {
		response.SqlFail(c, err.Error())
	}

	response.Success(c, posts)
}
