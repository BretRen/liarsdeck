package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"
)

type ChangelogHandler struct {
	Dir string
}

func NewChangelogHandler(dir string) *ChangelogHandler {
	if dir == "" {
		dir = "changelogs"
	}
	return &ChangelogHandler{Dir: dir}
}

var safeVersionRegex = regexp.MustCompile(`^[vV]?[0-9]+(\.[0-9]+)*(-[a-zA-Z0-9.]+)?$`)

// GetChangelogList 返回所有可用更新日志版本列表（按语义化版本降序排列）
func (h *ChangelogHandler) GetChangelogList(c echo.Context) error {
	entries, err := os.ReadDir(h.Dir)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]any{
			"success":  true,
			"versions": []string{},
		})
	}

	versions := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(strings.ToLower(name), ".md") {
			ver := strings.TrimSuffix(name, filepath.Ext(name))
			if safeVersionRegex.MatchString(ver) {
				versions = append(versions, ver)
			}
		}
	}

	// 语义化倒序排列（新版本在前）
	sort.Slice(versions, func(i, j int) bool {
		return compareSemver(versions[i], versions[j]) > 0
	})

	return c.JSON(http.StatusOK, map[string]any{
		"success":  true,
		"versions": versions,
	})
}

// GetChangelogContent 读取指定版本的更新日志 Markdown 原始内容
func (h *ChangelogHandler) GetChangelogContent(c echo.Context) error {
	version := strings.TrimSpace(c.Param("version"))
	if !safeVersionRegex.MatchString(version) {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "无效的版本号参数",
		})
	}

	filename := filepath.Join(h.Dir, version+".md")
	content, err := os.ReadFile(filename)
	if err != nil {
		// 尝试补全或去掉 v 前缀
		altVer := version
		if strings.HasPrefix(version, "v") {
			altVer = strings.TrimPrefix(version, "v")
		} else {
			altVer = "v" + version
		}
		content, err = os.ReadFile(filepath.Join(h.Dir, altVer+".md"))
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]any{
				"success": false,
				"error":   "未找到指定版本的更新日志",
			})
		}
		version = altVer
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"version": version,
		"content": string(content),
	})
}

// compareSemver 简易语义化版本比对
func compareSemver(v1, v2 string) int {
	clean1 := strings.TrimPrefix(strings.ToLower(v1), "v")
	clean2 := strings.TrimPrefix(strings.ToLower(v2), "v")

	parts1 := strings.Split(clean1, ".")
	parts2 := strings.Split(clean2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		p1 := 0
		p2 := 0
		if i < len(parts1) {
			p1 = parseLeadingInt(parts1[i])
		}
		if i < len(parts2) {
			p2 = parseLeadingInt(parts2[i])
		}
		if p1 != p2 {
			return p1 - p2
		}
	}
	return strings.Compare(clean1, clean2)
}

func parseLeadingInt(s string) int {
	res := 0
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			res = res*10 + int(ch-'0')
		} else {
			break
		}
	}
	return res
}
