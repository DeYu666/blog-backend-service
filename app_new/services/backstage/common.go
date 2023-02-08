package backstage

import (
	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/DeYu666/blog-backend-service/global"

	"github.com/gin-gonic/gin"
)

// AutoMigrate 自动更新表结构
func AutoMigrate(c *gin.Context) {
	// // blog
	// result := global.App.DB.AutoMigrate(&models.ChickenSoup{}, &models.BlogGeneralCategories{}, &models.BlogCategories{}, &models.BlogPost{}, &models.BlogTag{}, &models.AuthUser{}, &models.BlogPostPs{})
	// // book
	// result = global.App.DB.AutoMigrate(&models.BookContent{}, &models.BooksList{})
	// // diary
	// result = global.App.DB.AutoMigrate(&models.Diary{}, &models.DiaryPs{})
	// // user
	// result = global.App.DB.AutoMigrate(&models.AuthUser{})
	// // love
	// result = global.App.DB.AutoMigrate(&models.LoveInfo{})

	// memo
	// result := global.App.DB.AutoMigrate(&models.Memo{})

	// Cv
	result := global.App.DB.AutoMigrate(&models.ExperienceCv{}, &models.SkillCv{}, &models.ProjectCv{}, &models.ProjectCvPs{})

	c.JSON(200, result)
}
