package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/dyl0115/cup/scripts"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install <domain>",
	Short: "nginx + HTTPS 세팅",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]

		fmt.Println("📦 nginx + HTTPS 설치 중...")

		if _, err := exec.LookPath("nginx"); err == nil {
			fmt.Println("✅ nginx 이미 설치되어 있음. 스킵.")
			return nil
		}

		if err := runEmbeddedScript("install.sh", domain); err != nil {
			return fmt.Errorf("설치 실패: %v", err)
		}

		fmt.Printf("✅ 설치 완료! 도메인: %s\n", domain)
		return nil
	},
}

func runEmbeddedScript(name string, args ...string) error {
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
		return fmt.Errorf("스크립트 쓰기 실패: %Wv", err)
	}
	tmp.Close()

	if err := os.Chmod(tmp.Name(), 0755); err != nil {
		return fmt.Errorf("chmod 실패: %v", err)
	}

	cmdArgs := append([]string{tmp.Name()}, args...)
	c := exec.Command("bash", cmdArgs...)
	c.Stdout = newPrefixWriter("  ")
	c.Stderr = newPrefixWriter("  ")
	return c.Run()
}
