package model

import "time"

const (
	AutoCrawlOn = iota
	AutoCrawlOff
)

type AutoCrawlRule struct {
	ID             int64
	Domain         string    // 来源域名（不一定是顶级域名
	AllURL         string    // 全量url，关键词占位符使用%s
	IncrURL        string    // 增量url，关键词占位符使用%s
	Keywords       string    // 搜索关键词，多个逗号分隔
	ListSelector   string    // 列表选择器
	ResultSelector string    // 结果选择器，获取具体职位的 url
	PageField      string    // 分页字段名
	MaxPage        int       // 全量最多抓取多少页
	Ext            string    // 扩展信息，某些网站的特殊配置，json格式
	State          int       // 状态：0-自动抓取；1-停止抓取
	CreatedAt      time.Time // 创建时间
}
