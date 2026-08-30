package main

import (
	"flag"
	"log"

	"pdnode.com/play/liarsbar-web/internal/updater"
)

func main() {
	repoFlag := flag.String("repo", "BretRen/liarsdeck", "GitHub 仓库 (owner/repo)")
	portFlag := flag.String("port", "8095", "服务监听端口")
	flag.Parse()

	if err := updater.PerformSelfUpdate(*repoFlag, *portFlag); err != nil {
		log.Fatalf("❌ 更新失败: %v", err)
	}
}
