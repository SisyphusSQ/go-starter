package do

import (
	"time"
)

type AuditSQLBadUse struct {
	ID                int64     `gorm:"primaryKey;column:id" json:"id"`
	Cluster           string    `gorm:"column:cluster" json:"cluster"`
	CurCount          int64     `gorm:"column:cur_count" json:"cur_count"`
	ExcCount          int64     `gorm:"column:exc_count" json:"exc_count"`
	Digest            string    `gorm:"column:digest" json:"digest"`
	InstanceIP        string    `gorm:"column:instance_ip" json:"instance_ip"`
	ClientUser        string    `gorm:"column:client_user" json:"client_user"`
	ClientIP          string    `gorm:"column:client_ip" json:"client_ip"`
	DB                string    `gorm:"column:db" json:"db"`
	Cmd               string    `gorm:"column:cmd" json:"cmd"`
	Pattern           string    `gorm:"column:pattern" json:"pattern"`
	SQLExample        string    `gorm:"column:sql_example" json:"sql_example"`
	NoIndexUsed       int8      `gorm:"column:no_index_used" json:"no_index_used"`
	NoGoodIndexUsed   int8      `gorm:"column:no_good_index_used" json:"no_good_index_used"`
	SortScan          int8      `gorm:"column:sort_scan" json:"sort_scan"`
	AvgExecTime       int64     `gorm:"column:avg_exec_time" json:"avg_exec_time"`
	MaxExecTime       int64     `gorm:"column:max_exec_time" json:"max_exec_time"`
	AvgSentRows       int64     `gorm:"column:avg_sent_rows" json:"avg_sent_rows"`
	MaxSentRows       int64     `gorm:"column:max_sent_rows" json:"max_sent_rows"`
	AvgRowsExamined   int64     `gorm:"column:avg_rows_examined" json:"avg_rows_examined"`
	MaxRowsExamined   int64     `gorm:"column:max_rows_examined" json:"max_rows_examined"`
	AvgAffectedRows   int64     `gorm:"column:avg_affected_rows" json:"avg_affected_rows"`
	MaxAffectedRows   int64     `gorm:"column:max_affected_rows" json:"max_affected_rows"`
	AvgLockTime       int64     `gorm:"column:avg_lock_time" json:"avg_lock_time"`
	MaxLockTime       int64     `gorm:"column:max_lock_time" json:"max_lock_time"`
	AvgSortRows       int64     `gorm:"column:avg_sort_rows" json:"avg_sort_rows"`
	MaxSortRows       int64     `gorm:"column:max_sort_rows" json:"max_sort_rows"`
	LastExcTime       time.Time `gorm:"column:last_exc_time" json:"last_exc_time"`
	BadReason         string    `gorm:"column:bad_reason" json:"bad_reason"`
	CreatedAt         time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at" json:"updated_at"`
	ExecTimeCount     int64     `gorm:"-" json:"exec_time_count"`
	SentRowsCount     int64     `gorm:"-" json:"sent_rows_count"`
	RowsExaminedCount int64     `gorm:"-" json:"rows_examined_count"`
	AffectedRowsCount int64     `gorm:"-" json:"affected_rows_count"`
	LockTimeCount     int64     `gorm:"-" json:"lock_time_count"`
	SortRowsCount     int64     `gorm:"-" json:"sort_rows_count"`
	Port              int       `gorm:"-" json:"port"`
}

func (AuditSQLBadUse) TableName() string { return "audit_sql_bad_use" }
