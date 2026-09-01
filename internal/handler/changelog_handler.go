package handler

import (
	"io/fs"
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
	FS  fs.FS
}

func NewChangelogHandler(dir string, embeddedFS fs.FS) *ChangelogHandler {
	if dir == "" {
		dir = "changelogs"
	}
	return &ChangelogHandler{
		Dir: dir,
		FS:  embeddedFS,
	}
}

var safeVersionRegex = regexp.MustCompile(`^[vV]?[0-9]+(\.[0-9]+)*(-[a-zA-Z0-9.]+)?$`)

// GetChangelogList 返回所有可用更新日志版本列表（优先读取磁盘目录，回退至内置嵌入资源，按语义化版本降序排列）
func (h *ChangelogHandler) GetChangelogList(c echo.Context) error {
	versionMap := make(map[string]bool)

	// 1. 读取内置 embed.FS
	if h.FS != nil {
		_ = fs.WalkDir(h.FS, ".", func(path string, d fs.DirEntry, err error) error {
			if err == nil && !d.IsDir() && strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
				ver := strings.TrimSuffix(d.Name(), filepath.Ext(d.Name()))
				if safeVersionRegex.MatchString(ver) {
					versionMap[ver] = true
				}
			}
			return nil
		})
	}

	// 2. 读取本地磁盘目录（允许覆盖/扩充最新文件）
	if entries, err := os.ReadDir(h.Dir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if strings.HasSuffix(strings.ToLower(name), ".md") {
				ver := strings.TrimSuffix(name, filepath.Ext(name))
				if safeVersionRegex.MatchString(ver) {
					versionMap[ver] = true
				}
			}
		}
	}

	versions := make([]string, 0, len(versionMap))
	for v := range versionMap {
		versions = append(versions, v)
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

// GetChangelogContent 读取指定版本的更新日志 Markdown 原始内容（优先读取本地文件，回退至内置嵌入资源）
func (h *ChangelogHandler) GetChangelogContent(c echo.Context) error {
	version := strings.TrimSpace(c.Param("version"))
	if !safeVersionRegex.MatchString(version) {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "无效的版本号参数",
		})
	}

	// 尝试匹配可能的文件名变体 (带 v 或不带 v)
	candidates := []string{version, version + ".md"}
	if strings.HasPrefix(version, "v") {
		candidates = append(candidates, strings.TrimPrefix(version, "v"))
	} else {
		candidates = append(candidates, "v"+version)
	}

	var content []byte
	var readErr error
	matchedVer := version

	// 1. 优先尝试从本地磁盘读取
	for _, cand := range candidates {
		base := strings.TrimSuffix(cand, ".md")
		filePath := filepath.Join(h.Dir, base+".md")
		if data, err := os.ReadFile(filePath); err == nil {
			content = data
			matchedVer = base
			break
		}
	}

	// 2. 本地未找到时，尝试从内置 embed.FS 读取
	if len(content) == 0 && h.FS != nil {
		for _, cand := range candidates {
			base := strings.TrimSuffix(cand, ".md")
			// 尝试根路径或 changelogs 子路径
			for _, testPath := range []string{base + ".md", "changelogs/" + base + ".md"} {
				if data, err := fs.ReadFile(h.FS, testPath); err == nil {
					content = data
					matchedVer = base
					break
				}
			}
			if len(content) > 0 {
				break
			}
		}
	}

	if len(content) == 0 {
		return c.JSON(http.StatusNotFound, map[string]any{
			"success": false,
			"error":   "未找到指定版本的更新日志",
		})
	}

	_ = readErr
	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"version": matchedVer,
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
