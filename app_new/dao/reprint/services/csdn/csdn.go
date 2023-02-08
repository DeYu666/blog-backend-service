package csdn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/DeYu666/blog-backend-service/app_new/models/reprint"
)

type Option func(*reprint.ParamCSDN)

func Description(description string) Option {
	return func(param *reprint.ParamCSDN) {
		param.Description = description
	}
}

func NewParam(title, content, tags string, options ...func(*reprint.ParamCSDN)) (*reprint.ParamCSDN, error) {
	param := reprint.ParamCSDN{
		CoverImages:  make([]string, 0),
		ReadType:     "public",
		Type:         "original",
		Source:       "pc_postedit",
		NotAutoSaved: 1,
		IsNew:        1,
		Level:        "1",
		Title:        title,
		Content:      content,
		Tags:         tags,
	}

	for _, option := range options {
		option(&param)
	}

	return &param, nil
}

func PublishArticle(param *reprint.ParamCSDN, cookie string) (msg string, err error) {

	paramByte, err := json.Marshal(param)

	if err != nil {
		return "", fmt.Errorf("解析 CSDN 的参数出现了问题， 具体原因为: %v", err)
	}

	data := bytes.NewReader(paramByte)

	req, err := http.NewRequest("POST", "https://blog-console-api.csdn.net/v1/postedit/saveArticle", data)
	req.Header.Set("Cookie", cookie)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Content-Type", "application/json;")

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

	result := reprint.ReceiveCSDN{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		// handle error
		return "", fmt.Errorf("解析结果数据失败，具体原因为: %v", err)
	}

	if result.Code == 200 {
		return result.Msg, nil
	} else {
		return "", fmt.Errorf(result.Msg)
	}

}
