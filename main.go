package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"pdnode.com/play/liarsbar-web/internal/handler"
	"pdnode.com/play/liarsbar-web/internal/room"
)

//go:embed all:public
var embeddedPublic embed.FS

//go:embed changelogs/*.md
var embeddedChangelogs embed.FS

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
	changelogHandler := handler.NewChangelogHandler("changelogs", embeddedChangelogs)

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
	adminGroup.POST("/broadcast", adminHandler.Broadcast)
	adminGroup.POST("/stats", adminHandler.GetStats)

	// 公开版本查询 API
	e.GET("/api/version", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"version": handler.GetVersion(),
		})
	})

	// 更新日志 API
	e.GET("/api/changelogs", changelogHandler.GetChangelogList)
	e.GET("/api/changelogs/:version", changelogHandler.GetChangelogContent)

	// 公开大厅房间列表 API
	e.GET("/api/rooms", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"success": true,
			"rooms":   hub.GetPublicRooms(),
		})
	})

	// 静态文件与 SPA 页面托管：内置嵌入与本地开发双轨支持
	publicFS, err := fs.Sub(embeddedPublic, "public")
	if err != nil {
		log.Fatalf("无法解析内置前端资源: %v", err)
	}

	hasLocalPublic := false
	if fi, err := os.Stat("public/index.html"); err == nil && !fi.IsDir() {
		hasLocalPublic = true
	}

	var assetFS http.FileSystem
	if hasLocalPublic {
		assetFS = http.Dir("public")
	} else {
		assetFS = http.FS(publicFS)
	}
	fileServer := http.FileServer(assetFS)

	serveIndex := func(c echo.Context) error {
		if hasLocalPublic {
			return c.File("public/index.html")
		}
		indexHTML, err := fs.ReadFile(publicFS, "index.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Internal Error: index.html not found in binary")
		}
		return c.HTMLBlob(http.StatusOK, indexHTML)
	}

	e.GET("/assets/*", echo.WrapHandler(fileServer))
	e.GET("/favicon.ico", func(c echo.Context) error {
		if hasLocalPublic {
			if _, err := os.Stat("public/favicon.ico"); err == nil {
				return c.File("public/favicon.ico")
			}
		}
		if data, err := fs.ReadFile(publicFS, "favicon.ico"); err == nil {
			return c.Blob(http.StatusOK, "image/x-icon", data)
		}
		return c.NoContent(http.StatusNotFound)
	})

	e.GET("/", serveIndex)
	e.GET("/callback", serveIndex)
	e.RouteNotFound("/*", func(c echo.Context) error {
		p := strings.TrimPrefix(c.Request().URL.Path, "/")
		if p != "" {
			if hasLocalPublic {
				if _, err := os.Stat(filepath.Join("public", p)); err == nil {
					return c.File(filepath.Join("public", p))
				}
			} else {
				if f, err := publicFS.Open(p); err == nil {
					_ = f.Close()
					fileServer.ServeHTTP(c.Response(), c.Request())
					return nil
				}
			}
		}
		return serveIndex(c)
	})

	log.Printf("🃏 Liar's Deck (%s) 服务器已启动: http://localhost:%s", handler.GetVersion(), port)
	if os.Getenv("ADMIN_SECRET") != "" {
		log.Printf("🛡️ 管理员接口已启用 (ADMIN_SECRET 已配置)")
	} else {
		log.Printf("ℹ️ 未配置 ADMIN_SECRET 环境变量，管理员接口 (/api/admin/*) 已安全禁用")
	}

	if err := e.Start(":" + port); err != nil {
		log.Fatalf("服务器异常退出: %v", err)
	}
}
