package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/dyl0115/cup/scripts"
	"github.com/spf13/cobra"
)

var installEmail string

var installCmd = &cobra.Command{
	Use:   "install <domain>",
	Short: "nginx + HTTPS 세팅",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]

		if installEmail == "" {
			return fmt.Errorf("이메일이 필요합니다. --email 플래그를 사용해주세요.\n예시: cup install %s --email your@email.com", domain)
		}

		fmt.Println("📦 nginx + HTTPS 설치 중...")

		if _, err := exec.LookPath("nginx"); err == nil {
			fmt.Println("✅ nginx 이미 설치되어 있음. 스킵.")
			return nil
		}

		nginxTempConf, err := extractToTemp("nginx.temp.conf")
		if err != nil {
			return err
		}
		defer os.Remove(nginxTempConf)

		nginxConf, err := extractToTemp("nginx.conf")
		if err != nil {
			return err
		}
		defer os.Remove(nginxConf)

		if err := runEmbeddedScriptWithEnv("install.sh",
			[]string{domain, installEmail},
			[]string{
				"NGINX_TEMP_CONF=" + nginxTempConf,
				"NGINX_CONF=" + nginxConf,
			},
		); err != nil {
			return fmt.Errorf("설치 실패: %v", err)
		}

		fmt.Printf("✅ 설치 완료! 도메인: %s\n", domain)
		return nil
	},
}

func init() {
	installCmd.Flags().StringVarP(&installEmail, "email", "e", "", "Let's Encrypt 인증서 발급용 이메일 (필수)")
}

// extractToTemp는 embed된 파일을 임시파일로 추출하고 경로를 반환한다.
func extractToTemp(name string) (string, error) {
	data, err := scripts.FS.ReadFile(name)
	if err != nil {
		return "", fmt.Errorf("파일 로드 실패 [%s]: %v", name, err)
	}

	tmp, err := os.CreateTemp("", "cup-*-"+name)
	if err != nil {
		return "", fmt.Errorf("임시파일 생성 실패: %v", err)
	}

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return "", fmt.Errorf("파일 쓰기 실패: %v", err)
	}
	tmp.Close()

	if err := os.Chmod(tmp.Name(), 0755); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("chmod 실패: %v", err)
	}

	return tmp.Name(), nil
}

// runEmbeddedScript는 embed된 스크립트를 임시파일로 추출해서 실행한다.
func runEmbeddedScript(name string, args ...string) error {
	return runEmbeddedScriptWithEnv(name, args, nil)
}

// runEmbeddedScriptWithEnv는 환경변수를 추가로 주입해서 스크립트를 실행한다.
func runEmbeddedScriptWithEnv(name string, args []string, env []string) error {
	scriptPath, err := extractToTemp(name)
	if err != nil {
		return err
	}
	defer os.Remove(scriptPath)

	cmdArgs := append([]string{scriptPath}, args...)
	c := exec.Command("bash", cmdArgs...)
	c.Stdout = newPrefixWriter("  ")
	c.Stderr = newPrefixWriter("  ")
	c.Env = append(os.Environ(), env...)

	return c.Run()
}
