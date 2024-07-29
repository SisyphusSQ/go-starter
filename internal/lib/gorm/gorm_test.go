package gormv2

import (
	"fmt"
	"testing"
	"time"
)

var testSQL = "select now()"

func TestMySQL(t *testing.T) {
	var (
		err     error
		timeStr string
	)

	db := New(&GormConfig{
		Alias:        "sql_audit_test",
		Type:         "mysql",
		Server:       "sqlaudit.tidb.qiyi.database",
		Port:         4000,
		Database:     "sql_audit_test",
		User:         "sql_audit",
		Password:     "2m#qfqE!0e9g",
		MaxIdleConns: 200,
		MaxOpenConns: 500,
		Charset:      "utf8mb4",
		MaxLeftTime:  time.Second * 10,
	})

	if err = db.gorm.Statement.Error; err != nil {
		t.Error(err)
	}

	if err = db.gorm.Raw(testSQL).Scan(&timeStr).Error; err != nil {
		t.Error(err)
	}

	fmt.Printf("mysql now time: %s\n", timeStr)
}

func TestClickhouse(t *testing.T) {
	var (
		err     error
		timeStr string
	)

	db := New(&GormConfig{
		Alias:        "test",
		Type:         "clickhouse",
		Server:       "sqlrecord.proxy.ck.qiyi.database",
		Port:         8027,
		Database:     "idba",
		User:         "idba",
		Password:     "idbaAAA",
		MaxIdleConns: 2,
		MaxOpenConns: 2,
		MaxLeftTime:  time.Second * 10,
	})

	if err = db.gorm.Statement.Error; err != nil {
		t.Error(err)
	}

	if err = db.gorm.Raw(testSQL).Scan(&timeStr).Error; err != nil {
		t.Error(err)
	}

	fmt.Printf("clickhouse now time: %s\n", timeStr)
}
