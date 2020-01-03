package global

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/spf13/viper"
)

var once = new(sync.Once)

var (
	NeedAll   = flag.Bool("all", false, "是否全量抓取，默认否")
	WhichSite = flag.String("site", "", "抓取哪个站点(空标示所有站点)")
	Config    = flag.String("c", "config", "环境变量档案名称，默认 config")
)

func Init() {
	once.Do(func() {
		if !flag.Parsed() {
			flag.Parse()
		}

		rand.Seed(time.Now().UnixNano())

		viper.SetConfigName(*Config)
		viper.AddConfigPath("/etc/crawler/")
		viper.AddConfigPath("$HOME/.crawler")
		viper.AddConfigPath(App.RootDir + "/config")
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}

	})
}
