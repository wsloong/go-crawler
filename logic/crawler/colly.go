package crawler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/wsloong/go-crawler/dao"

	"github.com/gocolly/colly"
	"github.com/wsloong/go-crawler/model"
)

type CollyParse struct {
	isAll    bool                 // 是否全量抓取
	autoRule *model.AutoCrawlRule // 要自动抓取的网站列表
	d        time.Duration        // 抓取间隔，避免过快

	crawlRule *model.CrawlRule // 网站详情页配置信息
}

func NewCollyParser(isAll bool, autoRule *model.AutoCrawlRule, d time.Duration) *CollyParse {
	return &CollyParse{
		isAll:    isAll,
		autoRule: autoRule,
		d:        d,
	}
}

func (c *CollyParse) Work() {
	err := c.ListParse()
	if err != nil {
		log.Println(" colly parse work error: ", err)
	}
}

func (c *CollyParse) ListParse() error {
	keywords := strings.Split(c.autoRule.Keywords, ",")
	var err error

	for _, kw := range keywords {
		time.Sleep(c.d)

		log.Println("start list parse, the key is: ", kw)
		if c.isAll {
			err = c.parseOneKeyword(kw, c.autoRule.MaxPage)
		} else {
			err = c.parseOneKeyword(kw, 1)
		}

		if err != nil {
			log.Println("parse one keyword:", kw, "error:", err)
		}
	}

	return err
}

// parseOneKeyword 根据一个关键词进行解析
func (c *CollyParse) parseOneKeyword(kw string, maxPage int) error {
	var globalErr error

	pURL, err := url.Parse(c.autoRule.AllURL)
	if err != nil {
		return nil
	}

	cc := colly.NewCollector()
	cc.Limit(&colly.LimitRule{
		Parallelism: 1,
		RandomDelay: c.d * time.Second,
	})
	cc.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36")
	})
	cc.OnHTML(c.autoRule.ListSelector, func(e *colly.HTMLElement) {
		href, ok := e.DOM.Find(c.autoRule.ResultSelector).Attr("href")
		if ok {
			if !strings.HasPrefix(href, "http") {
				href = pURL.Scheme + "://" + c.autoRule.Domain + href
			}
			err = c.DetailParse(href)
			if err != nil {
				log.Println("detail parse error: ", err)
				return
			}
			log.Println("detail parse successfully!")
		} else {
			log.Println("parse job detail url error, href is not existes!")
		}
	})
	cc.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "; ERROR:", err)
		globalErr = err
		return
	})

	for i := 0; i < maxPage; i++ {
		allURL := fmt.Sprintf(c.autoRule.AllURL+"&%s=%d", kw, c.autoRule.PageField, i+1)
		// 开始抓取
		cc.Visit(allURL)
	}
	return globalErr
}

// DetailParse 获取职位详情页
func (c *CollyParse) DetailParse(detailURL string) error {
	log.Println("start crawl detail url: ", detailURL)

	var globalErr error

	job, _ := dao.DefaultJob.FindOneByURL(context.TODO(), detailURL)

	if job != nil && job.ID > 0 {
		log.Println("had exists")
		return nil
	}
	time.Sleep(c.d)

	var err error
	if c.crawlRule == nil {
		c.crawlRule, err = dao.DefaultRuleDAO.FindOne(c.autoRule.Domain)
		if err != nil {
			return err
		}
	}

	cc := colly.NewCollector()
	cc.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36")
	})
	cc.OnHTML("html", func(e *colly.HTMLElement) {
		doc := e.DOM

		jobName := doc.Find(c.crawlRule.JobName).Text()
		if jobName == "" {
			globalErr = errors.New("the job name is empty!")
			return
		}

		body, err := doc.Html()
		if err != nil {
			globalErr = err
			return
		}

		job = &model.JobInfo{URL: detailURL}

		job.Name = jobName
		job.Company = doc.Find(c.crawlRule.JobCompany).Text()

		job.City, err = c.parseContent(doc, c.crawlRule.JobCity, body)
		if err != nil {
			globalErr = err
			return
		}

		job.WorkExp, err = c.parseContent(doc, c.crawlRule.JobWorkExp, body)
		if err != nil {
			globalErr = err
			return
		}

		job.Education, err = c.parseContent(doc, c.crawlRule.JobEducation, body)
		if err != nil {
			globalErr = err
			return
		}

		job.Salary = doc.Find(c.crawlRule.JobSalary).Text()
		job.JD = doc.Find(c.crawlRule.JobJD).Text()
		job.Welfare = doc.Find(c.crawlRule.JobWelfare).Text()

		err = dao.DefaultJob.Create(job)
		if err != nil {
			globalErr = err
			return
		}
	})

	cc.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL: ", r.Request.URL, "failed with response:", r, "; Error:", err)
		globalErr = err
		return
	})
	cc.Visit(detailURL)
	return globalErr
}

func (c *CollyParse) parseContent(doc *goquery.Selection, selector, content string) (string, error) {
	var parsedContent string

	isSelector := doc.Is(selector)
	if isSelector {
		parsedContent = doc.Find(selector).Text()
	} else {
		reg, err := regexp.Compile(selector)
		if err != nil {
			return "", err
		}

		result := reg.FindStringSubmatch(content)
		if len(result) < 2 {
			return "", errors.New("the regex `" + selector + "` is illegal")
		}
		parsedContent = result[1]
	}
	return parsedContent, nil
}
