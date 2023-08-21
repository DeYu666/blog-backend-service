package client

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type TuChuang interface {
	UploadFile(image bytes.Buffer, fileName string) (string, error)
	UploadFileByBase64(imgBase64 string) (string, error)
	UploadFileByBase64AndFileName(imgBase64 string, fileName string) (string, error)
}

type QiNiuYunConfiguration struct {
	AccessKey string `mapstructure:"accessKey" json:"accessKey" yaml:"accessKey"`
	SecretKey string `mapstructure:"secretKey" json:"secretKey" yaml:"secretKey"`
	Bucket    string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
}

const (
	domain = "http://image.asa-zhang.top/"
)

type qiniuyun struct {
	token     string // token 信息
	expires   int64  // 截止时间
	accessKey string
	secretKey string
	bucket    string
}

func NewQiniuyun(accessKey, secretKey, bucket string) TuChuang {
	q := &qiniuyun{
		accessKey: accessKey,
		secretKey: secretKey,
		bucket:    bucket,
	}

	q.getToken()
	return q
}

func (q *qiniuyun) UploadFile(image bytes.Buffer, fileName string) (string, error) {

	timestamp := time.Now().Format("20060102150405")

	uniqueFileName := fmt.Sprintf("%s_%s", timestamp, fileName)

	fileName = "blog/post/" + uniqueFileName

	err := q.uploadFileByPath(&image, fileName)
	if err != nil {
		return "", err
	}

	return domain + fileName, err
}

func (q *qiniuyun) UploadFileByBase64(imgBase64 string) (string, error) {
	return q.UploadFileByBase64AndFileName(imgBase64, "")
}

func (q *qiniuyun) UploadFileByBase64AndFileName(imgBase64 string, fileName string) (string, error) {

	if fileName == "" {
		fileName = getFileNameFromData([]byte(imgBase64))
	}

	fileName = "blog/post/" + fileName

	ddd, _ := base64.StdEncoding.DecodeString(strings.Split(imgBase64, ",")[1]) //成图片文件并把文件写入到buffer

	err := q.uploadFileByPath(bytes.NewBuffer(ddd), fileName)
	if err != nil {
		return "", err
	}

	return domain + fileName, nil
}

func (q *qiniuyun) getToken() {
	// 判断 token 是否过期，过期则需重新申请, 默认的过期时间是 3600 s
	if q.token != "" && time.Now().Unix()-q.expires < 3600 {
		return
	}

	putPolicy := storage.PutPolicy{
		Scope: q.bucket,
	}
	mac := qbox.NewMac(q.accessKey, q.secretKey)
	upToken := putPolicy.UploadToken(mac) // 返回的数据：{"hash":"Ftgm-CkWePC90fzMBTRNmPMhGBcSV","key":"qiniu.jpg"}

	q.token = upToken
	q.expires = time.Now().Unix()
}

func (q *qiniuyun) uploadFileByPath(file *bytes.Buffer, fileName string) error {

	fmt.Println("uploadFileByPath", file, fileName)

	q.getToken()

	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuanan
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	fmt.Println(q.token)

	putExtra := storage.PutExtra{}

	dataLen := int64(file.Len())
	// fmt.Println(dataLen, fileName)
	err := formUploader.Put(context.Background(), &ret, q.token, fileName, file, dataLen, &putExtra)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// fmt.Println("qiniuyun return: ", ret.Key, ret.Hash)
	return nil
}

func getFileNameFromData(data []byte) string {

	Sha1Inst := md5.New()
	Sha1Inst.Write(data)
	fileName := Sha1Inst.Sum([]byte(""))
	tempData := binary.LittleEndian.Uint32(fileName)

	fileName1 := strconv.Itoa(int(tempData))

	return fileName1
}
