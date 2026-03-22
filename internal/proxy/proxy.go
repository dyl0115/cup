package proxy

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/dyl0115/cup/scripts"
)

// runScript는 embed된 스크립트를 임시파일로 추출해서 실행한다.
func runScript(name string, args ...string) error {
	data, err := scripts.FS.ReadFile(name)
	if err != nil {
		return fmt.Errorf("스크립트 로드 실패 [%s]: %v", name, err)
	}

	tmp, err := os.CreateTemp("", "cup-*.sh")
	if err != nil {
		return fmt.Errorf("임시파일 생성 실패: %v", err)
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.Write(data); err != nil {
		return fmt.Errorf("스크립트 쓰기 실패: %v", err)
	}
	tmp.Close()

	if err := os.Chmod(tmp.Name(), 0755); err != nil {
		return fmt.Errorf("chmod 실패: %v", err)
	}

	cmdArgs := append([]string{tmp.Name()}, args...)
	cmd := exec.Command("bash", cmdArgs...)
	cmd.Stdout = newPrefixWriter("  ")
	cmd.Stderr = newPrefixWriter("  ")
	return cmd.Run()
}

// Add는 nginx에 서비스를 등록한다.
func Add(id string, port int, path string) error {
	if err := runScript("add-server.sh", id, fmt.Sprintf("%d", port), path); err != nil {
		return fmt.Errorf("nginx 등록 실패: %v", err)
	}
	return nil
}

// Remove는 nginx에서 서비스를 제거한다.
func Remove(path string) error {
	if err := runScript("remove-server.sh", path); err != nil {
		return fmt.Errorf("nginx 제거 실패: %v", err)
	}
	return nil
}
