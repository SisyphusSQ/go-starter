package do

import (
	"time"

	"gorm.io/gorm"
)

type AuditWhitelist struct {
	ID               int64            `gorm:"primaryKey;column:id" json:"id"`
	Cluster          string           `gorm:"column:cluster" json:"cluster"`
	Digest           string           `gorm:"column:digest" json:"digest"`
	InspectionStatus InspectionStatus `gorm:"column:inspection_status" json:"inspection_status"`
	CreatedBy        string           `gorm:"column:created_by" json:"created_by"`
	CreatedReason    string           `gorm:"column:created_reason" json:"created_reason"`
	CreatedAt        time.Time        `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time        `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt        gorm.DeletedAt   `gorm:"column:deleted_at" json:"deleted_at"`
}

func (AuditWhitelist) TableName() string {
	return "audit_whitelist"
}

type InspectionStatus int8

const (
	Unused InspectionStatus = iota
	Used
)
