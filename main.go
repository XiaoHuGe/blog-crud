package main

import (
	"fmt"
	"github.com/fvbock/endless"
	"log"
	"syscall"
	"xhblog/models"
	"xhblog/routers"
	"xhblog/utils/gredis"
	"xhblog/utils/logging"
	"xhblog/utils/setting"
)

func main() {
	setting.Setup()
	models.Setup()
	logging.Setup()
	gredis.Setup()

	endless.DefaultReadTimeOut = setting.ServerSetting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.ServerSetting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20
	logging.Info("DefaultMaxHeaderBytes: ", 1 << 20)
	endlessPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)

	router := routers.InitRouter()
	server := endless.NewServer(endlessPoint, router)
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}
	//s := &http.Server{
	//	// Addr需要加':' :9000"
	//	Addr:           fmt.Sprintf(":%d", setting.HttpPort),
	//	Handler:        router,
	//	ReadTimeout:    setting.ReadTimeout,
	//	WriteTimeout:   setting.WriteTimeout,
	//	MaxHeaderBytes: 1 << 20,
	//}
	//s.ListenAndServe()
}
