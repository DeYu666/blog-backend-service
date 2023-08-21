package repository

import (
	"context"
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DeYu666/blog-backend-service/lib/test"
	"github.com/DeYu666/blog-backend-service/model"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type CategoryRepositorylSuite struct {
	suite.Suite
}

func (suite *CategoryRepositorylSuite) FetchCategoryRow(cate model.BlogCategories) []driver.Value {
	return []driver.Value{cate.ID.ID, cate.Name, cate.GeneralID}
}

func (s *CategoryRepositorylSuite) TestCountCategory() {
	model := NewCategory()

	mockFunc := func(mock sqlmock.Sqlmock, arg FindCategoriesArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT count(*) FROM `blog_categories` WHERE")

		query.WriteString(" 1 = 1")
		if arg.IDs != nil {
			query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(arg.IDs), ",")))
			for _, id := range arg.IDs {
				args = append(args, id)
			}
		}
		if arg.Names != nil {
			query.WriteString(fmt.Sprintf(" AND `name` IN (%s)", RepeatWithSep("?", len(arg.Names), ",")))
			for _, name := range arg.Names {
				args = append(args, name)
			}
		}
		if arg.GeneralCateIds != nil {
			query.WriteString(fmt.Sprintf(" AND `general_id` IN (%s)", RepeatWithSep("?", len(arg.GeneralCateIds), ",")))
			for _, generalCateId := range arg.GeneralCateIds {
				args = append(args, generalCateId)
			}
		}

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
		mock.ExpectCommit()
	}

	testCases := []struct {
		name     string
		arg      FindCategoriesArg
		expected int64
	}{
		{
			name: "normal",
			arg: FindCategoriesArg{
				IDs:     []uint{1, 2, 3},
				Names:   []string{"cate1", "cate2", "cate3"},
				Offset:  0,
				Limit:   10,
				NoLimit: false,
			},
			expected: 1,
		},
		{
			name: "no limit",
			arg: FindCategoriesArg{
				IDs:     []uint{1, 3},
				Names:   []string{"cate1", "cate2", "cate3"},
				Offset:  0,
				Limit:   10,
				NoLimit: true,
			},
			expected: 1,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc.arg)

				count, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					return model.CountCategory(context.Background(), tx, tc.arg)
				})
				s.NoError(err)
				s.Equal(tc.expected, count)

			})
		})
	}
}

func (s *CategoryRepositorylSuite) TestFindCategory() {

	type TestArg struct {
		name    string
		In      FindCategoriesArg
		Out     []model.BlogCategories
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT * FROM `blog_categories` WHERE")

		query.WriteString(" 1 = 1")
		if arg.In.IDs != nil {
			query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(arg.In.IDs), ",")))
			for _, id := range arg.In.IDs {
				args = append(args, id)
			}
		}
		if arg.In.Names != nil {
			query.WriteString(fmt.Sprintf(" AND `name` IN (%s)", RepeatWithSep("?", len(arg.In.Names), ",")))
			for _, name := range arg.In.Names {
				args = append(args, name)
			}
		}
		if arg.In.GeneralCateIds != nil {
			query.WriteString(fmt.Sprintf(" AND `general_id` IN (%s)", RepeatWithSep("?", len(arg.In.GeneralCateIds), ",")))
			for _, generalCateId := range arg.In.GeneralCateIds {
				args = append(args, generalCateId)
			}
		}

		if !arg.In.NoLimit {
			query.WriteString(" ORDER BY `id` DESC")
			query.WriteString(" LIMIT ? OFFSET ?")
			args = append(args, arg.In.Limit, arg.In.Offset)
		}

		rows := mock.NewRows([]string{"id", "name", "general_id"})

		for _, cate := range arg.Out {
			rows.AddRow(s.FetchCategoryRow(cate)...)
		}

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnRows(rows)
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			name: "normal",
			In: FindCategoriesArg{
				IDs:     []uint{1, 2, 3},
				Names:   []string{"cate1", "cate2", "cate3"},
				Offset:  0,
				Limit:   10,
				NoLimit: false,
			},
			Out: []model.BlogCategories{
				{
					ID:        model.ID{ID: 1},
					Name:      "cate1",
					GeneralID: 1,
					General:   model.BlogGeneralCategories{},
				},
			},
		},
		{
			name: "no limit",
			In: FindCategoriesArg{
				IDs:     []uint{1, 3},
				Names:   []string{"cate1", "cate2", "cate3"},
				Offset:  0,
				Limit:   10,
				NoLimit: true,
			},
			Out: []model.BlogCategories{
				{
					ID:        model.ID{ID: 1},
					Name:      "cate1",
					GeneralID: 1,
					General:   model.BlogGeneralCategories{},
				},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				result, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewCategory()
					return model.FindCategory(context.Background(), tx, tc.In)
				})
				s.NoError(err)

				cates, ok := result.([]model.BlogCategories)
				s.True(ok)

				s.Equal(len(tc.Out), len(cates))
			})
		})
	}
}

func (s *CategoryRepositorylSuite) TestGetCategories() {
	type TestArg struct {
		name    string
		In      uint
		Out     model.BlogCategories
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT * FROM `blog_categories` WHERE")

		query.WriteString(" id = ?")
		args = append(args, arg.In)

		rows := mock.NewRows([]string{"id", "name", "general_id"})

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnRows(rows.AddRow(s.FetchCategoryRow(arg.Out)...))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			name: "normal",
			In:   1,
			Out: model.BlogCategories{
				ID:        model.ID{ID: 1},
				Name:      "cate1",
				GeneralID: 1,
				General:   model.BlogGeneralCategories{},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				result, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewCategory()
					return model.GetCategories(context.Background(), tx, tc.In)
				})
				s.NoError(err)

				cates, ok := result.(model.BlogCategories)
				s.True(ok)

				s.Equal(tc.Out, cates)
			})
		})
	}
}

func (s *CategoryRepositorylSuite) TestCreateCategory() {
	type TestArg struct {
		name    string
		In      model.BlogCategories
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("INSERT INTO `blog_categories`")

		query.WriteString(" (`name`,`general_id`) VALUES (?,?)")
		args = append(args, arg.In.Name, arg.In.GeneralID)

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			name: "normal",
			In: model.BlogCategories{
				Name:      "cate1",
				GeneralID: 1,
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				_, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewCategory()
					return nil, model.CreateCategory(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func (s *CategoryRepositorylSuite) TestUpdateCategory() {
	type TestArg struct {
		name    string
		In      model.BlogCategories
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("UPDATE `blog_categories` SET")

		query.WriteString(" `name`=?,`general_id`=? WHERE id = ? AND `id` = ?")
		args = append(args, arg.In.Name, arg.In.GeneralID, arg.In.ID.ID, arg.In.ID.ID)

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			name: "normal",
			In: model.BlogCategories{
				ID:        model.ID{ID: 1},
				Name:      "cate1",
				GeneralID: 1,
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				_, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewCategory()
					return nil, model.UpdateCategory(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func (s *CategoryRepositorylSuite) TestDeleteCategory() {
	type TestArg struct {
		name    string
		In      uint
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("DELETE FROM `blog_categories` WHERE")

		query.WriteString(" id = ?")
		args = append(args, arg.In)

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			name: "normal",
			In:   1,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				_, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewCategory()
					return nil, model.DeleteCategory(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func TestCategoryRepositorylSuite(t *testing.T) {
	suite.Run(t, new(CategoryRepositorylSuite))
}
