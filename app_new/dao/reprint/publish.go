package reprint

import "fmt"

func PublicArticleTo(title, content, site string) error {
	if len(title) == 0 || len(content) == 0 || len(site) == 0 {
		return fmt.Errorf("参数不能为空！")
	}

	/*
	* 上传文章
	 */
	switch site {
	case "CSDN":
		break
	case "掘金":
		break
	case "博客园":
		break
	}
	return nil
}
