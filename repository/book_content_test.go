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

type BookContentRepositorylSuite struct {
	suite.Suite
}

func (suite *BookContentRepositorylSuite) FetchBookContentRow(bookContent model.BookContent) []driver.Value {
	return []driver.Value{bookContent.ID.ID, bookContent.BookContent, bookContent.BookId}
}

func (s *BookContentRepositorylSuite) TestCountBookContent() {
	model := NewBookContent()

	mockFunc := func(mock sqlmock.Sqlmock, arg FindBookContentArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT count(*) FROM `book_contents` WHERE")

		query.WriteString(" 1 = 1")
		if arg.IDs != nil {
			query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(arg.IDs), ",")))
			for _, id := range arg.IDs {
				args = append(args, id)
			}
		}
		if arg.ConentLikes != nil {
			query.WriteString(fmt.Sprintf(" AND `book_content` Like (%s)", RepeatWithSep("?", len(arg.ConentLikes), ",")))
			for _, name := range arg.ConentLikes {
				args = append(args, name)
			}
		}

		if arg.BookIDs != nil {
			query.WriteString(fmt.Sprintf(" AND `book_id` IN (%s)", RepeatWithSep("?", len(arg.BookIDs), ",")))
			for _, id := range arg.BookIDs {
				args = append(args, id)
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
		arg      FindBookContentArg
		expected int64
	}{
		{
			name: "normal",
			arg: FindBookContentArg{
				IDs:         []uint{1, 2, 3},
				ConentLikes: []string{"bookContent1", "bookContent2", "bookContent3"},
				Offset:      0,
				Limit:       10,
				NoLimit:     false,
			},
			expected: 1,
		},
		{
			name: "no limit",
			arg: FindBookContentArg{
				IDs:         []uint{1, 3},
				ConentLikes: []string{"bookContent1", "bookContent2", "bookContent3"},
				Offset:      0,
				Limit:       10,
				NoLimit:     true,
			},
			expected: 1,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc.arg)

				count, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					return model.CountBookContents(context.Background(), tx, tc.arg)
				})
				s.NoError(err)
				s.Equal(tc.expected, count)

			})
		})
	}
}

func (s *BookContentRepositorylSuite) TestFindBookContent() {

	type TestArg struct {
		name    string
		In      FindBookContentArg
		Out     []model.BookContent
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT * FROM `book_contents` WHERE")

		query.WriteString(" 1 = 1")
		if arg.In.IDs != nil {
			query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(arg.In.IDs), ",")))
			for _, id := range arg.In.IDs {
				args = append(args, id)
			}
		}
		if arg.In.ConentLikes != nil {
			query.WriteString(fmt.Sprintf(" AND `book_content` Like (%s)", RepeatWithSep("?", len(arg.In.ConentLikes), ",")))
			for _, name := range arg.In.ConentLikes {
				args = append(args, name)
			}
		}

		if !arg.In.NoLimit {
			query.WriteString(" ORDER BY `id` DESC")
			query.WriteString(" LIMIT ? OFFSET ?")
			args = append(args, arg.In.Limit, arg.In.Offset)
		}

		rows := mock.NewRows([]string{"id", "book_content", "book_id"})

		for _, bookContent := range arg.Out {
			rows.AddRow(s.FetchBookContentRow(bookContent)...)
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
			In: FindBookContentArg{
				IDs:         []uint{1, 2, 3},
				ConentLikes: []string{"bookContent1", "bookContent2", "bookContent3"},
				Offset:      0,
				Limit:       10,
				NoLimit:     false,
			},
			Out: []model.BookContent{
				{
					ID:          model.ID{ID: 1},
					BookContent: "bookContent1",
				},
			},
		},
		{
			name: "no limit",
			In: FindBookContentArg{
				IDs:         []uint{1, 3},
				ConentLikes: []string{"bookContent1", "bookContent2", "bookContent3"},
				Offset:      0,
				Limit:       10,
				NoLimit:     true,
			},
			Out: []model.BookContent{
				{
					ID:          model.ID{ID: 1},
					BookContent: "bookContent1",
				},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				result, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewBookContent()
					return model.FindBookContents(context.Background(), tx, tc.In)
				})
				s.NoError(err)

				bookContents, ok := result.([]model.BookContent)
				s.True(ok)

				s.Equal(len(tc.Out), len(bookContents))
			})
		})
	}
}

func (s *BookContentRepositorylSuite) TestGetBookContent() {
	type TestArg struct {
		name    string
		In      uint
		Out     model.BookContent
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT * FROM `book_contents` WHERE")

		query.WriteString(" `id` = ?")
		args = append(args, arg.In)
		query.WriteString(" ORDER BY `book_contents`.`id` LIMIT 1")

		rows := mock.NewRows([]string{"id", "book_content", "book_id"})

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnRows(rows.AddRow(s.FetchBookContentRow(arg.Out)...))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			name: "normal",
			In:   1,
			Out: model.BookContent{
				ID:          model.ID{ID: 1},
				BookContent: "bookContent1",
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				result, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewBookContent()
					return model.GetBookContent(context.Background(), tx, tc.In)
				})
				s.NoError(err)

				bookContents, ok := result.(model.BookContent)
				s.True(ok)

				s.Equal(tc.Out, bookContents)
			})
		})
	}
}

func (s *BookContentRepositorylSuite) TestCreateBookContent() {
	type TestArg struct {
		name    string
		In      model.BookContent
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("INSERT INTO `book_contents`")

		query.WriteString(" (`book_content`,`book_id`,`created_time`,`modified_time`) VALUES (?,?,?,?)")
		args = append(args, arg.In.BookContent, arg.In.BookId, arg.In.CreatedTime, arg.In.ModifiedTime)

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			name: "normal",
			In: model.BookContent{
				BookContent: "bookContent1",
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				_, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewBookContent()
					return nil, model.CreateBookContent(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func (s *BookContentRepositorylSuite) TestUpdateBookContent() {
	type TestArg struct {
		name    string
		In      model.BookContent
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("UPDATE `book_contents` SET")

		query.WriteString(" `id`=?,`book_content`=? WHERE `id` = ?")
		args = append(args, arg.In.ID.ID, arg.In.BookContent, arg.In.ID.ID)

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			name: "normal",
			In: model.BookContent{
				ID:          model.ID{ID: 1},
				BookContent: "bookContent1",
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				_, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewBookContent()
					return nil, model.UpdateBookContent(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func (s *BookContentRepositorylSuite) TestDeleteBookContent() {
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
		query.WriteString("DELETE FROM `book_contents` WHERE")

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
			name: "normal",
			In:   1,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				_, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewBookContent()
					return nil, model.DeleteBookContent(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func TestBookContentRepositorylSuite(t *testing.T) {
	suite.Run(t, new(BookContentRepositorylSuite))
}
