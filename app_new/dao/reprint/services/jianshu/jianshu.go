package jianshu

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/DeYu666/blog-backend-service/app_new/models/reprint"
)

func NewParam(title, content string, id string) (*reprint.ParamJianShu, error) {

	param := reprint.ParamJianShu{
		AutoSaveControl: 12,
		Title:           title,
		Content:         content,
		Id:              id,
	}

	return &param, nil
}

func PublishArticle(param *reprint.ParamJianShu, cookie string) (msg string, err error) {

	req, err := http.NewRequest("GET", "https://blog-console-api.csdn.net/v1/postedit/saveArticle", nil)
	req.Header.Set("Cookie", cookie)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Content-Type", "application/json;")
	req.Header.Set("Referer", "https://www.jianshu.com/writer")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送 request 请求失败，具体原因为: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("获取请求中 body 数据失败，具体原因为: %v", err)
	}

	fmt.Println(body)

	return "", nil
}
