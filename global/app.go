package global

import (
	"os"
	"sync"
	"time"

	"github.com/wsloong/go-crawler/util"

	"github.com/spf13/viper"
)

func init() {
	App.Name = os.Args[0]
	App.Version = "V1.0.0"
	App.LaunchTime = time.Now()

	App.RootDir = "."

	if !viper.InConfig("http.port") {
		App.RootDir = inferRootDir()
	}

	fileInfo, err := os.Stat(os.Args[0])
	if err != nil {
		panic(err)
	}

	App.Date = fileInfo.ModTime()
}

func inferRootDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var infer func(d string) string
	infer = func(d string) string {
		if util.Exist(d + "/config") {
			return d
		}
		return infer(cwd)
	}

	return infer(cwd)
}

var App = &app{}

type app struct {
	Name    string
	Version string
	Date    time.Time

	// 根目录
	RootDir string

	// 启动时间
	LaunchTime time.Time
	Uptime     time.Duration

	Locker sync.Mutex
}

func (a *app) SetUptime() {
	a.Locker.Lock()
	defer a.Locker.Unlock()
	a.Uptime = time.Now().Sub(a.LaunchTime)
}
