package backstage

import (
	"strconv"

	"github.com/DeYu666/blog-backend-service/app_new/common/response"
	"github.com/DeYu666/blog-backend-service/app_new/common/tool"
	"github.com/DeYu666/blog-backend-service/app_new/dao/blog"
	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/DeYu666/blog-backend-service/global"
	"github.com/gin-gonic/gin"
)

// AutoBlogMigrate 自动更新表结构
func AutoBlogMigrate(c *gin.Context) {
	result := global.App.DB.AutoMigrate(&models.ChickenSoup{}, &models.BlogGeneralCategories{}, &models.BlogCategories{}, &models.BlogPost{}, &models.BlogTag{}, &models.AuthUser{}, &models.BlogPostPs{})
	c.JSON(200, result)
}

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

func GetTags(c *gin.Context) {
	blogTags, err := blog.GetTag()
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, blogTags)
}

func GetPostLists(c *gin.Context) {

	blogPosts, err := blog.GetPosts(blog.PostOrderByDesc())
	if err != nil {
		response.SqlFail(c, err.Error())
	}

	for index, _ := range blogPosts {
		blogPosts[index].Body = ""
	}

	response.Success(c, blogPosts)
}

func GetPostByPostID(c *gin.Context) {
	param := c.Param("postId")

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

	response.Success(c, blogPosts[0])
}

func GetChickenSoup(c *gin.Context) {

	chickenSoups, err := blog.GetChickenSoups(blog.ChickenSoupOrderByDesc())
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	response.Success(c, chickenSoups)
}

// 添加

// AddGeneralCate 添加 总类
func AddGeneralCate(c *gin.Context) {

	generalCate, err := tool.BufferToStruct(c.Request.Body, &models.BlogGeneralCategories{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = blog.AddGeneralCate(generalCate)
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, "成功")
}

// AddCategory 添加 分类
func AddCategory(c *gin.Context) {

	category, err := tool.BufferToStruct(c.Request.Body, &models.BlogCategories{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = blog.AddCategory(category)
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, "成功")
}

// AddTags 添加 标签
func AddTags(c *gin.Context) {

	tag, err := tool.BufferToStruct(c.Request.Body, &models.BlogTag{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = blog.AddTag(tag)
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, "成功")

}

// AddPost 添加 文章
func AddPost(c *gin.Context) {

	post, err := tool.BufferToStruct(c.Request.Body, &models.BlogPost{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = blog.AddPost(post)
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, "成功")
}

// AddChickenSoup 添加 鸡汤
func AddChickenSoup(c *gin.Context) {
	chickenSoup, err := tool.BufferToStruct(c.Request.Body, &models.ChickenSoup{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = blog.AddChickenSoup(chickenSoup)
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, "成功")
}

// ModifyGeneralCate 修改 总类
func ModifyGeneralCate(c *gin.Context) {

	generalCate, err := tool.BufferToStruct(c.Request.Body, &models.BlogGeneralCategories{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = blog.ModifyGeneralCate(generalCate)
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, "成功")
}

// ModifyCategory 修改 分类
func ModifyCategory(c *gin.Context) {
	category, err := tool.BufferToStruct(c.Request.Body, &models.BlogCategories{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = blog.ModifyCategory(category)
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, "成功")
}

// ModifyPost 修改 文章
func ModifyPost(c *gin.Context) {

	post, err := tool.BufferToStruct(c.Request.Body, &models.BlogPost{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = blog.ModifyPost(post)
	if err != nil {
		response.SqlFail(c, err.Error())
	}

	response.Success(c, "成功")
}

// ModifyTags 修改 总类
func ModifyTags(c *gin.Context) {
	tag, err := tool.BufferToStruct(c.Request.Body, &models.BlogTag{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = blog.ModifyTag(tag)
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, "成功")

}

// ModifyChickenSoup 修改 鸡汤
func ModifyChickenSoup(c *gin.Context) {
	chickenSoup, err := tool.BufferToStruct(c.Request.Body, &models.ChickenSoup{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = blog.ModifyChickenSoup(chickenSoup)
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, "成功")
}

// DeleteGeneralCate 删除总类
func DeleteGeneralCate(c *gin.Context) {

	generalCate, err := tool.BufferToStruct(c.Request.Body, &models.BlogGeneralCategories{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = blog.DeleteGeneralCate(blog.GeneralCateId(generalCate.ID.ID))
	if err != nil {
		/**
		  能查到 关联文章的 id 么？ 能的话，最好返回去，前端可以做进一步的处理。
		*/
		response.Fail(c, 10002, "请先删除关联文章")
		return
	}

	response.Success(c, "成功")
}

// DeleteTags 删除标签
func DeleteTags(c *gin.Context) {

	tag, err := tool.BufferToStruct(c.Request.Body, &models.BlogTag{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = blog.DeleteTag(blog.TagId(tag.ID.ID))
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, "成功")
}

// DeleteCategory 删除分类
func DeleteCategory(c *gin.Context) {

	category, err := tool.BufferToStruct(c.Request.Body, &models.BlogCategories{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = blog.DeleteCategory(blog.CateId(category.ID.ID))
	if err != nil {
		/**
		  能查到 关联文章的 id 么？ 能的话，最好返回去，前端可以做进一步的处理。
		*/
		response.Fail(c, 10002, "请先删除关联文章")
		return
	}

	response.Success(c, "成功")

}

// DeletePost 删除文章
func DeletePost(c *gin.Context) {

	post, err := tool.BufferToStruct(c.Request.Body, &models.BlogPost{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = blog.DeletePost(post.ID.ID)
	if err != nil {
		response.SqlFail(c, err.Error())
	}

	response.Success(c, "成功")

}

// DeleteChickenSoup 删除鸡汤
func DeleteChickenSoup(c *gin.Context) {
	chickenSoup, err := tool.BufferToStruct(c.Request.Body, &models.ChickenSoup{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = blog.DeleteChickenSoup(blog.ChickenSoupId(chickenSoup.ID.ID))
	if err != nil {
		response.SqlFail(c, err.Error())
	}

	response.Success(c, "成功")

}

func GetBlogPostPs(c *gin.Context) {

	postPs, err := blog.GetPostPwd()
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, postPs)
}

func ModifyBlogPostPs(c *gin.Context) {

	postPs, err := tool.BufferToStruct(c.Request.Body, &models.BlogPostPs{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = blog.ModifyPostPwd(postPs)
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, "成功")
}
