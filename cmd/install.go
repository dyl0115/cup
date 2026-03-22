package cmd

import (
	"fmt"
	"os/exec"
)

const proxyDir = "/opt/dy-proxy-server"

func Install(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("사용법: cup install <domain>")
	}
	domain := args[0]

	fmt.Println("📦 dy-proxy-server 설치 중...")

	// 이미 설치되어 있으면 스킵
	if _, err := exec.LookPath("nginx"); err == nil {
		fmt.Println("✅ nginx 이미 설치되어 있음. 스킵.")
		return nil
	}

	// dy-proxy-server 클론
	clone := exec.Command("bash", "-c",
		fmt.Sprintf("git clone https://github.com/dyl0115/dy-proxy-server.git %s", proxyDir))
	clone.Stdout = newPrefixWriter("  ")
	clone.Stderr = newPrefixWriter("  ")
	if err := clone.Run(); err != nil {
		return fmt.Errorf("레포 클론 실패: %v", err)
	}

	// 실행 권한 부여
	chmod := exec.Command("bash", "-c", fmt.Sprintf("chmod +x %s/*.sh", proxyDir))
	if err := chmod.Run(); err != nil {
		return fmt.Errorf("chmod 실패: %v", err)
	}

	// install.sh 실행
	install := exec.Command("bash", fmt.Sprintf("%s/install.sh", proxyDir), domain)
	install.Stdout = newPrefixWriter("  ")
	install.Stderr = newPrefixWriter("  ")
	if err := install.Run(); err != nil {
		return fmt.Errorf("install.sh 실패: %v", err)
	}

	fmt.Printf("✅ 설치 완료! 도메인: %s\n", domain)
	return nil
}
