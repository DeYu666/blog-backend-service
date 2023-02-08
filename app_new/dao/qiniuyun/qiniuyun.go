package qiniuyun

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

func getToken() string {
	accessKey := "KIwd3YB6r2SPkDpzueop9cub9lMuHjj6vqCKDkw_"
	secretKey := "hWgHB90E5E8yQprzKu8-VNAZWS8-IImXtR3AdCXF"

	bucket := "myfreespacep"
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	// 返回的数据：{"hash":"Ftgm-CkWePC9fzMBTRNmPMhGBcSV","key":"qiniu.jpg"}
	return upToken
}

type SingleToen struct {
	token   string // token 信息
	expires int64  // 截止时间
}

var token_instance *SingleToen

func getToken_instance() *SingleToen {
	if token_instance == nil {
		token_instance = &SingleToen{
			token:   getToken(),
			expires: time.Now().Unix(),
		}
		return token_instance
	}

	// 判断 token 是否过期，过期则需重新申请, 默认的过期时间是 3600 s

	if time.Now().Unix()-token_instance.expires >= 3600 {
		token_instance = &SingleToen{
			token:   getToken(),
			expires: time.Now().Unix(),
		}
		return token_instance
	}

	return token_instance
}

/**
UploadImg 上传图片到七牛云，参数如下：
localFile 是要上传的文件的字节数组;
key 是要上传的文件访问路径;
*/
func UploadImg(localFile []byte, key string) error {
	upToken := getToken_instance()

	cfg := storage.Config{}
	cfg.Zone = &storage.ZoneHuanan
	cfg.UseHTTPS = false
	cfg.UseCdnDomains = false

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	err := formUploader.Put(context.Background(), &ret, upToken.token, key, bytes.NewBuffer(localFile), int64(len(localFile)), nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(ret.Key, ret.Hash)

	return nil
}

func UploadImgByUrl(pic []byte) error {
	// var pic = "填写你的base64后的字符串"

	upToken := getToken_instance()

	data := bytes.NewReader(pic)

	req, err := http.NewRequest("POST", "https://blog-console-api.csdn.net/v1/postedit/saveArticle", data)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Authorization", upToken.token)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("发送 request 请求失败，具体原因为: %v", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("获取请求中 body 数据失败，具体原因为: %v", err)
	}

	fmt.Println(string(body))

	return nil
}

func UploadImgByUrl2(data string, key string) error {

	ddd, _ := base64.StdEncoding.DecodeString(data) //成图片文件并把文件写入到buffer

	upToken := getToken_instance().token

	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuanan
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	dataLen := int64(len(ddd))
	err := formUploader.Put(context.Background(), &ret, upToken, key, bytes.NewBuffer(ddd), dataLen, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("qiniuyun return: ", ret.Key, ret.Hash)

	return nil
}
