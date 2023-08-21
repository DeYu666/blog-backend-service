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

type ChickenSoupRepositorySuite struct {
	suite.Suite
}

func (suite *ChickenSoupRepositorySuite) FetchChickenSoupsRow(cate model.ChickenSoup) []driver.Value {
	return []driver.Value{cate.ID.ID, cate.Sentence}
}

func (s *ChickenSoupRepositorySuite) TestCountChickenSoups() {

	mockFunc := func(mock sqlmock.Sqlmock, arg FindChickenSoupArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT count(*) FROM `chicken_soups` WHERE")

		query.WriteString(" 1 = 1")
		if arg.IDs != nil {
			query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(arg.IDs), ",")))
			for _, id := range arg.IDs {
				args = append(args, id)
			}
		}

		if arg.Sentence != nil {
			query.WriteString(fmt.Sprintf(" AND `sentence` IN (%s)", RepeatWithSep("?", len(arg.Sentence), ",")))
			for _, s := range arg.Sentence {
				args = append(args, s)
			}
		}

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
		mock.ExpectCommit()
	}

	testCases := []struct {
		sentence string
		arg      FindChickenSoupArg
		expected int64
	}{
		{
			sentence: "normal",
			arg: FindChickenSoupArg{
				IDs:      []uint{1, 2, 3},
				Sentence: []string{"cate1", "cate2", "cate3"},
				Offset:   0,
				Limit:    10,
				NoLimit:  false,
			},
			expected: 1,
		},
		{
			sentence: "no limit",
			arg: FindChickenSoupArg{
				IDs:      []uint{1, 3},
				Sentence: []string{"cate1", "cate2", "cate3"},
				Offset:   0,
				Limit:    10,
				NoLimit:  true,
			},
			expected: 1,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.sentence, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc.arg)

				count, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewChickenSoup()
					return model.CountChickenSoups(context.Background(), tx, tc.arg)
				})
				s.NoError(err)
				s.Equal(tc.expected, count)

			})
		})
	}
}

func (s *ChickenSoupRepositorySuite) TestFindChickenSoups() {

	type TestArg struct {
		sentence string
		In       FindChickenSoupArg
		Out      []model.ChickenSoup
		OutErr   bool
		Err      error
		Prepare  func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc   func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT * FROM `chicken_soups` WHERE")

		query.WriteString(" 1 = 1")
		if arg.In.IDs != nil {
			query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(arg.In.IDs), ",")))
			for _, id := range arg.In.IDs {
				args = append(args, id)
			}
		}
		if arg.In.Sentence != nil {
			query.WriteString(fmt.Sprintf(" AND `sentence` IN (%s)", RepeatWithSep("?", len(arg.In.Sentence), ",")))
			for _, sentence := range arg.In.Sentence {
				args = append(args, sentence)
			}
		}

		if !arg.In.NoLimit {
			query.WriteString(" ORDER BY `id` DESC")
			query.WriteString(" LIMIT ? OFFSET ?")
			args = append(args, arg.In.Limit, arg.In.Offset)
		}

		rows := mock.NewRows([]string{"id", "sentence"})

		for _, soup := range arg.Out {
			rows.AddRow(s.FetchChickenSoupsRow(soup)...)
		}

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnRows(rows)
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			sentence: "normal",
			In: FindChickenSoupArg{
				IDs:      []uint{1, 2, 3},
				Sentence: []string{"cate1", "cate2", "cate3"},
				Offset:   0,
				Limit:    10,
				NoLimit:  false,
			},
			Out: []model.ChickenSoup{
				{
					ID:       model.ID{ID: 1},
					Sentence: "cate1",
				},
			},
		},
		{
			sentence: "no limit",
			In: FindChickenSoupArg{
				IDs:      []uint{1, 3},
				Sentence: []string{"cate1", "cate2", "cate3"},
				Offset:   0,
				Limit:    10,
				NoLimit:  true,
			},
			Out: []model.ChickenSoup{
				{
					ID:       model.ID{ID: 1},
					Sentence: "cate1",
				},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.sentence, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				result, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewChickenSoup()
					return model.FindChickenSoups(context.Background(), tx, tc.In)
				})
				s.NoError(err)

				cates, ok := result.([]model.ChickenSoup)
				s.True(ok)

				s.Equal(len(tc.Out), len(cates))
			})
		})
	}
}

func (s *ChickenSoupRepositorySuite) TestGetCategories() {
	type TestArg struct {
		sentence string
		In       uint
		Out      model.ChickenSoup
		OutErr   bool
		Err      error
		Prepare  func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc   func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT * FROM `chicken_soups` WHERE")

		query.WriteString(" `id` = ?")
		args = append(args, arg.In)

		query.WriteString(" ORDER BY `chicken_soups`.`id` LIMIT 1")

		rows := mock.NewRows([]string{"id", "sentence"})

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnRows(rows.AddRow(s.FetchChickenSoupsRow(arg.Out)...))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			sentence: "normal",
			In:       1,
			Out: model.ChickenSoup{
				ID:       model.ID{ID: 1},
				Sentence: "cate1",
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.sentence, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				result, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewChickenSoup()
					return model.GetChickenSoups(context.Background(), tx, tc.In)
				})
				s.NoError(err)

				cates, ok := result.(model.ChickenSoup)
				s.True(ok)

				s.Equal(tc.Out, cates)
			})
		})
	}
}

func (s *ChickenSoupRepositorySuite) TestCreateChickenSoups() {
	type TestArg struct {
		sentence string
		In       model.ChickenSoup
		OutErr   bool
		Err      error
		Prepare  func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc   func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("INSERT INTO `chicken_soups`")

		query.WriteString(" (`sentence`) VALUES (?)")
		args = append(args, arg.In.Sentence)

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			sentence: "normal",
			In: model.ChickenSoup{
				Sentence: "cate1",
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.sentence, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				_, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewChickenSoup()
					return nil, model.CreateChickenSoup(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func (s *ChickenSoupRepositorySuite) TestUpdateChickenSoups() {
	type TestArg struct {
		sentence string
		In       model.ChickenSoup
		OutErr   bool
		Err      error
		Prepare  func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc   func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("UPDATE `chicken_soups` SET")

		query.WriteString(" `id`=?,`sentence`=? WHERE `id` = ?")
		args = append(args, arg.In.ID.ID, arg.In.Sentence, arg.In.ID.ID)

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			sentence: "normal",
			In: model.ChickenSoup{
				ID:       model.ID{ID: 1},
				Sentence: "cate1",
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.sentence, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				_, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewChickenSoup()
					return nil, model.UpdateChickenSoup(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func (s *ChickenSoupRepositorySuite) TestDeleteChickenSoups() {
	type TestArg struct {
		sentence string
		In       uint
		OutErr   bool
		Err      error
		Prepare  func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc   func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("DELETE FROM `chicken_soups` WHERE")

		query.WriteString(" `id` = ?")
		args = append(args, arg.In)

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			sentence: "normal",
			In:       1,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.sentence, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				_, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewChickenSoup()
					return nil, model.DeleteChickenSoup(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func TestChickenSoupRepositorylSuite(t *testing.T) {
	suite.Run(t, new(ChickenSoupRepositorySuite))
}
