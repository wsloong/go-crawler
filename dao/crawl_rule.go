package dao

import (
	"log"

	"github.com/wsloong/go-crawler/model"
)

type crawlRuleDAO struct{}

var DefaultRuleDAO = crawlRuleDAO{}

func (a crawlRuleDAO) FindOne(domain string) (*model.CrawlRule, error) {
	strSql := "SELECT id,domain,name,job_name,job_company,job_city,job_work_exp,job_salary,job_education,job_jd,job_welfare,created_at FROM crawl_rule WHERE domain=?"
	row := masterDB.QueryRow(strSql, domain)

	crawRule := &model.CrawlRule{}
	err := row.Scan(&crawRule.ID,
		&crawRule.Domain,
		&crawRule.Name,
		&crawRule.JobName,
		&crawRule.JobCompany,
		&crawRule.JobCity,
		&crawRule.JobWorkExp,
		&crawRule.JobSalary,
		&crawRule.JobEducation,
		&crawRule.JobJD,
		&crawRule.JobWelfare,
		&crawRule.CreatedAt)
	if err != nil {
		log.Println("crawl rule scan error: ", err)
		return nil, err
	}

	return crawRule, nil
}
