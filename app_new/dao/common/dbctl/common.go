package dbctl

import (
	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/DeYu666/blog-backend-service/global"
	"gorm.io/gorm"
)

type blogDataType interface {
	*models.BlogCategories | *models.BlogTag | *models.BlogGeneralCategories | *models.BlogPost | *models.ChickenSoup | *models.BooksList | *models.BookContent | *models.ExperienceCv | *models.ProjectCv | *models.SkillCv | *models.Diary | *models.AuthUser | *models.ProjectCvPs | *models.BlogPostPs | *models.DiaryPs | *models.LoveInfo | *models.Memo
	GetTableName() string
}

func AddDBData[V blogDataType](data V) error {
	db := global.App.DB
	result := db.Create(data)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetDBData[V blogDataType](tableType V, options ...func(option *gorm.DB)) ([]V, error) {
	db := global.App.DB

	db = db.Table(tableType.GetTableName())
	if db.Error != nil {
		return nil, db.Error
	}

	for _, option := range options {
		option(db)
	}

	var info []V

	result := db.Find(&info)
	if result.Error != nil {
		return nil, result.Error
	}

	return info, nil
}

func DeleteDBData[V blogDataType](data V, options ...func(option *gorm.DB)) error {
	db := global.App.DB
	db = db.Table(data.GetTableName())

	for _, option := range options {
		option(db)
	}

	result := db.Delete(data)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
