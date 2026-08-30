package updater

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type ReleaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type ReleaseInfo struct {
	TagName string         `json:"tag_name"`
	Name    string         `json:"name"`
	Body    string         `json:"body"`
	Assets  []ReleaseAsset `json:"assets"`
}

// PerformSelfUpdate 执行跨平台（Linux & Windows）安全热替换与无缝重启
func PerformSelfUpdate(repo string, port string) error {
	log.Printf("🚀 [Updater] 启动后台自更新流程 (目标仓库: %s, 平台: %s/%s)...", repo, runtime.GOOS, runtime.GOARCH)

	// 1. 获取 GitHub 最新发布信息
	release, err := FetchLatestRelease(repo)
	if err != nil {
		return fmt.Errorf("获取 GitHub Release 失败: %w", err)
	}

	assetURL, assetName := FindMatchingAsset(release.Assets, runtime.GOOS, runtime.GOARCH)
	if assetURL == "" {
		return fmt.Errorf("未找到适用于 %s-%s 的 Release 资产包", runtime.GOOS, runtime.GOARCH)
	}

	log.Printf("📥 [Updater] 发现新版本资产包: %s (%s)", assetName, release.TagName)

	// 2. 创建临时目录并下载
	tmpDir, err := os.MkdirTemp("", "liarsdeck_update_*")
	if err != nil {
		return fmt.Errorf("创建临时目录失败: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	downloadPath := filepath.Join(tmpDir, assetName)
	log.Printf("⏳ [Updater] 正在下载: %s ...", assetURL)
	if err := DownloadFile(assetURL, downloadPath); err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}
	log.Printf("✅ [Updater] 资产包下载完成")

	// 3. 解压
	extractDir := filepath.Join(tmpDir, "unpacked")
	_ = os.MkdirAll(extractDir, 0755)
	if strings.HasSuffix(assetName, ".zip") {
		err = Unzip(downloadPath, extractDir)
	} else if strings.HasSuffix(assetName, ".tar.gz") || strings.HasSuffix(assetName, ".tgz") {
		err = Untar(downloadPath, extractDir)
	} else {
		err = CopyFile(downloadPath, filepath.Join(extractDir, GetBinaryName()))
	}
	if err != nil {
		return fmt.Errorf("解压资产包失败: %w", err)
	}

	// 4. 定位当前工作目录与可执行文件路径
	cwd, _ := os.Getwd()
	execPath, err := os.Executable()
	if err != nil {
		execPath = filepath.Join(cwd, GetBinaryName())
	}

	binaryName := GetBinaryName()
	srcBinary := filepath.Join(extractDir, binaryName)
	if !FileExists(srcBinary) {
		// 尝试在解压子目录中寻找
		_ = filepath.Walk(extractDir, func(p string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && info.Name() == binaryName {
				srcBinary = p
			}
			return nil
		})
	}

	if !FileExists(srcBinary) {
		return fmt.Errorf("解压内容中未找到二进制可执行文件: %s", binaryName)
	}

	// 5. 备份旧二进制 (通过重命名机制：Linux & Windows 均支持在运行状态下将文件改名挪位)
	oldBackupPath := execPath + ".old." + strconv.FormatInt(time.Now().Unix(), 10)
	log.Printf("📂 [Updater] 正在将当前运行程序重命名让位至: %s", oldBackupPath)
	if err := os.Rename(execPath, oldBackupPath); err != nil {
		log.Printf("⚠️ [Updater] 重命名旧文件失败 (尝试直接覆盖): %v", err)
	}

	// 写入新版本二进制文件
	_ = os.Chmod(srcBinary, 0755)
	if err := CopyFile(srcBinary, execPath); err != nil {
		// 回滚
		_ = os.Rename(oldBackupPath, execPath)
		return fmt.Errorf("写入新二进制文件失败: %w", err)
	}
	_ = os.Chmod(execPath, 0755)

	// 更新 public 静态资源（若解压目录存在）
	srcPublic := filepath.Join(extractDir, "public")
	if !DirExists(srcPublic) {
		_ = filepath.Walk(extractDir, func(p string, info os.FileInfo, err error) error {
			if err == nil && info.IsDir() && info.Name() == "public" {
				srcPublic = p
			}
			return nil
		})
	}
	if DirExists(srcPublic) {
		dstPublic := filepath.Join(cwd, "public")
		_ = CopyDir(srcPublic, dstPublic)
		log.Printf("✅ [Updater] 前端静态资源 public/ 已同步替换")
	}

	// 6. 启动新版本子进程（传入 --wait-pid 确保旧进程退出并释放端口后再绑定）
	log.Printf("🚀 [Updater] 正在启动新版本服务进程: %s ...", execPath)
	cmd := exec.Command(execPath, fmt.Sprintf("--wait-pid=%d", os.Getpid()))
	cmd.Dir = cwd
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		// 启动失败，回滚二进制
		_ = os.Rename(oldBackupPath, execPath)
		return fmt.Errorf("启动新版本进程失败: %w", err)
	}

	newPID := cmd.Process.Pid
	log.Printf("✨ [Updater] 新版服务已拉起 (PID: %d)，旧进程即将退出以释放端口 :%s ...", newPID, port)
	time.Sleep(500 * time.Millisecond)
	os.Exit(0)
	return nil
}

// ── 辅助函数 ──

func FetchLatestRelease(repo string) (*ReleaseInfo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "LiarsDeck-Server")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var info ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

func FindMatchingAsset(assets []ReleaseAsset, goos, goarch string) (string, string) {
	keywords := []string{
		fmt.Sprintf("%s-%s", goos, goarch),
		fmt.Sprintf("%s_%s", goos, goarch),
		goos,
	}

	for _, kw := range keywords {
		for _, a := range assets {
			nameLower := strings.ToLower(a.Name)
			if strings.Contains(nameLower, strings.ToLower(kw)) {
				return a.BrowserDownloadURL, a.Name
			}
		}
	}

	if len(assets) > 0 {
		return assets[0].BrowserDownloadURL, assets[0].Name
	}
	return "", ""
}

func DownloadFile(url, dest string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "LiarsDeck-Server")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载返回 HTTP %d", resp.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func GetBinaryName() string {
	if runtime.GOOS == "windows" {
		return "liarsdeck.exe"
	}
	return "liarsdeck"
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func DirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func CopyDir(src, dst string) error {
	_ = os.MkdirAll(dst, 0755)
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			_ = CopyDir(srcPath, dstPath)
		} else {
			_ = CopyFile(srcPath, dstPath)
		}
	}
	return nil
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	destClean := filepath.Clean(dest) + string(os.PathSeparator)

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		cleanPath := filepath.Clean(fpath)

		// 防止 ZipSlip 路径穿越安全漏洞
		if !strings.HasPrefix(cleanPath+string(os.PathSeparator), destClean) && cleanPath != filepath.Clean(dest) {
			return fmt.Errorf("非法压缩包文件路径: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		_ = os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func Untar(src, dest string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	destClean := filepath.Clean(dest) + string(os.PathSeparator)

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dest, header.Name)
		cleanTarget := filepath.Clean(target)

		// 防止 ZipSlip 路径穿越安全漏洞
		if !strings.HasPrefix(cleanTarget+string(os.PathSeparator), destClean) && cleanTarget != filepath.Clean(dest) {
			return fmt.Errorf("非法 tar 文件路径: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			_ = os.MkdirAll(target, 0755)
		case tar.TypeReg:
			_ = os.MkdirAll(filepath.Dir(target), 0755)
			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}
	return nil
}
