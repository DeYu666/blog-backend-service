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

type MemoRepositorylSuite struct {
	suite.Suite
}

func (suite *MemoRepositorylSuite) FetchMemoRow(tag model.Memo) []driver.Value {
	return []driver.Value{tag.ID.ID, tag.Content, tag.Status, tag.Timestamps.CreatedTime, tag.Timestamps.ModifiedTime}
}

func (s *MemoRepositorylSuite) TestCountMemo() {
	model := NewMemo()

	mockFunc := func(mock sqlmock.Sqlmock, arg FindMemoArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT count(*) FROM `memos` WHERE")

		query.WriteString(" 1 = 1")
		if arg.IDs != nil {
			query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(arg.IDs), ",")))
			for _, id := range arg.IDs {
				args = append(args, id)
			}
		}
		if arg.ContentLikes != nil {
			query.WriteString(fmt.Sprintf(" AND `content` Like (%s)", RepeatWithSep("?", len(arg.ContentLikes), ",")))
			for _, name := range arg.ContentLikes {
				args = append(args, name)
			}
		}

		if arg.StatusLikes != 0 {
			query.WriteString(" AND `status` = ?")
			args = append(args, arg.StatusLikes)
		}

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
		mock.ExpectCommit()
	}

	testCases := []struct {
		name     string
		arg      FindMemoArg
		expected int64
	}{
		{
			name: "normal",
			arg: FindMemoArg{
				IDs:          []uint{1, 2, 3},
				ContentLikes: []string{"tag1", "tag2", "tag3"},
				Offset:       0,
				Limit:        10,
				NoLimit:      false,
			},
			expected: 1,
		},
		{
			name: "no limit",
			arg: FindMemoArg{
				IDs:          []uint{1, 3},
				ContentLikes: []string{"tag1", "tag2", "tag3"},
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
					return model.CountMemos(context.Background(), tx, tc.arg)
				})
				s.NoError(err)
				s.Equal(tc.expected, count)

			})
		})
	}
}

func (s *MemoRepositorylSuite) TestFindMemo() {

	type TestArg struct {
		name    string
		In      FindMemoArg
		Out     []model.Memo
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT * FROM `memos` WHERE")

		query.WriteString(" 1 = 1")
		if arg.In.IDs != nil {
			query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(arg.In.IDs), ",")))
			for _, id := range arg.In.IDs {
				args = append(args, id)
			}
		}
		if arg.In.ContentLikes != nil {
			query.WriteString(fmt.Sprintf(" AND `content` Like (%s)", RepeatWithSep("?", len(arg.In.ContentLikes), ",")))
			for _, name := range arg.In.ContentLikes {
				args = append(args, name)
			}
		}

		if arg.In.StatusLikes != 0 {
			query.WriteString(" AND `status` = ?")
			args = append(args, arg.In.StatusLikes)
		}
		if !arg.In.NoLimit {
			query.WriteString(" ORDER BY `id` DESC")
			query.WriteString(" LIMIT ? OFFSET ?")
			args = append(args, arg.In.Limit, arg.In.Offset)
		}

		rows := mock.NewRows([]string{"id", "content", "status", "created_time", "modified_time"})

		for _, tag := range arg.Out {
			rows.AddRow(s.FetchMemoRow(tag)...)
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
			In: FindMemoArg{
				IDs:          []uint{1, 2, 3},
				ContentLikes: []string{"tag1", "tag2", "tag3"},
				Offset:       0,
				Limit:        10,
				NoLimit:      false,
			},
			Out: []model.Memo{
				{
					ID:      model.ID{ID: 1},
					Content: "tag1",
					Status:  1,
					Timestamps: model.Timestamps{
						CreatedTime:  time.Now(),
						ModifiedTime: time.Now(),
					},
				},
			},
		},
		{
			name: "no limit",
			In: FindMemoArg{
				IDs:          []uint{1, 3},
				ContentLikes: []string{"tag1", "tag2", "tag3"},
				Offset:       0,
				Limit:        10,
				NoLimit:      true,
			},
			Out: []model.Memo{
				{
					ID:      model.ID{ID: 1},
					Content: "tag1",
					Status:  1,
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
					model := NewMemo()
					return model.FindMemos(context.Background(), tx, tc.In)
				})
				s.NoError(err)

				tags, ok := result.([]model.Memo)
				s.True(ok)

				s.Equal(len(tc.Out), len(tags))
			})
		})
	}
}

func (s *MemoRepositorylSuite) TestGetMemo() {
	type TestArg struct {
		name    string
		In      uint
		Out     model.Memo
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT * FROM `memos` WHERE")

		query.WriteString(" `id` = ?")
		args = append(args, arg.In)
		query.WriteString(" ORDER BY `memos`.`id` LIMIT 1")

		rows := mock.NewRows([]string{"id", "content", "status", "created_time", "modified_time"})

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnRows(rows.AddRow(s.FetchMemoRow(arg.Out)...))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			name: "normal",
			In:   1,
			Out: model.Memo{
				ID:      model.ID{ID: 1},
				Content: "tag1",
				Status:  1,
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
					model := NewMemo()
					return model.GetMemo(context.Background(), tx, tc.In)
				})
				s.NoError(err)

				tags, ok := result.(model.Memo)
				s.True(ok)

				s.Equal(tc.Out, tags)
			})
		})
	}
}

func (s *MemoRepositorylSuite) TestCreateMemo() {
	type TestArg struct {
		name    string
		In      model.Memo
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("INSERT INTO `memos`")

		query.WriteString(" (`content`,`status`,`created_time`,`modified_time`,`id`) VALUES (?,?,?,?,?)")
		args = append(args, arg.In.Content, arg.In.Status, arg.In.CreatedTime, arg.In.ModifiedTime, arg.In.ID.ID)

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			name: "normal",
			In: model.Memo{
				ID:      model.ID{ID: 1},
				Content: "tag1",
				Status:  1,
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
					model := NewMemo()
					return nil, model.CreateMemo(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func (s *MemoRepositorylSuite) TestUpdateMemo() {
	type TestArg struct {
		name    string
		In      model.Memo
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("UPDATE `memos` SET")

		query.WriteString(" `id`=?,`content`=?,`status`=?,`created_time`=?,`modified_time`=? WHERE `id` = ?")
		args = append(args, arg.In.ID.ID, arg.In.Content, arg.In.Status, arg.In.CreatedTime, arg.In.ModifiedTime, arg.In.ID.ID)

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			name: "normal",
			In: model.Memo{
				ID:      model.ID{ID: 1},
				Content: "tag1",
				Status:  1,
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
					model := NewMemo()
					return nil, model.UpdateMemo(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func (s *MemoRepositorylSuite) TestDeleteMemo() {
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
		query.WriteString("DELETE FROM `memos` WHERE")

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
					model := NewMemo()
					return nil, model.DeleteMemo(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func TestMemoRepositorylSuite(t *testing.T) {
	suite.Run(t, new(MemoRepositorylSuite))
}
