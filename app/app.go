package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/orbit-w/ogateway/app/gateway"
)

func Run() {
	stopper, err := gateway.Serve()
	if err != nil {
		panic(err)
	}

	gracefulShutdown(stopper)
}

// gracefulShutdown 优雅关闭服务
func gracefulShutdown(stopper func(ctx context.Context) error) {
	// 等待中断信号
	quit := make(chan os.Signal, 1)
	// 监听 SIGINT（Ctrl+C）和 SIGTERM 信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 创建一个5分钟超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if stopper != nil {
		if err := stopper(ctx); err != nil {
			log.Printf("Error stopping stopper: %v", err)
		}
	}

	log.Println("Server exiting")
}
