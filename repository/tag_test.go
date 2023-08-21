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

type TagRepositorylSuite struct {
	suite.Suite
}

func (suite *TagRepositorylSuite) FetchTagRow(tag model.BlogTag) []driver.Value {
	return []driver.Value{tag.ID.ID, tag.Name}
}

func (s *TagRepositorylSuite) TestCountTag() {
	model := NewTag()

	mockFunc := func(mock sqlmock.Sqlmock, arg FindTagArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT count(*) FROM `blog_tags` WHERE")

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

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
		mock.ExpectCommit()
	}

	testCases := []struct {
		name     string
		arg      FindTagArg
		expected int64
	}{
		{
			name: "normal",
			arg: FindTagArg{
				IDs:     []uint{1, 2, 3},
				Names:   []string{"tag1", "tag2", "tag3"},
				Offset:  0,
				Limit:   10,
				NoLimit: false,
			},
			expected: 1,
		},
		{
			name: "no limit",
			arg: FindTagArg{
				IDs:     []uint{1, 3},
				Names:   []string{"tag1", "tag2", "tag3"},
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
					return model.CountTags(context.Background(), tx, tc.arg)
				})
				s.NoError(err)
				s.Equal(tc.expected, count)

			})
		})
	}
}

func (s *TagRepositorylSuite) TestFindTag() {

	type TestArg struct {
		name    string
		In      FindTagArg
		Out     []model.BlogTag
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT * FROM `blog_tags` WHERE")

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

		if !arg.In.NoLimit {
			query.WriteString(" ORDER BY `id` DESC")
			query.WriteString(" LIMIT ? OFFSET ?")
			args = append(args, arg.In.Limit, arg.In.Offset)
		}

		rows := mock.NewRows([]string{"id", "name"})

		for _, tag := range arg.Out {
			rows.AddRow(s.FetchTagRow(tag)...)
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
			In: FindTagArg{
				IDs:     []uint{1, 2, 3},
				Names:   []string{"tag1", "tag2", "tag3"},
				Offset:  0,
				Limit:   10,
				NoLimit: false,
			},
			Out: []model.BlogTag{
				{
					ID:   model.ID{ID: 1},
					Name: "tag1",
				},
			},
		},
		{
			name: "no limit",
			In: FindTagArg{
				IDs:     []uint{1, 3},
				Names:   []string{"tag1", "tag2", "tag3"},
				Offset:  0,
				Limit:   10,
				NoLimit: true,
			},
			Out: []model.BlogTag{
				{
					ID:   model.ID{ID: 1},
					Name: "tag1",
				},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				result, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewTag()
					return model.FindTags(context.Background(), tx, tc.In)
				})
				s.NoError(err)

				tags, ok := result.([]model.BlogTag)
				s.True(ok)

				s.Equal(len(tc.Out), len(tags))
			})
		})
	}
}

func (s *TagRepositorylSuite) TestGetTag() {
	type TestArg struct {
		name    string
		In      uint
		Out     model.BlogTag
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT * FROM `blog_tags` WHERE")

		query.WriteString(" `id` = ?")
		args = append(args, arg.In)
		query.WriteString(" ORDER BY `blog_tags`.`id` LIMIT 1")

		rows := mock.NewRows([]string{"id", "name"})

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnRows(rows.AddRow(s.FetchTagRow(arg.Out)...))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			name: "normal",
			In:   1,
			Out: model.BlogTag{
				ID:   model.ID{ID: 1},
				Name: "tag1",
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				result, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewTag()
					return model.GetTag(context.Background(), tx, tc.In)
				})
				s.NoError(err)

				tags, ok := result.(model.BlogTag)
				s.True(ok)

				s.Equal(tc.Out, tags)
			})
		})
	}
}

func (s *TagRepositorylSuite) TestCreateTag() {
	type TestArg struct {
		name    string
		In      model.BlogTag
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("INSERT INTO `blog_tags`")

		query.WriteString(" (`name`) VALUES (?)")
		args = append(args, arg.In.Name)

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			name: "normal",
			In: model.BlogTag{
				Name: "tag1",
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				_, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewTag()
					return nil, model.CreateTag(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func (s *TagRepositorylSuite) TestUpdateTag() {
	type TestArg struct {
		name    string
		In      model.BlogTag
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("UPDATE `blog_tags` SET")

		query.WriteString(" `id`=?,`name`=? WHERE `id` = ?")
		args = append(args, arg.In.ID.ID, arg.In.Name, arg.In.ID.ID)

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			name: "normal",
			In: model.BlogTag{
				ID:   model.ID{ID: 1},
				Name: "tag1",
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				_, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewTag()
					return nil, model.UpdateTag(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func (s *TagRepositorylSuite) TestDeleteTag() {
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
		query.WriteString("DELETE FROM `blog_tags` WHERE")

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
					model := NewTag()
					return nil, model.DeleteTag(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func TestTagRepositorylSuite(t *testing.T) {
	suite.Run(t, new(TagRepositorylSuite))
}
