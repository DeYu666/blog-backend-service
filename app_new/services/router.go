package services

import (
	"github.com/DeYu666/blog-backend-service/app_new/services/backstage"
	"github.com/DeYu666/blog-backend-service/app_new/services/frontstage"
	"github.com/DeYu666/blog-backend-service/app_new/services/middleware"
	"github.com/DeYu666/blog-backend-service/app_new/services/tools/shuangpin"

	"github.com/gin-gonic/gin"
)

func RoutersBackStage(router *gin.Engine) {

	router.GET("/autoMigrate1997", backstage.AutoMigrate)

	router.POST("/login", backstage.Login)

	cv := router.Group("/cv").Use(middleware.JWTAuth("backstage"))
	{
		cv.GET("/autoCvMigrate", backstage.AutoCvMigrate)

		cv.GET("/getCvExperience", backstage.GetCvExperience)
		cv.POST("/addCvExperience", backstage.AddCvExperience)
		cv.POST("/modifyCvExperience", backstage.ModifyCvExperience)
		cv.POST("/deleteCvExperience", backstage.DeleteCvExperience)

		cv.GET("/getCvSkill", backstage.GetCvSkill)
		cv.POST("/addCvSkill", backstage.AddCvSkill)
		cv.POST("/modifyCvSkill", backstage.ModifyCvSkill)
		cv.POST("/deleteCvSkill", backstage.DeleteCvSkill)

		cv.GET("/getCvProject", backstage.GetCvProject)
		cv.POST("/addCvProject", backstage.AddCvProject)
		cv.POST("/modifyCvProject", backstage.ModifyCvProject)
		cv.POST("/deleteCvProject", backstage.DeleteCvProject)
		cv.GET("/getCvProjectById/:cvProjectId", backstage.GetCvProjectByID)

		cv.GET("/getCvProjectPs", backstage.GetCvProjectPs)
		cv.POST("/modifyCvProjectPs", backstage.ModifyCvProjectPs)
	}

	// 博客路由
	blog := router.Group("/blog").Use(middleware.JWTAuth("backstage"))
	{
		blog.GET("/autoBlogMigrate", backstage.AutoBlogMigrate)

		// 查看内容路由
		blog.GET("/getGeneralCate", backstage.GetGeneralCate)
		blog.GET("/getCategory", backstage.GetCategory)
		blog.GET("/getTags", backstage.GetTags)
		blog.GET("/getChickenSoup", backstage.GetChickenSoup)
		blog.GET("/getPostsList", backstage.GetPostLists)
		blog.GET("/getPostByPostId/:postId", backstage.GetPostByPostID)
		blog.GET("/getBlogPostPs", backstage.GetBlogPostPs)

		// 添加内容路由
		//blog.POST("/addGeneralCate", JWTAuthMiddleware(), frontbackstage.AddGeneralCate)
		blog.POST("/addGeneralCate", backstage.AddGeneralCate)
		blog.POST("/addTags", backstage.AddTags)
		blog.POST("/addCategory", backstage.AddCategory)
		blog.POST("/addChickenSoup", backstage.AddChickenSoup)
		blog.POST("/addPost", backstage.AddPost)

		// 修改内容路由
		blog.POST("/modifyGeneralCate", backstage.ModifyGeneralCate)
		blog.POST("/modifyTags", backstage.ModifyTags)
		blog.POST("/modifyCategory", backstage.ModifyCategory)
		blog.POST("/modifyChickenSoup", backstage.ModifyChickenSoup)
		blog.POST("/modifyPost", backstage.ModifyPost)
		blog.POST("/modifyBlogPostPs", backstage.ModifyBlogPostPs)

		// 删除内容路由
		blog.POST("/deleteGeneralCate", backstage.DeleteGeneralCate)
		blog.POST("/deleteTags", backstage.DeleteTags)
		blog.POST("/deleteCategory", backstage.DeleteCategory)
		blog.POST("/deleteChickenSoup", backstage.DeleteChickenSoup)
		blog.POST("/deletePost", backstage.DeletePost)
	}

	// 书籍路由
	book := router.Group("/book").Use(middleware.JWTAuth("backstage"))
	{
		book.GET("/autoBookMigrate", backstage.AutoBookMigrate)

		book.GET("/getBooksList", backstage.GetBooksList)
		book.POST("/addBooksList", backstage.AddBooksList)
		book.POST("/modifyBooksList", backstage.ModifyBooksList)
		book.POST("/deleteBooksList", backstage.DeleteBooksList)

		book.GET("/getBookContentByBookId/:bookId", backstage.GetBookContent)
		book.POST("/addBookContent", backstage.AddBookContent)
		book.POST("/modifyBookContent", backstage.ModifyBookContent)
		book.POST("/deleteBookContent", backstage.DeleteBookContent)
	}

	diary := router.Group("/diary").Use(middleware.JWTAuth("backstage"))
	{
		diary.GET("/autoDiaryMigrate", backstage.AutoDiaryMigrate)
		diary.GET("/getDiary", backstage.GetDiary)
		diary.POST("/addDiary", backstage.AddDiary)
		diary.POST("/modifyDiary", backstage.ModifyDiary)
		diary.POST("/deleteDiary", backstage.DeleteDiary)

		diary.GET("/getDiaryPs", backstage.GetDiaryPs)
		diary.POST("/modifyDiaryPs", backstage.ModifyDiaryPs)
	}

	user := router.Group("/user").Use(middleware.JWTAuth("backstage"))
	{
		user.GET("/getUser", backstage.GetUser)
		user.POST("/addUser", backstage.AddUser)
		user.POST("/modifyUser", backstage.ModifyUser)
		user.POST("/deleteUser", backstage.DeleteUser)
	}

	love := router.Group("/love").Use(middleware.JWTAuth("backstage"))
	{
		love.GET("/getLoveInfo", backstage.GetLoveInfo)
		love.POST("/addLoveInfo", backstage.AddLoveInfo)
		love.POST("/modifyLoveInfo", backstage.ModifyLoveInfo)
	}

	memo := router.Group("/memo").Use(middleware.JWTAuth("backstage"))
	{
		memo.GET("/getMemos/:statusId", backstage.GetMemos)
		memo.POST("/addMemoContent", backstage.AddMemo)
		memo.POST("/modifyMemoStatus", backstage.ModifyMemoStatus)
	}

	tuchuang := router.Group("/tuchuang")
	{
		tuchuang.POST("/uploadImageFromPost", backstage.UploadImgFromPost)
	}

}

func RoutersFrontStage(router *gin.Engine) {
	front := router.Group("/front")

	cv := front.Group("/cv")
	{
		cv.GET("/getCvExperience", frontstage.GetCvExperience)
		cv.GET("/getCvSkill", frontstage.GetCvSkill)
		cv.GET("/getCvProject", frontstage.GetCvProject)
		cv.GET("/getCvProjectById/:cvProjectId", frontstage.GetCvProjectByID)
	}

	blog := front.Group("/blog")
	{
		// 查看内容路由
		blog.GET("/getGeneralCate", frontstage.GetGeneralCate)
		blog.GET("/getCategory", frontstage.GetCategory)
		blog.GET("/getTags", frontstage.GetTags)
		blog.GET("/getPostsList", frontstage.GetPostLists)
		//blog.GET("/getPostByPostId/:postId", frontstage.GetPostByPostID)

		blog.GET("/getChickenSoup", frontstage.GetChickenSoup)                  // dv
		blog.GET("/getPostCategories", frontstage.GetPostCategories)            // dv
		blog.GET("/getPostsListByPage", frontstage.GetPostLists)                // dv
		blog.GET("/getPostListsByTagID/:tagID", frontstage.GetPostListsByTagID) // dv
		blog.GET("/getPostByPostID/:postID", frontstage.GetPostByPostID)        // dv
		blog.GET("/getAllPostList", frontstage.GetAllPostList)                  // dv
	}

	book := front.Group("/book")
	{
		book.GET("/getBookContentByBookId/:bookId", frontstage.GetBookContent)

		book.GET("/getBooksListByCateId", frontstage.GetBooksListByCateId) // dv
	}

	diary := front.Group("diary")
	{
		diary.GET("/getDiary", frontstage.GetDiaryByPs)
	}

}

func RoutersShuangPin(router *gin.Engine) {
	tools := router.Group("/tools")
	{
		tools.POST("/shuangpin", shuangpin.GetShuangpin)
	}
}
