package backstage

import (
	"github.com/DeYu666/blog-backend-service/app_new/common/response"
	"github.com/DeYu666/blog-backend-service/app_new/common/tool"
	"github.com/DeYu666/blog-backend-service/app_new/dao/diary"
	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/DeYu666/blog-backend-service/global"
	"github.com/gin-gonic/gin"
)

// AutoDiaryMigrate 自动更新表结构
func AutoDiaryMigrate(c *gin.Context) {
	result := global.App.DB.AutoMigrate(&models.Diary{}, &models.DiaryPs{})
	c.JSON(200, result)
}

func GetDiary(c *gin.Context) {

	diaries, err := diary.GetDiaries(diary.DiaryOrderByDesc())
	if err != nil {
		response.SqlFail(c, err.Error())
	}

	response.Success(c, diaries)
}

// AddDiary 添加日记内容
func AddDiary(c *gin.Context) {

	diaryContent, err := tool.BufferToStruct(c.Request.Body, &models.Diary{})
	if err != nil {
		response.ValidateFail(c, err.Error())
	}

	err = diary.AddDiary(diaryContent)
	if err != nil {
		response.SqlFail(c, err.Error())
	}

	response.Success(c, "成功")
}

// ModifyDiary 修改日记内容
func ModifyDiary(c *gin.Context) {

	diaryContent, err := tool.BufferToStruct(c.Request.Body, &models.Diary{})
	if err != nil {
		response.ValidateFail(c, err.Error())
	}

	err = diary.ModifyDiary(diaryContent)
	if err != nil {
		response.SqlFail(c, err.Error())
	}

	response.Success(c, "成功")
}

// DeleteDiary 删除日记内容
func DeleteDiary(c *gin.Context) {

	diaryContent, err := tool.BufferToStruct(c.Request.Body, &models.Diary{})
	if err != nil {
		response.ValidateFail(c, err.Error())
	}

	err = diary.DeleteDiary(diary.DiaryId(diaryContent.ID.ID))
	if err != nil {
		response.SqlFail(c, err.Error())
	}

	response.Success(c, "成功")
}

// 密码部分

func GetDiaryPs(c *gin.Context) {

	diaryPs, err := diary.GetDiaryPwd()
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, diaryPs)
}

// ModifyDiaryPs 修改日记观看密码
func ModifyDiaryPs(c *gin.Context) {

	diaryPs, err := tool.BufferToStruct(c.Request.Body, &models.DiaryPs{})
	if err != nil {
		response.ValidateFail(c, err.Error())
	}

	err = diary.ModifyDiaryPwd(diaryPs)
	if err != nil {
		response.SqlFail(c, err.Error())
	}

	response.Success(c, "成功")
}
