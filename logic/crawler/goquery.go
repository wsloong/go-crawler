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
	"github.com/wsloong/go-crawler/model"
	"github.com/wsloong/go-crawler/util"
)

type GoQueryParser struct {
	isAll    bool                 // 是否全量抓取
	autoRule *model.AutoCrawlRule // 要自动抓取的网站列表配置信息
	d        time.Duration        // 抓取间隔，避免过快

	crawlRule *model.CrawlRule // 网站详情页配置信息
}

func NewGoQueryParser(isAll bool, autoRule *model.AutoCrawlRule, d time.Duration) *GoQueryParser {
	return &GoQueryParser{
		isAll:    isAll,
		autoRule: autoRule,
		d:        d,
	}
}

// Work 抓取具体网站
func (q *GoQueryParser) Work() {
	err := q.ListParse()
	if err != nil {
		log.Println(" go query parse work error: ", err)
	}
}

// ListParse 网站职位列表解析
func (q *GoQueryParser) ListParse() error {
	var err error
	keywords := strings.Split(q.autoRule.Keywords, ",")

	for _, kw := range keywords {
		time.Sleep(q.d)
		log.Println("start list parse, the keyword is ", kw)
		if q.isAll {
			err = q.parseOneKeyword(kw, q.autoRule.MaxPage)
		} else {
			err = q.parseOneKeyword(kw, 1)
		}

		if err != nil {
			log.Println("parse one keyword: ", kw, "error: ", err)
		}
	}
	return err
}

// parseOneKeyword 根据一个关键词进行解析
func (q *GoQueryParser) parseOneKeyword(kw string, maxPage int) error {
	var globalErr error

	pURL, err := url.Parse(q.autoRule.AllURL)
	if err != nil {
		return err
	}

	for i := 0; i < maxPage; i++ {
		time.Sleep(q.d)

		allURL := fmt.Sprintf(q.autoRule.AllURL+"&%s=%d", kw, q.autoRule.PageField, i+1)
		resp, err := util.HTTPGet(allURL)
		if err != nil {
			log.Println("get url: ", allURL, "error: ", err)
			globalErr = err
			continue
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Println("new goquery document error: ", err)
			return err
		}

		listSel := doc.Find(q.autoRule.ListSelector)
		if listSel.Length() == 0 {
			log.Println("The ip may be forbidden!")
			globalErr = errors.New("the ip may be forbidden")
			continue
		}

		listSel.Each(func(i int, sel *goquery.Selection) {
			href, ok := sel.Find(q.autoRule.ResultSelector).Attr("href")
			if ok {
				if !strings.HasPrefix(href, "http") {
					href = pURL.Scheme + "://" + q.autoRule.Domain + href
				}
				err = q.DetailParse(href)
				if err != nil {
					log.Println("detail parse error: ", err)
					return
				}

				log.Println("Detail Parse successfully!")
			} else {
				log.Println("parse job detail url error, href is not exists!")
			}
		})
	}
	return globalErr
}

// DetailParse 职位详情页解析
func (q *GoQueryParser) DetailParse(detailURL string) error {
	log.Println("start crawl detail url: ", detailURL)

	job, _ := dao.DefaultJob.FindOneByURL(context.TODO(), detailURL)
	if job != nil && job.ID > 0 {
		log.Println("had exists!")
		return nil
	}

	time.Sleep(q.d)

	var err error
	if q.crawlRule == nil {
		q.crawlRule, err = dao.DefaultRuleDAO.FindOne(q.autoRule.Domain)
		if err != nil {
			return err
		}
	}

	resp, err := util.HTTPGet(detailURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	jobName := doc.Find(q.crawlRule.JobName).Text()
	if jobName == "" {
		return errors.New("the job name is empty")
	}

	body, err := doc.Html()
	if err != nil {
		return err
	}

	job = &model.JobInfo{URL: detailURL}
	job.Name = jobName
	job.Company = doc.Find(q.crawlRule.JobCompany).Text()

	job.City, err = q.parseContent(doc, q.crawlRule.JobCity, body)
	if err != nil {
		return err
	}

	job.WorkExp, err = q.parseContent(doc, q.crawlRule.JobWorkExp, body)
	if err != nil {
		return err
	}

	job.Education, err = q.parseContent(doc, q.crawlRule.JobEducation, body)
	if err != nil {
		return err
	}

	job.Salary = doc.Find(q.crawlRule.JobSalary).Text()
	job.JD = doc.Find(q.crawlRule.JobJD).Text()
	job.Welfare = doc.Find(q.crawlRule.JobWelfare).Text()

	return dao.DefaultJob.Create(job)

}

func (q *GoQueryParser) parseContent(doc *goquery.Document, selector, content string) (string, error) {
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
			return "", errors.New("The regex `" + selector + "` is illegal!")
		}
		parsedContent = result[1]
	}
	return parsedContent, nil
}
