package zhihu

import (
	"fmt"
	"net/http"
)

type Param struct {
	 
}

func PublishArticle(cookie string)  {

	req, _ := http.NewRequest("POST","http://www.zhihu.com/api/v4/document_convert", nil)
	req.Header.Set("Cookie", cookie)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Content-Type", "application/json;")
	req.Header.Set("Referer", "https://zhuanlan.zhihu.com/write")
	req.Header.Set("Origin", "https://zhuanlan.zhihu.com")
	req.Header.Set("X-Requested-With", "Fetch")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	fmt.Println(res)
}