package utils

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/DeYu666/blog-backend-service/model"
)

type dataType interface {
	*model.BlogCategories | *model.BlogTag | *model.BlogGeneralCategories | *model.BlogPost | *model.ChickenSoup | *model.BooksList | *model.BookContent | *model.ExperienceCv | *model.ProjectCv | *model.SkillCv | *model.Diary | *model.AuthUser | *model.ProjectCvPs | *model.BlogPostPs | *model.DiaryPs | *model.LoveInfo | *model.Memo
}

func BufferToStruct[V dataType](body io.Reader, model V) (V, error) {

	var bodyBytes []byte
	var err error
	if body != nil {
		bodyBytes, err = io.ReadAll(body)
		if err != nil {
			return nil, err
		}
	}

	fmt.Println("bodyBytes", string(bodyBytes))

	err = json.Unmarshal(bodyBytes, model)
	if err != nil {
		return nil, err
	}

	return model, nil
}
