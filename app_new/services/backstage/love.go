package backstage

import (
	"fmt"

	"github.com/DeYu666/blog-backend-service/app_new/common/response"
	"github.com/DeYu666/blog-backend-service/app_new/common/tool"
	"github.com/DeYu666/blog-backend-service/app_new/dao/love"
	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/gin-gonic/gin"
)

func GetLoveInfo(c *gin.Context) {

	loveInfo, err := love.GetLoveInfos()
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, loveInfo)
}

func AddLoveInfo(c *gin.Context) {

	loveInfo, err := tool.BufferToStruct(c.Request.Body, &models.LoveInfo{})

	fmt.Println(loveInfo)

	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = love.AddLoveInfo(loveInfo)
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, "成功")
}

func ModifyLoveInfo(c *gin.Context) {
	loveInfo, err := tool.BufferToStruct(c.Request.Body, &models.LoveInfo{})

	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = love.ModifyLoveInfo(loveInfo)
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, "成功")
}
