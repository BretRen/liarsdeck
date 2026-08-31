package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"pdnode.com/play/liarsbar-web/internal/room"
	"pdnode.com/play/liarsbar-web/internal/updater"
)

// Version 服务端运行版本号，支持在编译时通过 -ldflags "-X 'pdnode.com/play/liarsbar-web/internal/handler.Version=v...'" 动态注入
var Version = "dev"

// GetVersion 动态获取当前运行版本（优先使用编译时注入的 Version，其次读取 git tag / buildinfo）
func GetVersion() string {
	if Version != "" && Version != "dev" && Version != "v2.0.6" {
		return Version
	}
	if bi, ok := debug.ReadBuildInfo(); ok {
		if bi.Main.Version != "" && bi.Main.Version != "(devel)" {
			return bi.Main.Version
		}
		for _, s := range bi.Settings {
			if s.Key == "vcs.revision" && len(s.Value) >= 7 {
				return "dev-" + s.Value[:7]
			}
		}
	}
	// 开发态回退：从本地 git 获取最新 tag
	cmd := exec.Command("git", "describe", "--tags", "--always")
	if out, err := cmd.Output(); err == nil {
		tag := strings.TrimSpace(string(out))
		if tag != "" {
			return tag
		}
	}
	return "v2.3.0"
}

type AdminHandler struct {
	Hub    *room.Hub
	Secret string
	Port   string
}

func NewAdminHandler(hub *room.Hub, port string) *AdminHandler {
	secret := os.Getenv("ADMIN_SECRET")
	return &AdminHandler{
		Hub:    hub,
		Secret: secret,
		Port:   port,
	}
}

type AdminAuthPayload struct {
	Secret string `json:"secret"`
}

func (h *AdminHandler) checkAuth(c echo.Context) (bool, error) {
	if h.Secret == "" {
		return false, echo.NewHTTPError(http.StatusForbidden, "服务端未配置 ADMIN_SECRET 环境变量，管理员接口已安全禁用 / Admin API is disabled: ADMIN_SECRET environment variable is not set")
	}
	var req AdminAuthPayload
	if err := c.Bind(&req); err != nil {
		return false, err
	}
	return req.Secret == h.Secret, nil
}

func (h *AdminHandler) Auth(c echo.Context) error {
	if h.Secret == "" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "服务端未配置 ADMIN_SECRET 环境变量，管理员接口已安全禁用 / Admin API disabled: ADMIN_SECRET is not set"})
	}
	ok, err := h.checkAuth(c)
	if err != nil || !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "管理密钥错误 / Invalid admin secret"})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"authenticated": true,
		"version":       GetVersion(),
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

	currVer := GetVersion()

	if resp.StatusCode == http.StatusNotFound {
		return c.JSON(http.StatusOK, map[string]any{
			"current_version": currVer,
			"latest_version":  "暂无 Release",
			"has_update":      false,
			"release_name":    "尚未发布 Release 版本",
			"release_body":    "GitHub 仓库 (BretRen/liarsdeck) 暂未发布任何 Release。\n推送版本 Tag (例如: git tag v2.0.0 && git push origin v2.0.0) 后 GitHub Actions 将自动构建并发布。",
		})
	}

	if resp.StatusCode != http.StatusOK {
		return c.JSON(http.StatusOK, map[string]any{
			"current_version": currVer,
			"latest_version":  "—",
			"has_update":      false,
			"release_name":    fmt.Sprintf("GitHub 响应异常 (HTTP %d)", resp.StatusCode),
			"release_body":    "可能触发了 GitHub API 访问速率限制，请稍后重试。",
		})
	}

	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "解析 GitHub 数据失败"})
	}

	latestTag, _ := data["tag_name"].(string)
	hasUpdate := latestTag != "" && latestTag != currVer

	return c.JSON(http.StatusOK, map[string]any{
		"current_version": currVer,
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

	// 启动纯 Go 原生后台自更新流程，无需系统具备 go 编译器或外部工具
	go func() {
		if err := updater.PerformSelfUpdate("BretRen/liarsdeck", h.Port); err != nil {
			fmt.Printf("❌ [Updater] 自更新执行失败: %v\n", err)
		}
	}()

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"status":  "started",
		"message": "已在后台启动下载并热替换更新，预计 3~5 秒内自动重启完成！",
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

func (h *AdminHandler) Broadcast(c echo.Context) error {
	if h.Secret == "" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "服务端未配置 ADMIN_SECRET 环境变量，管理员接口已安全禁用"})
	}

	var req struct {
		Secret  string `json:"secret"`
		Message string `json:"message"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "请求参数解析失败"})
	}

	if req.Secret != h.Secret {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "管理密钥错误 / Invalid admin secret"})
	}

	msg := strings.TrimSpace(req.Message)
	if msg == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "广播内容不能为空"})
	}

	roomCount := h.Hub.BroadcastGlobal("server_broadcast", map[string]any{
		"message":   msg,
		"timestamp": time.Now().Unix(),
	})

	return c.JSON(http.StatusOK, map[string]any{
		"success":    true,
		"status":     "ok",
		"message":    fmt.Sprintf("广播已成功推送至全服 %d 个房间！", roomCount),
		"room_count": roomCount,
	})
}
