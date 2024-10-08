package do

import "time"

type TaskResult struct {
	ID         int64      `gorm:"primaryKey;column:id" json:"id"`
	UUID       string     `gorm:"column:uuid" json:"uuid"`
	TaskName   string     `gorm:"column:task_name" json:"task_name"`
	ExecIP     string     `gorm:"column:exec_ip" json:"exec_ip"`
	TaskStatus TaskStatus `gorm:"column:task_status" json:"task_status"`
	TaskCost   int64      `gorm:"column:task_cost" json:"task_cost"`
	ErrMsg     string     `gorm:"column:err_msg" json:"err_msg"`
	CreatedAt  time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at" json:"updated_at"`
	Start      time.Time  `gorm:"-" json:"start"`
}

func (TaskResult) TableName() string {
	return "task_result"
}

type TaskStatus int

const (
	_ TaskStatus = iota
	Processing
	Finish
	Error
)
