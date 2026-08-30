package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"pdnode.com/play/liarsbar-web/internal/handler"
	"pdnode.com/play/liarsbar-web/internal/room"
)

func main() {
	waitPID := flag.Int("wait-pid", 0, "等待旧服务进程退出并释放端口后再绑定")
	flag.Parse()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8095"
	}

	// 如果传入了 --wait-pid，说明是由旧版本自更新拉起的新版本，等待旧进程释放端口
	if *waitPID > 0 {
		log.Printf("⏳ [AutoRestart] 正在等待旧服务进程 (PID: %d) 退出并释放端口 :%s...", *waitPID, port)
		for i := 0; i < 50; i++ {
			time.Sleep(100 * time.Millisecond)
			l, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
			if err == nil {
				_ = l.Close()
				break
			}
		}
		time.Sleep(200 * time.Millisecond)
		log.Printf("✅ [AutoRestart] 端口 :%s 已就绪，新版本服务开始监听！", port)
	}

	hub := room.NewHub()
	wsHandler := handler.NewWSHandler(hub)
	adminHandler := handler.NewAdminHandler(hub, port)

	e := echo.New()
	e.HideBanner = true

	// 使用现代且推荐的 RequestLoggerWithConfig 替代已废弃的 LoggerWithConfig
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		LogLatency:  true,
		LogMethod:   true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error != nil {
				log.Printf("%s | %d | %v | %s %s | error: %v",
					v.StartTime.Format(time.RFC3339),
					v.Status,
					v.Latency,
					v.Method,
					v.URI,
					v.Error,
				)
			} else {
				log.Printf("%s | %d | %v | %s %s",
					v.StartTime.Format(time.RFC3339),
					v.Status,
					v.Latency,
					v.Method,
					v.URI,
				)
			}
			return nil
		},
	}))
	e.Use(middleware.Recover())

	// WebSocket 对战接口
	e.GET("/ws", wsHandler.HandleWebSocket)

	// 管理员与自动更新 API
	adminGroup := e.Group("/api/admin")
	adminGroup.POST("/auth", adminHandler.Auth)
	adminGroup.POST("/check-update", adminHandler.CheckUpdate)
	adminGroup.POST("/trigger-update", adminHandler.TriggerUpdate)
	adminGroup.POST("/stats", adminHandler.GetStats)

	// 静态文件与 SPA 页面托管
	e.Static("/", "public")
	e.File("/", "public/index.html")

	log.Printf("🃏 Liar's Deck 服务器已启动: http://localhost:%s", port)
	if os.Getenv("ADMIN_SECRET") != "" {
		log.Printf("🛡️ 管理员接口已启用 (ADMIN_SECRET 已配置)")
	} else {
		log.Printf("ℹ️ 未配置 ADMIN_SECRET 环境变量，管理员接口 (/api/admin/*) 已安全禁用")
	}

	if err := e.Start(":" + port); err != nil {
		log.Fatalf("服务器异常退出: %v", err)
	}
}
