package model

import "time"

type JobInfo struct {
	ID        uint64
	Name      string    // 职位名称
	Company   string    // 公司
	City      string    // 所在城市（地区）
	Salary    string    // 薪资
	Education string    // 学历要求
	WorkExp   string    // 工作年限要求
	JD        string    // 职位描述
	Welfare   string    // 福利
	URL       string    // 职位详情URL
	CreatedAt time.Time // 创建时间
}
