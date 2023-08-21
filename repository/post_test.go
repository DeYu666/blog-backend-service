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

type PostRepositorylSuite struct {
	suite.Suite
}

func (suite *PostRepositorylSuite) FetchPostRow(blog model.BlogPost) []driver.Value {
	return []driver.Value{blog.ID.ID, blog.Title, blog.Body, blog.CategoryID, blog.Excerpt, blog.Timestamps.CreatedTime, blog.Timestamps.ModifiedTime}
}

func (s *PostRepositorylSuite) TestCountPost() {

	mockFunc := func(mock sqlmock.Sqlmock, arg FindPostArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT count(*) FROM `blog_posts` WHERE")

		query.WriteString(" 1 = 1")
		if arg.IDs != nil {
			query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(arg.IDs), ",")))
			for _, id := range arg.IDs {
				args = append(args, id)
			}
		}

		if arg.Titles != nil {
			query.WriteString(fmt.Sprintf(" AND `title` IN (%s)", RepeatWithSep("?", len(arg.Titles), ",")))
			for _, title := range arg.Titles {
				args = append(args, title)
			}
		}

		if arg.CategoryIDs != nil {
			query.WriteString(fmt.Sprintf(" AND `category_id` IN (%s)", RepeatWithSep("?", len(arg.CategoryIDs), ",")))
			for _, categoryID := range arg.CategoryIDs {
				args = append(args, categoryID)
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
		arg      FindPostArg
		expected int64
	}{
		{
			name: "normal",
			arg: FindPostArg{
				IDs:         []uint{1, 2, 3},
				Titles:      []string{"blog1", "blog2", "blog3"},
				CategoryIDs: []uint{1, 2, 3},
				Offset:      0,
				Limit:       10,
				NoLimit:     false,
			},
			expected: 1,
		},
		{
			name: "no limit",
			arg: FindPostArg{
				IDs:         []uint{1, 2, 3},
				Titles:      []string{"blog1", "blog2", "blog3"},
				CategoryIDs: []uint{1, 2, 3},
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
					model := NewPost()
					return model.CountPosts(context.Background(), tx, tc.arg)
				})
				s.NoError(err)
				s.Equal(tc.expected, count)

			})
		})
	}
}

func (s *PostRepositorylSuite) TestFindPost() {

	type TestArg struct {
		name    string
		In      FindPostArg
		Out     []model.BlogPost
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT * FROM `blog_posts` WHERE")

		query.WriteString(" 1 = 1")
		if arg.In.IDs != nil {
			query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(arg.In.IDs), ",")))
			for _, id := range arg.In.IDs {
				args = append(args, id)
			}
		}

		if arg.In.Titles != nil {
			query.WriteString(fmt.Sprintf(" AND `title` IN (%s)", RepeatWithSep("?", len(arg.In.Titles), ",")))
			for _, title := range arg.In.Titles {
				args = append(args, title)
			}
		}

		if arg.In.CategoryIDs != nil {
			query.WriteString(fmt.Sprintf(" AND `category_id` IN (%s)", RepeatWithSep("?", len(arg.In.CategoryIDs), ",")))
			for _, categoryID := range arg.In.CategoryIDs {
				args = append(args, categoryID)
			}
		}

		if !arg.In.NoLimit {
			query.WriteString(" ORDER BY `id` DESC")
			query.WriteString(" LIMIT ? OFFSET ?")
			args = append(args, arg.In.Limit, arg.In.Offset)
		}
		rows := mock.NewRows([]string{"id", "title", "body", "category_id", "excerpt", "created_time", "modified_time"})

		for _, blog := range arg.Out {
			rows.AddRow(s.FetchPostRow(blog)...)
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
			In: FindPostArg{
				IDs:     []uint{1, 2, 3},
				Titles:  []string{"blog1", "blog2", "blog3"},
				Offset:  0,
				Limit:   10,
				NoLimit: false,
			},
			Out: []model.BlogPost{
				{
					ID:         model.ID{ID: 1},
					Title:      "blog1",
					Body:       "body1",
					CategoryID: 1,
					Excerpt:    "excerpt1",
					AuthorID:   1,
				},
			},
		},
		{
			name: "no limit",
			In: FindPostArg{
				IDs:     []uint{1, 2, 3},
				Titles:  []string{"blog1", "blog2", "blog3"},
				Offset:  0,
				Limit:   10,
				NoLimit: true,
			},
			Out: []model.BlogPost{
				{
					ID:         model.ID{ID: 1},
					Title:      "blog1",
					Body:       "body1",
					CategoryID: 1,
					Excerpt:    "excerpt1",
					AuthorID:   1,
				},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				result, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewPost()
					return model.FindPosts(context.Background(), tx, tc.In)
				})
				s.NoError(err)

				blogs, ok := result.([]model.BlogPost)
				s.True(ok)

				s.Equal(len(tc.Out), len(blogs))
			})
		})
	}
}

func (s *PostRepositorylSuite) TestGetCategories() {
	type TestArg struct {
		name    string
		In      uint
		Out     model.BlogPost
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT * FROM `blog_posts` WHERE")

		query.WriteString(" `id` = ?")
		args = append(args, arg.In)

		query.WriteString(" ORDER BY `blog_posts`.`id` LIMIT 1")

		rows := mock.NewRows([]string{"id", "title", "body", "category_id", "excerpt", "created_time", "modified_time"})

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnRows(rows.AddRow(s.FetchPostRow(arg.Out)...))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			name: "normal",
			In:   1,
			Out: model.BlogPost{
				ID:         model.ID{ID: 1},
				Title:      "blog1",
				Body:       "body1",
				CategoryID: 1,
				Excerpt:    "excerpt1",
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				result, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewPost()
					return model.GetPost(context.Background(), tx, tc.In)
				})
				s.NoError(err)

				blogs, ok := result.(model.BlogPost)
				s.True(ok)

				s.Equal(tc.Out, blogs)
			})
		})
	}
}

func (s *PostRepositorylSuite) TestCreatePost() {
	type TestArg struct {
		name    string
		In      model.BlogPost
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("INSERT INTO `blog_posts`")

		query.WriteString(" (`created_time`,`modified_time`,`title`,`body`,`excerpt`,`category_id`,`author_id`,`views`,`likes`,`cover_url`,`title_url`,`is_open`,`id`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)")
		args = append(args, arg.In.Timestamps.CreatedTime, arg.In.Timestamps.ModifiedTime, arg.In.Title, arg.In.Body, arg.In.Excerpt, arg.In.CategoryID, arg.In.AuthorID, arg.In.Views, arg.In.Likes, arg.In.CoverURL, arg.In.TitleURL, arg.In.IsOpen, arg.In.ID.ID)

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			name: "normal",
			In: model.BlogPost{
				ID:         model.ID{ID: 1},
				Title:      "blog1",
				Body:       "body1",
				CategoryID: 1,
				Excerpt:    "excerpt1",
				AuthorID:   1,
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				_, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewPost()
					return nil, model.CreatePost(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func (s *PostRepositorylSuite) TestUpdatePost() {
	type TestArg struct {
		name    string
		In      model.BlogPost
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("UPDATE `blog_posts` SET")

		query.WriteString(" `id`=?,`title`=?,`body`=?,`excerpt`=?,`category_id`=?,`author_id`=? WHERE `id` = ?")
		args = append(args, arg.In.ID.ID, arg.In.Title, arg.In.Body, arg.In.Excerpt, arg.In.CategoryID, arg.In.AuthorID, arg.In.ID.ID)

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			name: "normal",
			In: model.BlogPost{
				ID:         model.ID{ID: 1},
				Title:      "blog1",
				Body:       "body1",
				CategoryID: 1,
				Excerpt:    "excerpt1",
				AuthorID:   1,
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				_, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewPost()
					return nil, model.UpdatePost(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func (s *PostRepositorylSuite) TestDeletePost() {
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
		query.WriteString("DELETE FROM `blog_posts` WHERE")

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
					model := NewPost()
					return nil, model.DeletePost(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func TestPostRepositorylSuite(t *testing.T) {
	suite.Run(t, new(PostRepositorylSuite))
}
