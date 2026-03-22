package proxy

import (
	"fmt"
	"os/exec"
)

const proxyDir = "/opt/dy-proxy-server"

// Add는 nginx에 서비스를 등록한다.
func Add(id string, port int, path string) error {
	script := fmt.Sprintf("%s/add-server.sh", proxyDir)
	cmd := exec.Command("bash", script, id, fmt.Sprintf("%d", port), path)
	cmd.Stdout = newPrefixWriter("  ")
	cmd.Stderr = newPrefixWriter("  ")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("nginx 등록 실패: %v", err)
	}
	return nil
}

// Remove는 nginx에서 서비스를 제거한다.
func Remove(path string) error {
	script := fmt.Sprintf("%s/remove-server.sh", proxyDir)
	cmd := exec.Command("bash", script, path)
	cmd.Stdout = newPrefixWriter("  ")
	cmd.Stderr = newPrefixWriter("  ")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("nginx 제거 실패: %v", err)
	}
	return nil
}
