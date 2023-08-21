package repository

import (
	"context"
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DeYu666/blog-backend-service/lib/test"
	"github.com/DeYu666/blog-backend-service/model"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type DiaryRepositorylSuite struct {
	suite.Suite
}

func (suite *DiaryRepositorylSuite) FetchDiaryRow(blog model.Diary) []driver.Value {
	return []driver.Value{blog.ID.ID, blog.Content, blog.IsOpen, blog.Timestamps.CreatedTime, blog.Timestamps.ModifiedTime}
}

func (s *DiaryRepositorylSuite) TestCountDiary() {

	mockFunc := func(mock sqlmock.Sqlmock, arg FindDiaryArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT count(*) FROM `diaries` WHERE")

		query.WriteString(" 1 = 1")
		if arg.IDs != nil {
			query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(arg.IDs), ",")))
			for _, id := range arg.IDs {
				args = append(args, id)
			}
		}

		if arg.ContentLikes != nil {
			query.WriteString(fmt.Sprintf(" AND `content` Like (%s)", RepeatWithSep("?", len(arg.ContentLikes), ",")))
			for _, content := range arg.ContentLikes {
				args = append(args, content)
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
		arg      FindDiaryArg
		expected int64
	}{
		{
			name: "normal",
			arg: FindDiaryArg{
				IDs:          []uint{1, 2, 3},
				ContentLikes: []string{"blog1", "blog2", "blog3"},
				Offset:       0,
				Limit:        10,
				NoLimit:      false,
			},
			expected: 1,
		},
		{
			name: "no limit",
			arg: FindDiaryArg{
				IDs:          []uint{1, 2, 3},
				ContentLikes: []string{"blog1", "blog2", "blog3"},
				Offset:       0,
				Limit:        10,
				NoLimit:      true,
			},
			expected: 1,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc.arg)

				count, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewDiary()
					return model.CountDiarys(context.Background(), tx, tc.arg)
				})
				s.NoError(err)
				s.Equal(tc.expected, count)

			})
		})
	}
}

func (s *DiaryRepositorylSuite) TestFindDiary() {

	type TestArg struct {
		name    string
		In      FindDiaryArg
		Out     []model.Diary
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT * FROM `diaries` WHERE")

		query.WriteString(" 1 = 1")
		if arg.In.IDs != nil {
			query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(arg.In.IDs), ",")))
			for _, id := range arg.In.IDs {
				args = append(args, id)
			}
		}

		if arg.In.ContentLikes != nil {
			query.WriteString(fmt.Sprintf(" AND `content` Like (%s)", RepeatWithSep("?", len(arg.In.ContentLikes), ",")))
			for _, content := range arg.In.ContentLikes {
				args = append(args, content)
			}
		}

		if !arg.In.NoLimit {
			query.WriteString(" ORDER BY `id` DESC")
			query.WriteString(" LIMIT ? OFFSET ?")
			args = append(args, arg.In.Limit, arg.In.Offset)
		}
		rows := mock.NewRows([]string{"id", "content", "is_open", "created_time", "modified_time"})

		for _, blog := range arg.Out {
			rows.AddRow(s.FetchDiaryRow(blog)...)
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
			In: FindDiaryArg{
				IDs:          []uint{1, 2, 3},
				ContentLikes: []string{"blog1", "blog2", "blog3"},
				Offset:       0,
				Limit:        10,
				NoLimit:      false,
			},
			Out: []model.Diary{
				{
					ID:      model.ID{ID: 1},
					Content: "blog1",
					IsOpen:  true,
					Timestamps: model.Timestamps{
						CreatedTime:  time.Now(),
						ModifiedTime: time.Now(),
					},
				},
			},
		},
		{
			name: "no limit",
			In: FindDiaryArg{
				IDs:          []uint{1, 2, 3},
				ContentLikes: []string{"blog1", "blog2", "blog3"},
				Offset:       0,
				Limit:        10,
				NoLimit:      true,
			},
			Out: []model.Diary{
				{
					ID:      model.ID{ID: 1},
					Content: "blog1",
					IsOpen:  true,
					Timestamps: model.Timestamps{
						CreatedTime:  time.Now(),
						ModifiedTime: time.Now(),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				result, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewDiary()
					return model.FindDiarys(context.Background(), tx, tc.In)
				})
				s.NoError(err)

				blogs, ok := result.([]model.Diary)
				s.True(ok)

				s.Equal(len(tc.Out), len(blogs))
			})
		})
	}
}

func (s *DiaryRepositorylSuite) TestGetCategories() {
	type TestArg struct {
		name    string
		In      uint
		Out     model.Diary
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT * FROM `diaries` WHERE")

		query.WriteString(" `id` = ?")
		args = append(args, arg.In)

		query.WriteString(" ORDER BY `diaries`.`id` LIMIT 1")

		rows := mock.NewRows([]string{"id", "content", "is_open", "created_time", "modified_time"})

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnRows(rows.AddRow(s.FetchDiaryRow(arg.Out)...))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			name: "normal",
			In:   1,
			Out: model.Diary{
				ID:      model.ID{ID: 1},
				Content: "blog1",
				IsOpen:  true,
				Timestamps: model.Timestamps{
					CreatedTime:  time.Now(),
					ModifiedTime: time.Now(),
				},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				result, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewDiary()
					return model.GetDiary(context.Background(), tx, tc.In)
				})
				s.NoError(err)

				blogs, ok := result.(model.Diary)
				s.True(ok)

				s.Equal(tc.Out, blogs)
			})
		})
	}
}

func (s *DiaryRepositorylSuite) TestCreateDiary() {
	type TestArg struct {
		name    string
		In      model.Diary
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("INSERT INTO `diaries`")

		query.WriteString(" (`content`,`is_open`,`created_time`,`modified_time`,`id`) VALUES (?,?,?,?,?)")
		args = append(args, arg.In.Content, arg.In.IsOpen, arg.In.Timestamps.CreatedTime, arg.In.Timestamps.ModifiedTime, arg.In.ID.ID)

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			name: "normal",
			In: model.Diary{
				ID:      model.ID{ID: 1},
				Content: "blog1",
				IsOpen:  true,
				Timestamps: model.Timestamps{
					CreatedTime:  time.Now(),
					ModifiedTime: time.Now(),
				},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				_, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewDiary()
					return nil, model.CreateDiary(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func (s *DiaryRepositorylSuite) TestUpdateDiary() {
	type TestArg struct {
		name    string
		In      model.Diary
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("UPDATE `diaries` SET")

		query.WriteString(" `id`=?,`content`=?,`is_open`=?,`created_time`=?,`modified_time`=? WHERE `id` = ?")
		args = append(args, arg.In.ID.ID, arg.In.Content, arg.In.IsOpen, arg.In.Timestamps.CreatedTime, arg.In.Timestamps.ModifiedTime, arg.In.ID.ID)

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			name: "normal",
			In: model.Diary{
				ID:      model.ID{ID: 1},
				Content: "blog1",
				IsOpen:  true,
				Timestamps: model.Timestamps{
					CreatedTime:  time.Now(),
					ModifiedTime: time.Now(),
				},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				_, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewDiary()
					return nil, model.UpdateDiary(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func (s *DiaryRepositorylSuite) TestDeleteDiary() {
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
		query.WriteString("DELETE FROM `diaries` WHERE")

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
					model := NewDiary()
					return nil, model.DeleteDiary(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func TestDiaryRepositorylSuite(t *testing.T) {
	suite.Run(t, new(DiaryRepositorylSuite))
}
