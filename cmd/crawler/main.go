package main

import (
	"log"
	"net/http"

	"github.com/spf13/viper"
	"github.com/wsloong/go-crawler/api"
)

func main() {
	// 注册路由
	api.RegisterRouters()

	viper.SetDefault("http.port", "9090")
	host := viper.GetString("http.host")
	port := viper.GetString("http.port")

	addr := host + ":" + port

	log.Println("service listen at: ", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
