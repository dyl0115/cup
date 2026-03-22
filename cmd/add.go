package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dyl0115/cup/internal/cupfile"
	"github.com/dyl0115/cup/internal/executor"
	"github.com/dyl0115/cup/internal/proxy"
	"github.com/dyl0115/cup/internal/registry"
	"github.com/spf13/cobra"
)

const cloneBaseDir = "/opt/cup-services"

var addCmd = &cobra.Command{
	Use:   "add <github-url 또는 로컬경로>",
	Short: "cupFile.yaml 읽고 서비스 등록",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]

		var absDir string

		if isGitHubURL(input) {
			// GitHub URL이면 자동으로 clone
			dir, err := cloneRepo(input)
			if err != nil {
				return err
			}
			absDir = dir
		} else {
			// 로컬 경로면 그대로 사용
			dir, err := filepath.Abs(input)
			if err != nil {
				return fmt.Errorf("경로 변환 실패: %v", err)
			}
			absDir = dir
		}

		fmt.Printf("📂 %s/cupFile.yaml 읽는 중...\n", absDir)
		cf, err := cupfile.Load(absDir)
		if err != nil {
			return err
		}

		fmt.Printf("🚀 [%s] INIT 실행 중...\n", cf.ID)
		if err := executor.Run(cf.Init); err != nil {
			return err
		}

		fmt.Printf("🌐 nginx에 등록 중... (port: %d, path: %s)\n", cf.Port, cf.Path)
		if err := proxy.Add(cf.ID, cf.Port, cf.Path); err != nil {
			return err
		}

		if err := registry.Save(cf.ID, absDir); err != nil {
			return fmt.Errorf("registry 저장 실패: %v", err)
		}

		fmt.Printf("✅ [%s] 등록 완료! (%s)\n", cf.ID, absDir)
		return nil
	},
}

// isGitHubURL은 입력이 GitHub URL인지 확인한다.
func isGitHubURL(input string) bool {
	return strings.HasPrefix(input, "https://github.com") ||
		strings.HasPrefix(input, "git@github.com")
}

// cloneRepo는 GitHub URL을 /opt/cup-services/<레포명> 으로 clone한다.
// 이미 존재하면 git pull로 업데이트한다.
func cloneRepo(url string) (string, error) {
	// URL에서 레포명 추출 (예: https://github.com/dyl0115/my-go-server → my-go-server)
	repoName := extractRepoName(url)
	if repoName == "" {
		return "", fmt.Errorf("GitHub URL에서 레포명을 추출할 수 없습니다: %s", url)
	}

	destDir := filepath.Join(cloneBaseDir, repoName)

	if err := os.MkdirAll(cloneBaseDir, 0755); err != nil {
		return "", fmt.Errorf("디렉토리 생성 실패: %v", err)
	}

	if _, err := os.Stat(destDir); err == nil {
		// 이미 존재하면 git pull
		fmt.Printf("📦 [%s] 이미 존재함. git pull 중...\n", repoName)
		c := exec.Command("git", "-C", destDir, "pull")
		c.Stdout = newPrefixWriter("  ")
		c.Stderr = newPrefixWriter("  ")
		if err := c.Run(); err != nil {
			return "", fmt.Errorf("git pull 실패: %v", err)
		}
	} else {
		// 없으면 git clone
		fmt.Printf("📦 [%s] git clone 중...\n", repoName)
		c := exec.Command("git", "clone", url, destDir)
		c.Stdout = newPrefixWriter("  ")
		c.Stderr = newPrefixWriter("  ")
		if err := c.Run(); err != nil {
			return "", fmt.Errorf("git clone 실패: %v", err)
		}
	}

	return destDir, nil
}

// extractRepoName은 GitHub URL에서 레포명을 추출한다.
// https://github.com/dyl0115/my-go-server     → my-go-server
// https://github.com/dyl0115/my-go-server.git → my-go-server
// git@github.com:dyl0115/my-go-server.git     → my-go-server
func extractRepoName(url string) string {
	// .git 제거
	url = strings.TrimSuffix(url, ".git")
	// 마지막 / 이후가 레포명
	parts := strings.Split(url, "/")
	if len(parts) == 0 {
		return ""
	}
	// git@github.com:dyl0115/my-go-server 형태 처리
	last := parts[len(parts)-1]
	if colonIdx := strings.LastIndex(last, ":"); colonIdx != -1 {
		last = last[colonIdx+1:]
	}
	return last
}
