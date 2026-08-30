package main

import (
	"log"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"pdnode.com/play/liarsbar-web/internal/handler"
	"pdnode.com/play/liarsbar-web/internal/room"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8095"
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
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("服务器异常退出: %v", err)
	}
}
