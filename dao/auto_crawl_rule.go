package dao

import (
	"log"

	"github.com/wsloong/go-crawler/model"
)

type autoRuleDAO struct{}

var DefaultAutoRule = autoRuleDAO{}

func (a autoRuleDAO) FindOne(domain string) (*model.AutoCrawlRule, error) {
	strSql := `SELECT 
			id,domain,all_url,incr_url,keywords,list_selector,result_selector,page_field,max_page,ext,created_at
		FROM auto_crawl_rule WHERE domain=?`
	row := masterDB.QueryRow(strSql, domain)

	autoRule := &model.AutoCrawlRule{}
	err := row.Scan(&autoRule.ID,
		&autoRule.Domain,
		&autoRule.AllURL,
		&autoRule.IncrURL,
		&autoRule.Keywords,
		&autoRule.ListSelector,
		&autoRule.ResultSelector,
		&autoRule.PageField,
		&autoRule.MaxPage,
		&autoRule.Ext,
		&autoRule.CreatedAt)

	if err != nil {
		log.Println("auto crawl rule scan error:", err)
		return nil, err
	}

	return autoRule, nil
}

func (a autoRuleDAO) Find() ([]*model.AutoCrawlRule, error) {
	strSql := `SELECT 
			id,domain,all_url,incr_url,keywords,list_selector,result_selector,page_field,max_page,ext,created_at 
		FROM auto_crawl_rule WHERE state=?`
	rows, err := masterDB.Query(strSql, model.AutoCrawlOn)
	if err != nil {
		return nil, err
	}

	autoRuleSlice := make([]*model.AutoCrawlRule, 0, 10)
	for rows.Next() {
		autoRule := &model.AutoCrawlRule{}
		dest := []interface{}{
			&autoRule.ID,
			&autoRule.Domain,
			&autoRule.AllURL,
			&autoRule.IncrURL,
			&autoRule.Keywords,
			&autoRule.ListSelector,
			&autoRule.ResultSelector,
			&autoRule.PageField,
			&autoRule.MaxPage,
			&autoRule.Ext,
			&autoRule.CreatedAt,
		}

		err = rows.Scan(dest...)
		if err != nil {
			log.Println("auto crawl rule scan error:", err)
			continue
		}

		autoRuleSlice = append(autoRuleSlice, autoRule)
	}

	return autoRuleSlice, nil
}
