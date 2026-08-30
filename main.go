package main

import (
	"log"
	"os"

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

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} | ${status} | ${latency_human} | ${method} ${uri}\n",
	}))
	e.Use(middleware.Recover())

	// WebSocket 对战接口
	e.GET("/ws", wsHandler.HandleWebSocket)

	// 静态文件与 SPA 页面托管
	e.Static("/", "public")
	e.File("/", "public/index.html")

	log.Printf("🃏 Liar's Deck 服务器已启动: http://localhost:%s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("服务器异常退出: %v", err)
	}
}
