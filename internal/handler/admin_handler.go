package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/labstack/echo/v4"
	"pdnode.com/play/liarsbar-web/internal/room"
)

type AdminHandler struct {
	Hub     *room.Hub
	Secret  string
	Port    string
	Version string
}

func NewAdminHandler(hub *room.Hub, port string) *AdminHandler {
	secret := os.Getenv("ADMIN_SECRET")
	if secret == "" {
		secret = "liarsbar2026"
	}
	return &AdminHandler{
		Hub:     hub,
		Secret:  secret,
		Port:    port,
		Version: "v2.0.0",
	}
}

type AdminAuthPayload struct {
	Secret string `json:"secret"`
}

func (h *AdminHandler) checkAuth(c echo.Context) (bool, error) {
	var req AdminAuthPayload
	if err := c.Bind(&req); err != nil {
		return false, err
	}
	return req.Secret == h.Secret, nil
}

func (h *AdminHandler) Auth(c echo.Context) error {
	ok, err := h.checkAuth(c)
	if err != nil || !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "管理密钥错误 / Invalid admin secret"})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"authenticated": true,
		"version":       h.Version,
	})
}

func (h *AdminHandler) CheckUpdate(c echo.Context) error {
	ok, err := h.checkAuth(c)
	if err != nil || !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "管理密钥错误 / Invalid admin secret"})
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/BretRen/liarsdeck/releases/latest", nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	req.Header.Set("User-Agent", "LiarsDeck-Server")

	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusBadGateway, map[string]string{"error": "无法连接 GitHub: " + err.Error()})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.JSON(resp.StatusCode, map[string]string{"error": "GitHub API 返回错误"})
	}

	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "解析 GitHub 数据失败"})
	}

	latestTag, _ := data["tag_name"].(string)
	hasUpdate := latestTag != "" && latestTag != h.Version

	return c.JSON(http.StatusOK, map[string]any{
		"current_version": h.Version,
		"latest_version":  latestTag,
		"has_update":      hasUpdate,
		"release_name":    data["name"],
		"release_body":    data["body"],
		"published_at":    data["published_at"],
	})
}

func (h *AdminHandler) TriggerUpdate(c echo.Context) error {
	ok, err := h.checkAuth(c)
	if err != nil || !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "管理密钥错误 / Invalid admin secret"})
	}

	pid := os.Getpid()
	binaryName := "update"
	if runtime.GOOS == "windows" {
		binaryName = "update.exe"
	}

	// 如果当前目录有 update 二进制，则调用；若无则通过 go run ./cmd/update 调用
	var cmd *exec.Cmd
	if _, err := os.Stat(binaryName); err == nil {
		cmd = exec.Command("./"+binaryName, fmt.Sprintf("--pid=%d", pid), fmt.Sprintf("--port=%s", h.Port))
	} else {
		cmd = exec.Command("go", "run", "./cmd/update", fmt.Sprintf("--pid=%d", pid), fmt.Sprintf("--port=%s", h.Port))
	}

	// 后台分离启动
	if err := cmd.Start(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "启动更新程序失败: " + err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status":  "started",
		"message": "更新程序已在后台拉起，稍后将自动替换并重启服务！",
	})
}

func (h *AdminHandler) GetStats(c echo.Context) error {
	ok, err := h.checkAuth(c)
	if err != nil || !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "管理密钥错误 / Invalid admin secret"})
	}

	h.Hub.Lock()
	defer h.Hub.Unlock()

	totalRooms := len(h.Hub.Rooms)
	totalPlayers := 0
	for _, r := range h.Hub.Rooms {
		r.Game.Lock()
		for _, p := range r.Game.State.Players {
			if !p.IsSpectator {
				totalPlayers++
			}
		}
		r.Game.Unlock()
	}

	return c.JSON(http.StatusOK, map[string]any{
		"total_rooms":   totalRooms,
		"total_players": totalPlayers,
		"current_time":  time.Now().Format("2006-01-02 15:04:05"),
	})
}
