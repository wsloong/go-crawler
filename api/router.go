package api

import "github.com/wsloong/go-crawler/api/handlers"

func RegisterRouters() {
	handlers.JobHandler.RegisterRoutes()
}
