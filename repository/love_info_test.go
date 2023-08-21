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

type LoveInfoRepositorylSuite struct {
	suite.Suite
}

func (suite *LoveInfoRepositorylSuite) FetchLoveInfoRow(LoveInfo model.LoveInfo) []driver.Value {
	return []driver.Value{LoveInfo.ID.ID, LoveInfo.KnownTime, LoveInfo.ConfessionTime, LoveInfo.LoveName, LoveInfo.ExtraInfo}
}

func (s *LoveInfoRepositorylSuite) TestCountLoveInfo() {
	model := NewLoveInfo()

	mockFunc := func(mock sqlmock.Sqlmock, arg FindLoveInfoArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT count(*) FROM `love_infos` WHERE")

		query.WriteString(" 1 = 1")
		if arg.IDs != nil {
			query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(arg.IDs), ",")))
			for _, id := range arg.IDs {
				args = append(args, id)
			}
		}
		if arg.LoveNames != nil {
			query.WriteString(fmt.Sprintf(" AND `love_name` IN (%s)", RepeatWithSep("?", len(arg.LoveNames), ",")))
			for _, name := range arg.LoveNames {
				args = append(args, name)
			}
		}
		if arg.ExtraInfoLikes != nil {
			query.WriteString(fmt.Sprintf(" AND `extra_info` Like (%s)", RepeatWithSep("?", len(arg.ExtraInfoLikes), ",")))
			for _, name := range arg.ExtraInfoLikes {
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
		arg      FindLoveInfoArg
		expected int64
	}{
		{
			name: "normal",
			arg: FindLoveInfoArg{
				IDs:       []uint{1, 2, 3},
				LoveNames: []string{"LoveInfo1", "LoveInfo2", "LoveInfo3"},
				Offset:    0,
				Limit:     10,
				NoLimit:   false,
			},
			expected: 1,
		},
		{
			name: "no limit",
			arg: FindLoveInfoArg{
				IDs:       []uint{1, 3},
				LoveNames: []string{"LoveInfo1", "LoveInfo2", "LoveInfo3"},
				Offset:    0,
				Limit:     10,
				NoLimit:   true,
			},
			expected: 1,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc.arg)

				count, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					return model.CountLoveInfos(context.Background(), tx, tc.arg)
				})
				s.NoError(err)
				s.Equal(tc.expected, count)

			})
		})
	}
}

func (s *LoveInfoRepositorylSuite) TestFindLoveInfo() {

	type TestArg struct {
		name    string
		In      FindLoveInfoArg
		Out     []model.LoveInfo
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT * FROM `love_infos` WHERE")

		query.WriteString(" 1 = 1")
		if arg.In.IDs != nil {
			query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(arg.In.IDs), ",")))
			for _, id := range arg.In.IDs {
				args = append(args, id)
			}
		}
		if arg.In.LoveNames != nil {
			query.WriteString(fmt.Sprintf(" AND `love_name` IN (%s)", RepeatWithSep("?", len(arg.In.LoveNames), ",")))
			for _, name := range arg.In.LoveNames {
				args = append(args, name)
			}
		}
		if arg.In.ExtraInfoLikes != nil {
			query.WriteString(fmt.Sprintf(" AND `extra_info` Like (%s)", RepeatWithSep("?", len(arg.In.ExtraInfoLikes), ",")))
			for _, name := range arg.In.ExtraInfoLikes {
				args = append(args, name)
			}
		}

		if !arg.In.NoLimit {
			query.WriteString(" ORDER BY `id` DESC")
			query.WriteString(" LIMIT ? OFFSET ?")
			args = append(args, arg.In.Limit, arg.In.Offset)
		}

		rows := mock.NewRows([]string{"id", "known_time", "confession_time", "love_name", "extra_info"})

		for _, LoveInfo := range arg.Out {
			rows.AddRow(s.FetchLoveInfoRow(LoveInfo)...)
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
			In: FindLoveInfoArg{
				IDs:       []uint{1, 2, 3},
				LoveNames: []string{"LoveInfo1", "LoveInfo2", "LoveInfo3"},
				Offset:    0,
				Limit:     10,
				NoLimit:   false,
			},
			Out: []model.LoveInfo{
				{
					ID:             model.ID{ID: 1},
					LoveName:       "LoveInfo1",
					KnownTime:      time.Now(),
					ConfessionTime: time.Now(),
					ExtraInfo:      "ExtraInfo1",
				},
			},
		},
		{
			name: "no limit",
			In: FindLoveInfoArg{
				IDs:       []uint{1, 3},
				LoveNames: []string{"LoveInfo1", "LoveInfo2", "LoveInfo3"},
				Offset:    0,
				Limit:     10,
				NoLimit:   true,
			},
			Out: []model.LoveInfo{
				{
					ID:             model.ID{ID: 1},
					LoveName:       "LoveInfo1",
					KnownTime:      time.Now(),
					ConfessionTime: time.Now(),
					ExtraInfo:      "ExtraInfo1",
				},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				result, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewLoveInfo()
					return model.FindLoveInfos(context.Background(), tx, tc.In)
				})
				s.NoError(err)

				LoveInfos, ok := result.([]model.LoveInfo)
				s.True(ok)

				s.Equal(len(tc.Out), len(LoveInfos))
			})
		})
	}
}

func (s *LoveInfoRepositorylSuite) TestGetLoveInfo() {
	type TestArg struct {
		name    string
		In      uint
		Out     model.LoveInfo
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("SELECT * FROM `love_infos` WHERE")

		query.WriteString(" `id` = ?")
		args = append(args, arg.In)
		query.WriteString(" ORDER BY `love_infos`.`id` LIMIT 1")

		rows := mock.NewRows([]string{"id", "known_time", "confession_time", "love_name", "extra_info"})

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnRows(rows.AddRow(s.FetchLoveInfoRow(arg.Out)...))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			name: "normal",
			In:   1,
			Out: model.LoveInfo{
				ID:             model.ID{ID: 1},
				LoveName:       "LoveInfo1",
				KnownTime:      time.Now(),
				ConfessionTime: time.Now(),
				ExtraInfo:      "ExtraInfo1",
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				result, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewLoveInfo()
					return model.GetLoveInfo(context.Background(), tx, tc.In)
				})
				s.NoError(err)

				LoveInfos, ok := result.(model.LoveInfo)
				s.True(ok)

				s.Equal(tc.Out, LoveInfos)
			})
		})
	}
}

func (s *LoveInfoRepositorylSuite) TestCreateLoveInfo() {
	type TestArg struct {
		name    string
		In      model.LoveInfo
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("INSERT INTO `love_infos`")

		query.WriteString(" (`known_time`,`confession_time`,`love_name`,`extra_info`,`id`) VALUES (?,?,?,?,?)")
		args = append(args, arg.In.KnownTime, arg.In.ConfessionTime, arg.In.LoveName, arg.In.ExtraInfo, arg.In.ID.ID)

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			name: "normal",
			In: model.LoveInfo{
				ID:             model.ID{ID: 1},
				LoveName:       "LoveInfo1",
				KnownTime:      time.Now(),
				ConfessionTime: time.Now(),
				ExtraInfo:      "ExtraInfo1",
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				_, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewLoveInfo()
					return nil, model.CreateLoveInfo(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func (s *LoveInfoRepositorylSuite) TestUpdateLoveInfo() {
	type TestArg struct {
		name    string
		In      model.LoveInfo
		OutErr  bool
		Err     error
		Prepare func(mock sqlmock.Sqlmock, args TestArg)
		TxFunc  func(arg TestArg) TxFunc
	}

	mockFunc := func(mock sqlmock.Sqlmock, arg TestArg) {
		query := strings.Builder{}
		args := []interface{}{}
		query.WriteString("UPDATE `love_infos` SET")

		query.WriteString(" `id`=?,`known_time`=?,`confession_time`=?,`love_name`=?,`extra_info`=? WHERE `id` = ?")
		args = append(args, arg.In.ID.ID, arg.In.KnownTime, arg.In.ConfessionTime, arg.In.LoveName, arg.In.ExtraInfo, arg.In.ID.ID)

		values, err := test.ToDriveVaules(args)
		s.NoError(err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query.String())).WithArgs(values...).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	testCases := []TestArg{
		{
			name: "normal",
			In: model.LoveInfo{
				ID:             model.ID{ID: 1},
				LoveName:       "LoveInfo1",
				KnownTime:      time.Now(),
				ConfessionTime: time.Now(),
				ExtraInfo:      "ExtraInfo1",
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			test.WarpRepositoryTest(&s.Suite, func(db sqlmock.Sqlmock) {

				mockFunc(db, tc)

				_, err := test.MockTx(func(c context.Context, tx *gorm.DB) (interface{}, error) {
					model := NewLoveInfo()
					return nil, model.UpdateLoveInfo(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func (s *LoveInfoRepositorylSuite) TestDeleteLoveInfo() {
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
		query.WriteString("DELETE FROM `love_infos` WHERE")

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
					model := NewLoveInfo()
					return nil, model.DeleteLoveInfo(context.Background(), tx, tc.In)
				})
				s.NoError(err)
			})
		})
	}
}

func TestLoveInfoRepositorylSuite(t *testing.T) {
	suite.Run(t, new(LoveInfoRepositorylSuite))
}
