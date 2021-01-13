package main

import (
	"GoGin/ginpackage/pkg/logging"
	"GoGin/ginpackage/pkg/setting"
	"GoGin/ginpackage/routers"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// 使用endless进行监控
// func main() {

// 	endless.DefaultReadTimeOut = setting.ReadTimeout
// 	endless.DefaultWriteTimeOut = setting.WriteTimeout
// 	endless.DefaultMaxHeaderBytes = 1 << 20
// 	endPoint := fmt.Sprintf(":%d", setting.HTTPPort)

// 	server := endless.NewServer(endPoint, routers.InitRouter())
// 	server.BeforeBegin = func(add string) {
// 		log.Printf("Actual pid is %d", syscall.Getpid())
// 	}

// 	err := server.ListenAndServe()
// 	if err != nil {
// 		log.Printf("Server err:%v", err)
// 	}
// }

// 使用http.Server的Shutdown方法
func main() {
	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := s.ListenAndServe(); err != nil {
			logging.Error("服务启动失败。监听报错：", err)
		}

	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		logging.Fatal("Server Shutdown:", err)
	}
	logging.Info("Server exiting")
}
