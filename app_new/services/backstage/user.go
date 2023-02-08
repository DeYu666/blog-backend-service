package backstage

import (
	"github.com/DeYu666/blog-backend-service/app_new/common/response"
	"github.com/DeYu666/blog-backend-service/app_new/common/tool"
	"github.com/DeYu666/blog-backend-service/app_new/dao/user"
	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/DeYu666/blog-backend-service/global"
	"github.com/gin-gonic/gin"
)

// AutoCvMigrate 自动更新表结构
func AutoUserMigrate(c *gin.Context) {
	result := global.App.DB.AutoMigrate(&models.AuthUser{})
	c.JSON(200, result)
}

// ----- 用户 ----

func GetUser(c *gin.Context) {

	users, err := user.GetUsers()
	if err != nil {
		response.SqlFail(c, err.Error())
	}

	response.Success(c, users)
}

// AddUser 添加用户
func AddUser(c *gin.Context) {

	userData, err := tool.BufferToStruct(c.Request.Body, &models.AuthUser{})
	if err != nil {
		response.ValidateFail(c, err.Error())
	}

	err = user.AddUser(userData)
	if err != nil {
		response.SqlFail(c, err.Error())
	}

	response.Success(c, "成功")
}

// ModifyUser 修改用户
func ModifyUser(c *gin.Context) {

	userData, err := tool.BufferToStruct(c.Request.Body, &models.AuthUser{})
	if err != nil {
		response.ValidateFail(c, err.Error())
	}

	err = user.ModifyUser(userData)
	if err != nil {
		response.SqlFail(c, err.Error())
	}

	response.Success(c, "成功")
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {

	userData, err := tool.BufferToStruct(c.Request.Body, &models.AuthUser{})
	if err != nil {
		response.ValidateFail(c, err.Error())
	}

	err = user.DeleteUser(user.UserId(userData.ID))
	if err != nil {
		response.SqlFail(c, err.Error())
	}

	response.Success(c, "成功")
}
