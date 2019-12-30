package model

import "time"

type CrawlRule struct {
	ID           uint64
	Domain       string    // 来源域名（不一定是顶级域名）
	Name         string    // 来源名称
	JobName      string    // 职位名称规则
	JobCompany   string    // 职位公司规则
	JobCity      string    // 职位所在城市规则
	JobWorkExp   string    // 职位工作年限要求
	JobSalary    string    // 职位薪资规则
	JobEducation string    // 职位学历规则
	JobJd        string    // 职位 JD 规则
	JobWelfare   string    // 职位福利规则
	CreatedAt    time.Time // 创建时间
}
