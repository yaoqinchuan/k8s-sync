package model

import "github.com/gogf/gf/v2/os/gtime"

type AsyncTaskModel struct {
	Id                  int64       `json:"id" v:"required"`
	TaskName            string      `json:"task_name" v:"required"`
	TaskAttributes      string      `json:"task_attributes" v:"required"`
	Ip                  string      `json:"ip" v:"required"`
	IpHeartBeatLifeTime *gtime.Time `json:"ip_heart_beat_life_time" v:"required"`
	Status              string      `json:"status" v:"required"`
	ErrorInfo           string      `json:"error_info" v:"required"`
	RetryTime           int         `json:"retry_time" v:"required"`
	TotalRetryTime      int         `json:"total_retry_time" v:"required"`
	TaskStartTime       *gtime.Time `json:"task_start_time" v:"required"`
	TaskEndTime         *gtime.Time `json:"task_end_time" v:"required"`
	TaskTimeoutTime     float64     `json:"task_timeout_time"  v:"required"`
	CreateAt            *gtime.Time `json:"create_at" v:"required"`
	UpdateAt            *gtime.Time `json:"update_at" v:"required"`
	Deleted             int         `json:"deleted" v:"required"`
}
