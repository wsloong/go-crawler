package crawler

import (
	"log"
	"time"

	"github.com/spf13/viper"
	"github.com/wsloong/go-crawler/util"

	"github.com/wsloong/go-crawler/dao"
	"github.com/wsloong/go-crawler/model"
)

type Parser interface {
	ListParse() error
	DetailParse(url string) error
}

func DoCrawl(parser string, isAll bool, whichSite string) error {
	log.Println("start crawling ...")

	var autoRuleSlice []*model.AutoCrawlRule
	var err error

	if whichSite == "" {
		autoRuleSlice, err = dao.DefaultAutoRule.Find()
		if err != nil {
			log.Println("find auto rule error: ", err)
			return err
		}
	} else {
		autoRule, err := dao.DefaultAutoRule.FindOne(whichSite)
		if err != nil {
			log.Println("find one auto rule error: ", err)
			return err
		}
		autoRuleSlice = []*model.AutoCrawlRule{autoRule}
	}

	min := util.Min(len(autoRuleSlice), viper.GetInt("crawl.concurrency_num"))
	pool := NewPool(min)

	d, err := time.ParseDuration(viper.GetString("crawl.sleep"))
	if err != nil {
		d = 20 * time.Second
	}

	for _, autoRule := range autoRuleSlice {
		var worker Worker

		if parser == "goquery" {
			worker = NewGoQueryParser(isAll, autoRule, d)
		} else {
			worker = NewCollyParser(isAll, autoRule, d)
		}
		pool.Run(worker)
	}

	pool.Shutdown()
	return nil
}
