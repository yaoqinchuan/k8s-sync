package do

import "github.com/gogf/gf/v2/os/gtime"

type AsyncTask struct {
	Id                    int64       `json:"id"       `
	TaskName              string      `json:"task_name" `
	TaskAttributes        string      `json:"task_attributes" `
	Ip                    string      `json:"ip" `
	IpHeartBeatLifePeriod int         `json:"ip_heart_beat_life_period" `
	Status                string      `json:"status"`
	ErrorInfo             string      `json:"error_info"`
	RetryTime             int         `json:"retry_time"`
	TotalRetryTime        int         `json:"total_retry_time"`
	TaskStartTime         *gtime.Time `json:"task_start_time" `
	TaskTimeoutTime       float64     `json:"task_timeout_time" `
	TaskEndTime           *gtime.Time `json:"task_end_time" `
	CreateAt              *gtime.Time `json:"create_at"`
	UpdateAt              *gtime.Time `json:"update_at"`
	Deleted               int         `json:"deleted"`
}
