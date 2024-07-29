package do

import "time"

type AuditSQLDigest struct {
	ID          int64     `gorm:"primaryKey;column:id" json:"id"`
	Cluster     string    `gorm:"column:cluster" json:"cluster"`
	InstanceIP  string    `gorm:"column:instance_ip" json:"instance_ip"`
	Digest      string    `gorm:"column:digest" json:"digest"`
	Pattern     string    `gorm:"column:pattern" json:"pattern"`
	SQLExample  string    `gorm:"column:sql_example" json:"sql_example"`
	DB          string    `gorm:"column:db" json:"db"`
	Cmd         string    `gorm:"column:cmd" json:"cmd"`
	CurCount    int64     `gorm:"column:cur_count" json:"cur_count"`
	ExcCount    int64     `gorm:"column:exc_count" json:"exc_count"`
	LastExcTime time.Time `gorm:"column:last_exc_time" json:"last_exc_time"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (AuditSQLDigest) TableName() string { return "audit_sql_digest" }
