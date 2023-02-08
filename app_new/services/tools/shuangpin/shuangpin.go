package shuangpin

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/DeYu666/blog-backend-service/app_new/common/response"
	"github.com/DeYu666/blog-backend-service/app_new/dao/shuangpin"
	"github.com/DeYu666/blog-backend-service/app_new/models/tools"
	"github.com/gin-gonic/gin"
)

type RetrieveDataByShuangPin struct {
	Cate string `json:"cate"`
	Num  int    `json:"num"`
}

func bufferToStruct(body io.Reader) (*RetrieveDataByShuangPin, error) {

	var bodyBytes []byte
	var err error
	if body != nil {
		bodyBytes, err = ioutil.ReadAll(body)
		if err != nil {
			return nil, err
		}
	}

	fmt.Println(string(bodyBytes))

	model := &RetrieveDataByShuangPin{}

	err = json.Unmarshal(bodyBytes, model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

type ResponseDataShuangPin struct {
	Words []string `json:"word"`
	Pys   []string `json:"pinying"`
	Sps   []string `json:"shuangpin"`
}

func GetShuangpin(c *gin.Context) {

	retreive, err := bufferToStruct(c.Request.Body)
	if err != nil {
		response.ValidateFail(c, err.Error())
	}

	cate := retreive.Cate

	var words []string
	if cate == "all" || cate == "" {
		words, _ = shuangpin.GetWords([]string{}, retreive.Num)
	} else {
		words, _ = shuangpin.GetWords([]string{cate}, retreive.Num)
	}

	var pYs, sPs []string

	for _, word := range words {
		pY, _ := shuangpin.GetPingYin(word)
		sp := shuangpin.PinYinTurnShuangPin(pY)
		sp = shuangpin.CheckResult(word, pY, sp)

		var tPy, tSp string
		for _, p := range pY {
			tPy += p + " "
		}
		for _, s := range sp {
			tSp += s
		}
		pYs = append(pYs, tPy)
		sPs = append(sPs, tSp)
	}

	result := &tools.ShuangPin{Words: words, PYs: pYs, SPs: sPs}

	response.Success(c, result)
}
