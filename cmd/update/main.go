package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
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

func main() {
	repoFlag := flag.String("repo", "BretRen/liarsdeck", "GitHub 仓库 (owner/repo)")
	portFlag := flag.String("port", "8095", "服务监听端口")
	pidFlag := flag.Int("pid", 0, "当前运行的服务进程 PID")
	checkOnlyFlag := flag.Bool("check-only", false, "仅检查最新版本")
	flag.Parse()

	fmt.Println("==================================================")
	fmt.Println("🚀 Liar's Deck 自动更新与热替换程序")
	fmt.Printf("📦 目标仓库: %s | 系统: %s/%s\n", *repoFlag, runtime.GOOS, runtime.GOARCH)
	fmt.Println("==================================================")

	// 1. 获取 GitHub 最新 Release 信息
	release, err := fetchLatestRelease(*repoFlag)
	if err != nil {
		fmt.Printf("❌ 获取 GitHub 最新 Release 失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("🌟 最新发布版本: %s (%s)\n", release.TagName, release.Name)

	if *checkOnlyFlag {
		fmt.Println("✅ 检查版本完成。")
		return
	}

	// 2. 匹配适合当前操作系统与架构的 Asset 压缩包
	assetURL, assetName := findMatchingAsset(release.Assets, runtime.GOOS, runtime.GOARCH)
	if assetURL == "" {
		fmt.Printf("❌ 未找到适用于 %s-%s 的 Release 产物！\n", runtime.GOOS, runtime.GOARCH)
		fmt.Println("可用的文件列表:")
		for _, a := range release.Assets {
			fmt.Printf("  - %s\n", a.Name)
		}
		os.Exit(1)
	}

	fmt.Printf("📥 发现匹配资产包: %s\n🔗 下载链接: %s\n", assetName, assetURL)

	// 3. 创建临时工作目录并下载
	tmpDir, err := os.MkdirTemp("", "liarsdeck_update_*")
	if err != nil {
		fmt.Printf("❌ 创建临时目录失败: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	downloadPath := filepath.Join(tmpDir, assetName)
	fmt.Println("⏳ 正在下载最新更新包...")
	if err := downloadFile(assetURL, downloadPath); err != nil {
		fmt.Printf("❌ 下载失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ 下载完成！")

	// 4. 解压到暂存目录
	extractDir := filepath.Join(tmpDir, "unpacked")
	_ = os.MkdirAll(extractDir, 0755)
	fmt.Println("📦 正在解压安装包...")
	if strings.HasSuffix(assetName, ".zip") {
		err = unzip(downloadPath, extractDir)
	} else if strings.HasSuffix(assetName, ".tar.gz") || strings.HasSuffix(assetName, ".tgz") {
		err = untar(downloadPath, extractDir)
	} else {
		err = copyFile(downloadPath, filepath.Join(extractDir, getBinaryName()))
	}
	if err != nil {
		fmt.Printf("❌ 解压失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ 解压就绪！")

	// 5. 停止旧进程（如果指定了 PID）
	if *pidFlag > 0 {
		fmt.Printf("🛑 正在停止旧服务进程 (PID: %d)...\n", *pidFlag)
		killProcess(*pidFlag)
		time.Sleep(1 * time.Second)
	}

	// 6. 替换当前工作目录的文件
	cwd, _ := os.Getwd()
	fmt.Printf("📂 正在替换程序文件至: %s\n", cwd)

	binaryName := getBinaryName()
	srcBinary := filepath.Join(extractDir, binaryName)
	dstBinary := filepath.Join(cwd, binaryName)

	if fileExists(srcBinary) {
		_ = os.Chmod(srcBinary, 0755)
		if err := copyFile(srcBinary, dstBinary); err != nil {
			fmt.Printf("⚠️ 替换二进制文件失败 (可能需要重命名覆盖): %v\n", err)
			_ = os.Rename(dstBinary, dstBinary+".old."+strconv.FormatInt(time.Now().Unix(), 10))
			_ = copyFile(srcBinary, dstBinary)
		}
		_ = os.Chmod(dstBinary, 0755)
		fmt.Printf("✅ 主程序已替换: %s\n", binaryName)
	}

	// 替换 public 静态目录（如果存在）
	srcPublic := filepath.Join(extractDir, "public")
	dstPublic := filepath.Join(cwd, "public")
	if dirExists(srcPublic) {
		_ = copyDir(srcPublic, dstPublic)
		fmt.Println("✅ 前端资源 public/ 已更新！")
	}

	// 7. 启动新版服务进程
	fmt.Println("🚀 正在启动新版 Liar's Deck 服务...")
	cmd := exec.Command(dstBinary)
	cmd.Dir = cwd
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Printf("❌ 启动新服务失败: %v\n", err)
		os.Exit(1)
	}

	newPID := cmd.Process.Pid
	fmt.Printf("✨ 新服务已启动！PID: %d\n", newPID)

	// 8. 健康检查验证服务是否正常响应
	fmt.Println("🩺 正在进行服务健康检查...")
	healthURL := fmt.Sprintf("http://127.0.0.1:%s/", *portFlag)
	healthy := false

	for i := 0; i < 15; i++ {
		time.Sleep(800 * time.Millisecond)
		resp, err := http.Get(healthURL)
		if err == nil && resp.StatusCode == 200 {
			_ = resp.Body.Close()
			healthy = true
			break
		}
		if resp != nil {
			_ = resp.Body.Close()
		}
		fmt.Print(".")
	}
	fmt.Println()

	if healthy {
		fmt.Println("==================================================")
		fmt.Printf("🎉 更新成功！Liar's Deck (%s) 正在运行于端口 :%s\n", release.TagName, *portFlag)
		fmt.Println("==================================================")
	} else {
		fmt.Println("⚠️ 健康检查超时，但新进程已启动。请手动核对日志。")
	}
}

// ── 辅助函数 ──

func fetchLatestRelease(repo string) (*ReleaseInfo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "LiarsDeck-Updater")

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

func findMatchingAsset(assets []ReleaseAsset, goos, goarch string) (string, string) {
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

func downloadFile(url, dest string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "LiarsDeck-Updater")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败 HTTP %d", resp.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func getBinaryName() string {
	if runtime.GOOS == "windows" {
		return "liarsdeck.exe"
	}
	return "liarsdeck"
}

func killProcess(pid int) {
	proc, err := os.FindProcess(pid)
	if err == nil {
		_ = proc.Kill()
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func copyFile(src, dst string) error {
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

func copyDir(src, dst string) error {
	_ = os.MkdirAll(dst, 0755)
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			_ = copyDir(srcPath, dstPath)
		} else {
			_ = copyFile(srcPath, dstPath)
		}
	}
	return nil
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
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

func untar(src, dest string) error {
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
