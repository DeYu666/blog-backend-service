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

type BookShelfRepositorylSuite struct {
	suite.Suite
}

func (suite *BookShelfRepositorylSuite) FetchBookShelfRow(bookShelf model.BooksList) []driver.Value {
	return []driver.Value{bookShelf.ID.ID, bookShelf.BookName, bookShelf.BookStatus, bookShelf.Abstract, bookShelf.Timestamps.CreatedTime, bookShelf.Timestamps.ModifiedTime}
}

func (s *BookShelfRepositorylSuite) TestCountBookShelf() {
	model := NewBookShelf()

	mockFunc := func(mock sqlmock.Sqlmock, arg FindBookShelfArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT count(*) FROM `books_lists` WHERE")

		query.WriteString(" 1 = 1")
		if arg.IDs != nil {
			query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(arg.IDs), ",")))
			for _, id := range arg.IDs {
				args = append(args, id)
			}
		}
		if arg.BookNameLikes != nil {
			query.WriteString(fmt.Sprintf(" AND `book_name` Like (%s)", RepeatWithSep("?", len(arg.BookNameLikes), ",")))
			for _, name := range arg.BookNameLikes {
				args = append(args, name)
			}
		}

		if arg.BookNames != nil {
			query.WriteString(fmt.Sprintf(" AND `book_name` IN (%s)", RepeatWithSep("?", len(arg.BookNames), ",")))
			for _, id := range arg.BookNames {
				args = append(args, id)
			}
		}

		if arg.BookStatuses != nil {
			query.WriteString(fmt.Sprintf(" AND `book_status` IN (%s)", RepeatWithSep("?", len(arg.BookStatuses), ",")))
			for _, id := range arg.BookStatuses {
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
		arg      FindBookShelfArg
		expected int64
	}{
		{
			name: "normal",
			arg: FindBookShelfArg{
				IDs:           []uint{1, 2, 3},
				BookNameLikes: []string{"bookShelf1", "bookShelf2", "bookShelf3"},
				Offset:        0,
				Limit:         10,
				NoLimit:       false,
			},
			expected: 1,
		},
		{
			name: "no limit",
			arg: FindBookShelfArg{
				IDs:           []uint{1, 3},
				BookNameLikes: []string{"bookShelf1", "bookShelf2", "bookShelf3"},
				Offset:        0,
				Limit:         10,
				NoLimit:       true,
			},
			expected: 1,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc.arg)

				count, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					return model.CountBookShelfs(context.Background(), tx, tc.arg)
				})
				s.NoError(err)
				s.Equal(tc.expected, count)

			})
		})
	}
}

func (s *BookShelfRepositorylSuite) TestFindBookShelf() {

	type TestArg struct {
		name    string
		In      FindBookShelfArg
		Out     []model.BooksList
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT * FROM `books_lists` WHERE")

		query.WriteString(" 1 = 1")
		if arg.In.IDs != nil {
			query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(arg.In.IDs), ",")))
			for _, id := range arg.In.IDs {
				args = append(args, id)
			}
		}
		if arg.In.BookNameLikes != nil {
			query.WriteString(fmt.Sprintf(" AND `book_name` Like (%s)", RepeatWithSep("?", len(arg.In.BookNameLikes), ",")))
			for _, name := range arg.In.BookNameLikes {
				args = append(args, name)
			}
		}

		if arg.In.BookStatuses != nil {
			query.WriteString(fmt.Sprintf(" AND `book_status` IN (%s)", RepeatWithSep("?", len(arg.In.BookStatuses), ",")))
			for _, id := range arg.In.BookStatuses {
				args = append(args, id)
			}
		}

		if !arg.In.NoLimit {
			query.WriteString(" ORDER BY `id` DESC")
			query.WriteString(" LIMIT ? OFFSET ?")
			args = append(args, arg.In.Limit, arg.In.Offset)
		}

		rows := mock.NewRows([]string{"id", "book_name", "book_status", "abstract", "created_time", "modified_time"})

		for _, bookShelf := range arg.Out {
			rows.AddRow(s.FetchBookShelfRow(bookShelf)...)
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
			In: FindBookShelfArg{
				IDs:           []uint{1, 2, 3},
				BookNameLikes: []string{"bookShelf1", "bookShelf2", "bookShelf3"},
				Offset:        0,
				Limit:         10,
				NoLimit:       false,
			},
			Out: []model.BooksList{
				{
					ID:         model.ID{ID: 1},
					BookName:   "bookShelf1",
					BookStatus: "bookStatus1",
				},
			},
		},
		{
			name: "no limit",
			In: FindBookShelfArg{
				IDs:           []uint{1, 3},
				BookNameLikes: []string{"bookShelf1", "bookShelf2", "bookShelf3"},
				Offset:        0,
				Limit:         10,
				NoLimit:       true,
			},
			Out: []model.BooksList{
				{
					ID:         model.ID{ID: 1},
					BookName:   "bookShelf1",
					BookStatus: "bookStatus1",
				},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				result, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewBookShelf()
					return model.FindBookShelfs(context.Background(), tx, tc.In)
				})
				s.NoError(err)

				bookShelfs, ok := result.([]model.BooksList)
				s.True(ok)

				s.Equal(len(tc.Out), len(bookShelfs))
			})
		})
	}
}

func (s *BookShelfRepositorylSuite) TestGetBookShelf() {
	type TestArg struct {
		name    string
		In      uint
		Out     model.BooksList
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT * FROM `books_lists` WHERE")

		query.WriteString(" `id` = ?")
		args = append(args, arg.In)
		query.WriteString(" ORDER BY `books_lists`.`id` LIMIT 1")

		rows := mock.NewRows([]string{"id", "book_name", "book_status", "abstract", "created_time", "modified_time"})

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnRows(rows.AddRow(s.FetchBookShelfRow(arg.Out)...))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			name: "normal",
			In:   1,
			Out: model.BooksList{
				ID:         model.ID{ID: 1},
				BookName:   "bookShelf1",
				BookStatus: "bookStatus1",
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				result, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewBookShelf()
					return model.GetBookShelf(context.Background(), tx, tc.In)
				})
				s.NoError(err)

				bookShelfs, ok := result.(model.BooksList)
				s.True(ok)

				s.Equal(tc.Out, bookShelfs)
			})
		})
	}
}

func (s *BookShelfRepositorylSuite) TestCreateBookShelf() {
	type TestArg struct {
		name    string
		In      model.BooksList
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("INSERT INTO `books_lists`")

		query.WriteString(" (`book_name`,`book_status`,`abstract`,`created_time`,`modified_time`) VALUES (?,?,?,?,?)")
		args = append(args, arg.In.BookName, arg.In.BookStatus, arg.In.Abstract, arg.In.CreatedTime, arg.In.ModifiedTime)

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			name: "normal",
			In: model.BooksList{
				BookName:   "bookShelf1",
				BookStatus: "bookStatus1",
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				_, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewBookShelf()
					return nil, model.CreateBookShelf(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func (s *BookShelfRepositorylSuite) TestUpdateBookShelf() {
	type TestArg struct {
		name    string
		In      model.BooksList
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("UPDATE `books_lists` SET")

		query.WriteString(" `id`=?,`book_name`=?,`book_status`=? WHERE `id` = ?")
		args = append(args, arg.In.ID.ID, arg.In.BookName, arg.In.BookStatus, arg.In.ID.ID)

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			name: "normal",
			In: model.BooksList{
				ID:         model.ID{ID: 1},
				BookName:   "bookShelf1",
				BookStatus: "bookStatus1",
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				_, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewBookShelf()
					return nil, model.UpdateBookShelf(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func (s *BookShelfRepositorylSuite) TestDeleteBookShelf() {
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
		query.WriteString("DELETE FROM `books_lists` WHERE")

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
					model := NewBookShelf()
					return nil, model.DeleteBookShelf(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func TestBookShelfRepositorylSuite(t *testing.T) {
	suite.Run(t, new(BookShelfRepositorylSuite))
}
