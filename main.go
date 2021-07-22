package main

import (
	"GoGin/config"
	"GoGin/config/logging"

	"GoGin/models"
	"GoGin/routers"
	"fmt"
	"github.com/fvbock/endless"
	"log"
	"syscall"
)

func main() {
	config.Setup()
	models.Setup()
	logging.Setup()

	endless.DefaultReadTimeOut = config.ServerSetting.ReadTimeout
	endless.DefaultWriteTimeOut = config.ServerSetting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20

	endPoint := fmt.Sprintf(":%d", config.ServerSetting.HttpPort)
	server := endless.NewServer(endPoint, routers.InitRouter())
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}

	//router := routers.InitRouter()
	//
	//s := &http.Server{
	//	Addr:           fmt.Sprintf(":%d", config.HTTPPort),
	//	Handler:        router,
	//	ReadTimeout:    config.ReadTimeout,
	//	WriteTimeout:   config.WriteTimeout,
	//	MaxHeaderBytes: 1 << 20,
	//}
	//
	//s.ListenAndServe()
}
