package tool

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/DeYu666/blog-backend-service/app_new/models"
)

type dataType interface {
	*models.BlogCategories | *models.BlogTag | *models.BlogGeneralCategories | *models.BlogPost | *models.ChickenSoup | *models.BooksList | *models.BookContent | *models.ExperienceCv | *models.ProjectCv | *models.SkillCv | *models.Diary | *models.AuthUser | *models.ProjectCvPs | *models.BlogPostPs | *models.DiaryPs | *models.LoveInfo | *models.Memo
}

func BufferToStruct[V dataType](body io.Reader, model V) (V, error) {

	var bodyBytes []byte
	var err error
	if body != nil {
		bodyBytes, err = ioutil.ReadAll(body)
		if err != nil {
			return nil, err
		}
	}

	err = json.Unmarshal(bodyBytes, model)
	if err != nil {
		return nil, err
	}

	return model, nil
}
