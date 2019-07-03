package main

import (
	"fmt"
	"net/http"
	"xhblog/utils/setting"
	"xhblog/routers"
)

func main() {
	router := routers.InitRouter()
	s := &http.Server{
		// Addr需要加':' :9000"
		Addr: fmt.Sprintf(":%d", setting.HttpPort),
		Handler: router,
		ReadTimeout: setting.ReadTimeout,
		WriteTimeout: setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
