package blog

import (
	"testing"

	"github.com/DeYu666/blog-backend-service/app_new/models"
)

func TestChickenSoup(t *testing.T) {
	testName := "test_golang"
	data := &models.ChickenSoup{}
	data.Sentence = testName

	err := AddChickenSoup(data)
	if err != nil {
		t.Errorf("insert Data error, it`s %v", err)
		return
	}

	info, err := GetChickenSoups(ChickenSoupSentence(testName))

	if err != nil || len(info) < 1 {
		t.Errorf("select Data error, error is %v, info is %v", err, info)
		return
	}

	id := info[0].ID.ID

	data = info[0]
	data.Sentence = "test_modify_golang"

	err = ModifyChickenSoup(data)

	if err != nil {
		t.Errorf("modify Data error, it`s %v", err)
		return
	}

	info, err = GetChickenSoups(ChickenSoupId(id))

	if err != nil || len(info) != 1 || info[0].Sentence != "test_modify_golang" {
		t.Errorf("modify Data error, error`s %v, info is %v", err, info)
	}

	err = DeleteChickenSoup(ChickenSoupId(id))

	if err != nil {
		t.Errorf("delete Data error, it`s %v", err)
		return
	}

	info, err = GetChickenSoups(ChickenSoupId(id))

	if err != nil || len(info) != 0 {
		t.Errorf("delete Data error, error`s %v, info is %v", err, info)
		return
	}
}
