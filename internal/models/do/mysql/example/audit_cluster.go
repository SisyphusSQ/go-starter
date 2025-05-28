package example

import "time"

type AuditCluster struct {
	ID               int64            `gorm:"primaryKey;column:id" json:"id"`
	Cluster          string           `gorm:"column:cluster" json:"cluster"`
	InspectionStatus InspectionStatus `gorm:"column:inspection_status" json:"inspection_status"`
	CreatedAt        time.Time        `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time        `gorm:"column:updated_at" json:"updated_at"`
}

func (AuditCluster) TableName() string {
	return "audit_cluster"
}

type InspectionStatus int8

const (
	Unused InspectionStatus = iota
	Used
)
