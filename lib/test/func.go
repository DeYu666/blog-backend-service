package test

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DeYu666/blog-backend-service/lib/client"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func WarpRepositoryTest(suite *suite.Suite, f func(sqlmock.Sqlmock)) {
	db, mock := mockDB()
	defer db.Close()

	f(mock)
}

func mockDB() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "postgres",
		SkipInitializeWithVersion: true,
		Conn:                      db,
	})

	gdb, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic(err)
	}

	client.Mysql.SetupDB(gdb)

	return db, mock
}

func MockTx(f func(c context.Context, tx *gorm.DB) (interface{}, error)) (interface{}, error) {

	var ret interface{}
	var err error

	client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		c := context.Background()

		ret, err = f(c, tx)

		return nil
	})

	return ret, err
}

func ToDriveVaules(arr []interface{}) ([]driver.Value, error) {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("arr must be slice")
	}

	l := v.Len()

	ret := make([]driver.Value, l)

	for i := 0; i < l; i++ {
		ret[i] = v.Index(i).Interface()
	}

	return ret, nil
}
