package delivery

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/DeYu666/blog-backend-service/config"
	"github.com/DeYu666/blog-backend-service/delivery/middleware"
	"github.com/DeYu666/blog-backend-service/delivery/response"
	"github.com/DeYu666/blog-backend-service/lib/client"
	blog "github.com/DeYu666/blog-backend-service/lib/log"
	"github.com/DeYu666/blog-backend-service/model"
	"github.com/DeYu666/blog-backend-service/service"
	"github.com/DeYu666/blog-backend-service/utils"
	"github.com/gin-gonic/gin"
)

type BlogHandler struct {
	blogService service.BlogService
	tuChuang    client.TuChuang
}

func NewBlogHandler(router *gin.Engine) {
	handle := &BlogHandler{
		blogService: service.NewBlogService(),
		tuChuang:    client.NewQiniuyun(config.Conf.QiNiuYun.AccessKey, config.Conf.QiNiuYun.SecretKey, config.Conf.QiNiuYun.Bucket),
	}

	authRouter := router.Group("/inner").Use(middleware.JWTAuth())

	router.GET("/blog/general_cate", handle.GetGeneralCate)
	authRouter.POST("/blog/general_cate", handle.CreateGeneralCate)
	authRouter.PUT("/blog/general_cate", handle.UpdateGeneralCate)
	authRouter.DELETE("/blog/general_cate/:id", handle.DeleteGeneralCate)

	router.GET("/blog/cate", handle.GetCate)
	authRouter.POST("/blog/cate", handle.CreateCate)
	authRouter.PUT("/blog/cate", handle.UpdateCate)
	authRouter.DELETE("/blog/cate/:id", handle.DeleteCate)

	router.GET("/blog/post", handle.GetPostLists)
	router.GET("/blog/post/:id", handle.GetPostByPostID)
	authRouter.POST("/blog/post", handle.CreatePost)
	authRouter.PUT("/blog/post", handle.UpdatePost)
	authRouter.DELETE("/blog/post/:id", handle.DeletePost)

	authRouter.GET("/blog/post_ps", handle.GetPostPs)
	authRouter.POST("/blog/post_ps", handle.CreatePostPs)
	authRouter.PUT("/blog/post_ps", handle.UpdatePostPs)
	authRouter.DELETE("/blog/post_ps", handle.DeletePostPs)

	router.GET("/blog/tag", handle.GetTags)
	authRouter.POST("/blog/tag", handle.CreateTag)
	authRouter.PUT("/blog/tag", handle.UpdateTag)
	authRouter.DELETE("/blog/tag/:id", handle.DeleteTag)

	router.GET("/blog/chicken_soup", handle.GetChickenSoup)
	authRouter.POST("/blog/chicken_soup", handle.CreateChickenSoup)
	authRouter.PUT("/blog/chicken_soup", handle.UpdateChickenSoup)
	authRouter.DELETE("/blog/chicken_soup/:id", handle.DeleteChickenSoup)

	router.POST("/blog/uploadImage", handle.UploadImage)
}

func (h *BlogHandler) GetGeneralCate(c *gin.Context) {

	log := blog.Extract(c)

	generalCate, err := h.blogService.GetGeneralCategories(c)
	if err != nil {
		log.Error("GetGeneralCate, err " + err.Error())
		response.FailByError(c, response.Errors.ServeError)
		return
	}

	response.Success(c, generalCate)
}

func (h *BlogHandler) CreateGeneralCate(c *gin.Context) {

	log := blog.Extract(c)

	generalCate, err := utils.BufferToStruct(c.Request.Body, &model.BlogGeneralCategories{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = h.blogService.CreateGeneralCate(c, *generalCate)
	if err != nil {
		log.Error("CreateGeneralCate, err " + err.Error())
		response.FailByError(c, response.Errors.ServeError)
		return
	}

	response.Success(c, "success")
}

func (h *BlogHandler) UpdateGeneralCate(c *gin.Context) {

	log := blog.Extract(c)

	generalCate, err := utils.BufferToStruct(c.Request.Body, &model.BlogGeneralCategories{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = h.blogService.UpdateGeneralCate(c, *generalCate)
	if err != nil {
		log.Error("UpdateGeneralCate, err " + err.Error())
		response.FailByError(c, response.Errors.ServeError)
		return
	}

	response.Success(c, "success")
}

func (h *BlogHandler) DeleteGeneralCate(c *gin.Context) {

	log := blog.Extract(c)

	param := c.Param("id")
	cateId, err := strconv.Atoi(param)
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	log.Sugar().Debug("DeleteGeneralCate, cateId " + param)

	err = h.blogService.DeleteGeneralCate(c, uint(cateId))
	if err != nil {
		log.Error("DeleteGeneralCate, err " + err.Error())
		response.FailByError(c, response.Errors.BusinessError)
		return
	}

	response.Success(c, "success")
}

func (h *BlogHandler) GetCate(c *gin.Context) {

	log := blog.Extract(c)

	cate, err := h.blogService.GetCategories(c)
	if err != nil {
		log.Error("GetCate, err " + err.Error())
		response.FailByError(c, response.Errors.ServeError)
		return
	}

	response.Success(c, cate)
}

func (h *BlogHandler) CreateCate(c *gin.Context) {

	log := blog.Extract(c)

	cate, err := utils.BufferToStruct(c.Request.Body, &model.BlogCategories{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = h.blogService.CreateCate(c, *cate)
	if err != nil {
		log.Error("CreateCate, err " + err.Error())
		response.FailByError(c, response.Errors.ServeError)
		return
	}

	response.Success(c, cate)
}

func (h *BlogHandler) UpdateCate(c *gin.Context) {

	log := blog.Extract(c)

	cate, err := utils.BufferToStruct(c.Request.Body, &model.BlogCategories{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = h.blogService.UpdateCate(c, *cate)
	if err != nil {
		log.Error("UpdateCate, err " + err.Error())
		response.FailByError(c, response.Errors.ServeError)
		return
	}

	response.Success(c, "success")
}

func (h *BlogHandler) DeleteCate(c *gin.Context) {

	log := blog.Extract(c)

	param := c.Param("id")
	cateId, err := strconv.Atoi(param)
	if err != nil {
		log.Error("DeleteCate, err " + err.Error())
		response.ValidateFail(c, err.Error())
		return
	}

	err = h.blogService.DeleteCate(c, uint(cateId))
	if err != nil {
		log.Error("DeleteCate, err " + err.Error())
		response.FailByError(c, response.Errors.ServeError)
		return
	}

	response.Success(c, "success")
}

func (h *BlogHandler) GetPostLists(c *gin.Context) {

	log := blog.Extract(c)

	categoryID, err := strconv.Atoi(c.Query("categoryId"))
	if err != nil {
		categoryID = -1
	}
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 0
	}

	var cateIds []uint
	if categoryID >= 0 {
		cateIds = append(cateIds, uint(categoryID))
	}

	postArrWithCount, err := h.blogService.GetPostLists(c, cateIds, int32(offset), int32(limit))
	if err != nil {
		log.Error("GetPostLists, err " + err.Error())
		response.FailByError(c, response.Errors.ServeError)
		return
	}

	response.Success(c, postArrWithCount)
}

func (h *BlogHandler) GetPostByPostID(c *gin.Context) {

	log := blog.Extract(c)

	param := c.Param("id")
	postID, err := strconv.Atoi(param)
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	post, err := h.blogService.GetPostByPostID(c, uint(postID))
	if err != nil {
		log.Error("GetPost, err " + err.Error())
		response.FailByError(c, response.Errors.ServeError)
		return
	}

	response.Success(c, post)
}

func (h *BlogHandler) CreatePost(c *gin.Context) {

	log := blog.Extract(c)

	post, err := utils.BufferToStruct(c.Request.Body, &model.BlogPost{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	authorIdStr := c.MustGet("author_id").(string)

	authorId, err := strconv.Atoi(authorIdStr)
	if err != nil {
		log.Error("CreatePost, err " + err.Error())
		response.ValidateFail(c, err.Error())
		return
	}

	post.AuthorID = authorId

	err = h.blogService.CreatePost(c, *post)
	if err != nil {
		log.Error("CreatePost, err " + err.Error())
		response.FailByError(c, response.Errors.ServeError)
		return
	}

	response.Success(c, "success")
}

func (h *BlogHandler) UpdatePost(c *gin.Context) {

	log := blog.Extract(c)

	post, err := utils.BufferToStruct(c.Request.Body, &model.BlogPost{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	authorIdStr := c.MustGet("author_id").(string)

	authorId, err := strconv.Atoi(authorIdStr)
	if err != nil {
		log.Error("CreatePost, err " + err.Error())
		response.ValidateFail(c, err.Error())
		return
	}

	post.AuthorID = authorId

	log.Sugar().Debugf("post: %v", post)

	err = h.blogService.UpdatePost(c, *post)
	if err != nil {
		log.Error("UpdatePost, err " + err.Error())
		response.FailByError(c, response.Errors.ServeError)
		return
	}

	response.Success(c, "success")
}

func (h *BlogHandler) DeletePost(c *gin.Context) {

	log := blog.Extract(c)

	param := c.Param("id")
	postId, err := strconv.Atoi(param)
	if err != nil {
		log.Error("DeleteCate, err " + err.Error())
		response.ValidateFail(c, err.Error())
		return
	}

	err = h.blogService.DeletePost(c, uint(postId))
	if err != nil {
		log.Error("DeletePost, err " + err.Error())
		response.FailByError(c, response.Errors.ServeError)
		return
	}

	response.Success(c, "success")
}

func (h *BlogHandler) GetPostPs(c *gin.Context) {

	log := blog.Extract(c)

	postP, err := h.blogService.GetPostPs(c)
	if err != nil {
		log.Error("GetPostPs, err " + err.Error())
		response.FailByError(c, response.Errors.ServeError)
		return
	}

	response.Success(c, postP)
}

func (h *BlogHandler) CreatePostPs(c *gin.Context) {

	log := blog.Extract(c)

	postP, err := utils.BufferToStruct(c.Request.Body, &model.BlogPostPs{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = h.blogService.CreatePostPs(c, *postP)
	if err != nil {
		log.Error("CreatePostPs, err " + err.Error())
		response.FailByError(c, response.Errors.ServeError)
		return
	}

	response.Success(c, "success")
}

func (h *BlogHandler) UpdatePostPs(c *gin.Context) {

	log := blog.Extract(c)

	postP, err := utils.BufferToStruct(c.Request.Body, &model.BlogPostPs{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = h.blogService.UpdatePostPs(c, *postP)
	if err != nil {
		log.Error("UpdatePostPs, err " + err.Error())
		response.FailByError(c, response.Errors.ServeError)
		return
	}

	response.Success(c, "success")
}

func (h *BlogHandler) DeletePostPs(c *gin.Context) {

	log := blog.Extract(c)

	postP, err := utils.BufferToStruct(c.Request.Body, &model.BlogPostPs{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = h.blogService.DeletePostPs(c, postP.ID.ID)
	if err != nil {
		log.Error("DeletePostPs, err " + err.Error())
		response.FailByError(c, response.Errors.ServeError)
		return
	}

	response.Success(c, "success")
}

func (h *BlogHandler) GetTags(c *gin.Context) {

	log := blog.Extract(c)

	tags, err := h.blogService.GetTags(c)
	if err != nil {
		log.Error("GetTags, err " + err.Error())
		response.FailByError(c, response.Errors.ServeError)
		return
	}

	response.Success(c, tags)
}

func (h *BlogHandler) CreateTag(c *gin.Context) {

	log := blog.Extract(c)

	tag, err := utils.BufferToStruct(c.Request.Body, &model.BlogTag{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = h.blogService.CreateTag(c, *tag)
	if err != nil {
		log.Error("CreateTag, err " + err.Error())
		response.FailByError(c, response.Errors.ServeError)
		return
	}

	response.Success(c, "success")
}

func (h *BlogHandler) UpdateTag(c *gin.Context) {

	log := blog.Extract(c)

	tag, err := utils.BufferToStruct(c.Request.Body, &model.BlogTag{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = h.blogService.UpdateTag(c, *tag)
	if err != nil {
		log.Error("UpdateTag, err " + err.Error())
		response.FailByError(c, response.Errors.ServeError)
		return
	}

	response.Success(c, "success")
}

func (h *BlogHandler) DeleteTag(c *gin.Context) {

	log := blog.Extract(c)

	param := c.Param("id")
	tagId, err := strconv.Atoi(param)
	if err != nil {
		log.Error("DeleteCate, err " + err.Error())
		response.ValidateFail(c, err.Error())
		return
	}

	err = h.blogService.DeleteTag(c, uint(tagId))
	if err != nil {
		log.Error("DeleteTag, err " + err.Error())
		response.FailByError(c, response.Errors.ServeError)
		return
	}

	response.Success(c, "success")
}

func (h *BlogHandler) GetChickenSoup(c *gin.Context) {

	log := blog.Extract(c)

	chickenSoup, err := h.blogService.GetChickenSoups(c)
	if err != nil {
		log.Error("GetChickenSoup, err " + err.Error())
		response.FailByError(c, response.Errors.ServeError)
		return
	}

	response.Success(c, chickenSoup)
}

func (h *BlogHandler) CreateChickenSoup(c *gin.Context) {

	log := blog.Extract(c)

	chickenSoup, err := utils.BufferToStruct(c.Request.Body, &model.ChickenSoup{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = h.blogService.CreateChickenSoup(c, *chickenSoup)
	if err != nil {
		log.Error("CreateChickenSoup, err " + err.Error())
		response.FailByError(c, response.Errors.ServeError)
		return
	}

	response.Success(c, "success")

}

func (h *BlogHandler) UpdateChickenSoup(c *gin.Context) {

	log := blog.Extract(c)

	chickenSoup, err := utils.BufferToStruct(c.Request.Body, &model.ChickenSoup{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = h.blogService.UpdateChickenSoup(c, *chickenSoup)
	if err != nil {
		log.Error("UpdateChickenSoup, err " + err.Error())
		response.FailByError(c, response.Errors.ServeError)
		return
	}

	response.Success(c, "success")

}

func (h *BlogHandler) DeleteChickenSoup(c *gin.Context) {

	log := blog.Extract(c)

	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		log.Error("DeleteCate, err " + err.Error())
		response.ValidateFail(c, err.Error())
		return
	}

	err = h.blogService.DeleteChickenSoup(c, uint(id))
	if err != nil {
		log.Error("DeleteChickenSoup, err " + err.Error())
		response.FailByError(c, response.Errors.ServeError)
		return
	}

	response.Success(c, "success")

}

type ResponseUploadImg struct {
	FileName string `json:"file_name"`
	Url      string `json:"url"`
}

func (h *BlogHandler) UploadImage(c *gin.Context) {

	log := blog.Extract(c)

	form, err := c.MultipartForm()
	if err != nil {
		log.Error("UploadImage, err " + err.Error())
		response.FailByError(c, response.Errors.ServeError)
		return
	}
	files := form.File["file[]"]

	var urls []ResponseUploadImg

	for _, file := range files {
		uploadedFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer uploadedFile.Close()

		buffer := bytes.NewBuffer(nil)
		if _, err := io.Copy(buffer, uploadedFile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		fileName := file.Filename

		url, err := h.tuChuang.UploadFile(*buffer, fileName)
		if err != nil {
			log.Error("UploadImage, err " + err.Error())
			response.FailByError(c, response.Errors.ServeError)
			return
		}
		log.Info("UploadImage, url " + url)
		urls = append(urls, ResponseUploadImg{fileName, url})
	}

	response.Success(c, urls)

	// response.Success(c, url)
}
