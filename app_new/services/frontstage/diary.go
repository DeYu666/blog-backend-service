package frontstage

import (
	"github.com/DeYu666/blog-backend-service/app_new/common/response"
	"github.com/DeYu666/blog-backend-service/app_new/dao/diary"
	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/gin-gonic/gin"
)

func GetDiaryByPs(c *gin.Context) {

	diaryPassword := c.Query("diaryPassword")

	diaryPs, err := diary.GetDiaryPwd()
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	diaries, err := diary.GetDiaries(diary.DiaryOrderByDesc())
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	if isPsTrue(diaryPassword, diaryPs) {
		response.Success(c, diaries)
		return
	}

	for i, data := range diaries {
		if data.IsOpen == false {
			lenContent := len(data.Content)
			diaries[i].Content = strDouble("*", lenContent)
		}
	}

	response.Success(c, diaries)

}

// 判断密码是否正确
func isPsTrue(pass string, diaryPs []*models.DiaryPs) bool {
	for _, data := range diaryPs {
		if data.Password == pass {
			return true
		}
	}
	return false
}

// strDouble 将字符串翻倍，参数：字符串、翻倍的倍数
func strDouble(str string, n int) string {
	if n == 0 {
		return str
	}
	res := str
	for i := 0; i < n; i++ {
		res += str
	}
	return res
}
