package blog

import (
	"testing"

	"github.com/DeYu666/blog-backend-service/app_new/models"
)

func TestPost(t *testing.T) {
	tags, err := GetPosts()

	if len(tags) < 1 {
		t.Errorf("select post data error, it is %v, error is %v", tags, err)
		return
	}

	testName := "test_golang"

	cate, _ := GetCategories()

	post := &models.BlogPost{
		Title:      testName,
		Body:       "123",
		Excerpt:    "qwe",
		CategoryID: int(cate[0].ID.ID),
		AuthorID:   1,
		Views:      0,
		Likes:      0,
		CoverURL:   "",
		TitleURL:   "",
		IsOpen:     false,
	}

	err = AddPost(post)
	if err != nil {
		t.Errorf("insert Data error, it`s %v", err)
		return
	}

	info, err := GetPosts(PostTitle(testName))
	if err != nil || len(info) < 1 {
		t.Errorf("select Data error, error is %v, info is %v", err, info)
		return
	}

	id := info[0].ID.ID

	post = info[0]
	post.Title = "test_modify_golang"

	err = ModifyPost(post)
	if err != nil {
		t.Errorf("modify Data error, it`s %v", err)
		return
	}

	info, err = GetPosts(PostId(id))
	if err != nil || len(info) != 1 || info[0].Title != "test_modify_golang" {
		t.Errorf("modify Data error, error`s %v, info is %v", err, info)
	}

	err = DeletePost(id)
	if err != nil {
		t.Errorf("delete Data error, it`s %v", err)
		return
	}

	info, err = GetPosts(PostId(id))
	if err != nil || len(info) != 0 {
		t.Errorf("delete Data error, error`s %v, info is %v", err, info)
		return
	}
}

func TestPostPwd(t *testing.T) {
	PostPwd, err := GetPostPwd()

	if len(PostPwd) < 1 {
		t.Errorf("select Post password data error, it is %v, error is %v", PostPwd, err)
		return
	}

	testName := "test_golang"

	PostPwd1 := &models.BlogPostPs{Password: testName}

	err = AddPostPwd(PostPwd1)
	if err != nil {
		t.Errorf("insert Data error, it`s %v", err)
	}

	info, err := GetPostPwd(PostPassWord(testName))
	if err != nil || len(info) < 1 {
		t.Errorf("select Data error, error is %v, info is %v", err, info)
		return
	}

	id := info[0].ID.ID

	PostPwd1 = info[0]
	PostPwd1.Password = "test_modify_golang"

	err = ModifyPostPwd(PostPwd1)
	if err != nil {
		t.Errorf("modify Data error, it`s %v", err)
		return
	}

	info, err = GetPostPwd(PostPwdId(id))
	if err != nil || len(info) != 1 || info[0].Password != "test_modify_golang" {
		t.Errorf("modify Data error, error`s %v, info is %v", err, info)
	}

	err = DeletePostPwd(PostPwdId(id))
	if err != nil {
		t.Errorf("delete Data error, it`s %v", err)
		return
	}

	info, err = GetPostPwd(PostPwdId(id))
	if err != nil || len(info) != 0 {
		t.Errorf("delete Data error, error`s %v, info is %v", err, info)
		return
	}
}
