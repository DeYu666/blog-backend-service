package qiniuyun_test

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"testing"
)

func TestUploadImage(t *testing.T) {

	imageUrl := "http://image.asa-zhang.top/2.jpg"

	req, _ := http.NewRequest("GET", imageUrl, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.4 Safari/605.1.15")

	client := &http.Client{}
	resp, _ := client.Do(req)
	fmt.Println(resp)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	Sha1Inst := md5.New()
	Sha1Inst.Write(body)
	fileName := Sha1Inst.Sum([]byte(""))
	data := binary.LittleEndian.Uint32(fileName)

	fileName1 := strconv.Itoa(int(data))

	fmt.Println(fileName1 + ".jpg")
	// err := qiniuyun.UploadImg(body, "test/test8.jpg")
	// fmt.Println(err)
}
