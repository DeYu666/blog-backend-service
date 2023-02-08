package backstage

import (
	"strconv"
	"time"

	"github.com/DeYu666/blog-backend-service/app_new/common/response"
	"github.com/DeYu666/blog-backend-service/app_new/common/tool"
	"github.com/DeYu666/blog-backend-service/app_new/dao/memo"
	"github.com/DeYu666/blog-backend-service/app_new/models"

	"github.com/gin-gonic/gin"
)

type ResultMemos struct {
	CreateTime time.Time     `json:"create_time"`
	Content    []models.Memo `json:"content"`
}

func GetMemos(c *gin.Context) {

	param := c.Param("statusId") // 0 表示未完成； 1 表示已完成；2 表示全部

	memoStatusId, err := strconv.Atoi(param)
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	memos, err := memo.GetMemos(memo.MemoStatus(memoStatusId), memo.MemoOrderByDesc())
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	if len(memos) < 1 {
		response.Success(c, nil)
		return
	}
	resultData := make([]ResultMemos, 1)

	resultData[0].CreateTime = memos[0].CreatedTime

	for _, m := range memos {
		if !isSameDay(resultData[len(resultData)-1].CreateTime, m.CreatedTime) {
			resultData = append(resultData, ResultMemos{CreateTime: m.CreatedTime})
		}
		resultData[len(resultData)-1].Content = append(resultData[len(resultData)-1].Content, *m)
	}

	response.Success(c, resultData)
}

func isSameDay(a time.Time, b time.Time) bool {
	if a.Year() == b.Year() {
		if a.Month().String() == b.Month().String() {
			if a.Day() == b.Day() {
				return true
			}
		}
	}

	return false
}

// AddMemo 添加备忘录内容
func AddMemo(c *gin.Context) {

	momoContent, err := tool.BufferToStruct(c.Request.Body, &models.Memo{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = memo.AddMemo(momoContent)
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, "成功")
}

// ModifyMemoStatus 修改备忘录状态
func ModifyMemoStatus(c *gin.Context) {

	memoContent, err := tool.BufferToStruct(c.Request.Body, &models.Memo{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	updateMemos, err := memo.GetMemos(memo.MemoContent(memoContent.Content))
	if len(updateMemos) < 1 || err != nil {
		response.ValidateFail(c, "查询不到该备忘录内容")
		return
	}

	updateMemo := updateMemos[0]
	updateMemo.Status = 1

	err = memo.ModifyMemo(updateMemo)
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, "成功")
}

// DeleteMemo 删除备忘录内容
func DeleteMemo(c *gin.Context) {

	memoContent, err := tool.BufferToStruct(c.Request.Body, &models.Memo{})
	if err != nil {
		response.ValidateFail(c, err.Error())
	}

	err = memo.DeleteMemo(memo.MemoId(memoContent.ID.ID))
	if err != nil {
		response.SqlFail(c, err.Error())
	}

	response.Success(c, "成功")
}
