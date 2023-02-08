package backstage

import (
	"github.com/DeYu666/blog-backend-service/app_new/common/response"
	"github.com/DeYu666/blog-backend-service/app_new/common/tool"
	"github.com/DeYu666/blog-backend-service/app_new/dao/user"
	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	userData, err := tool.BufferToStruct(c.Request.Body, &models.AuthUser{})
	if err != nil {
		response.ValidateFail(c, err.Error())
	}

	result, err := user.IsExist(userData)
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	if !result {
		response.Fail(c, 10003, "用户名或密码错误")
	}

	tokenData, err, _ := models.JwtService.CreateToken("backstage", userData)
	if err != nil {
		response.Fail(c, 10004, err.Error())
		return
	}

	response.Success(c, tokenData)
}
