package main

import (
	"github.com/robfig/cron"
	"github.com/spf13/viper"
	"github.com/wsloong/go-crawler/global"
	"github.com/wsloong/go-crawler/logic/crawler"
)

func ServeBackGround() {
	parser := viper.GetString("crawl.parser")

	if *global.NeedAll {
		go crawler.DoCrawl(parser, true, *global.WhichSite)
	}

	// 定时增量
	c := cron.New()

	viper.SetDefault("crawl.spec", "0 0 */1 * * ?")
	c.AddFunc(viper.GetString("crawl.spec"), func() {
		crawler.DoCrawl(parser, false, *global.WhichSite)
	})
	c.Start()
}
