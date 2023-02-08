package backstage

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/DeYu666/blog-backend-service/app_new/common/response"
	"github.com/DeYu666/blog-backend-service/app_new/dao/qiniuyun"

	"github.com/gin-gonic/gin"
)

type RetrieveImg struct {
	ImgBase64 string `json:"imgBase64"`
}

func UploadImgFromPost(c *gin.Context) {

	data := &RetrieveImg{}

	bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(bodyBytes, data)
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	Sha1Inst := md5.New()
	Sha1Inst.Write([]byte(data.ImgBase64))
	fileName := Sha1Inst.Sum([]byte(""))
	temp_data := binary.LittleEndian.Uint32(fileName)

	fileName1 := strconv.Itoa(int(temp_data))

	fmt.Println(strings.Split(data.ImgBase64, ",")[1])
	fmt.Println(fileName1 + ".jpg")

	path := "blog/post/" + fileName1

	err = qiniuyun.UploadImgByUrl2(strings.Split(data.ImgBase64, ",")[1], path)

	if err != nil {
		fmt.Println(err)

		response.Fail(c, 4000401, err.Error())

		return
	}

	// err = qiniuyun.UploadImgByUrl([]byte(strings.Split(data.ImgBase64, ",")[1]), fileName1+".jpg")

	response.Success(c, "http://image.asa-zhang.top/"+path)
}
