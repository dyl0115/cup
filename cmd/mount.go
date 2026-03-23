package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var mountOnedrivePath string

var mountCmd = &cobra.Command{
	Use:   "mount",
	Short: "클라우드 스토리지 마운트 관리",
}

var mountAddCmd = &cobra.Command{
	Use:   "add <provider> <token>",
	Short: "클라우드 스토리지 마운트 추가",
	Long: `클라우드 스토리지를 서버에 마운트합니다.

지원 provider:
  onedrive  Microsoft OneDrive

토큰 발급 방법 (로컬 PC에서):
  rclone authorize "onedrive"
  → 브라우저 로그인 후 터미널에 출력된 JSON을 복사`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		provider := args[0]
		token := args[1]

		switch provider {
		case "onedrive":
			return mountOnedrive(token, mountOnedrivePath)
		default:
			return fmt.Errorf("지원하지 않는 provider: %s\n지원 목록: onedrive", provider)
		}
	},
}

func mountOnedrive(token, mountPath string) error {
	fmt.Println("📂 OneDrive 마운트 중...")

	if err := runEmbeddedScriptWithEnv("mount-onedrive.sh",
		[]string{token, mountPath},
		nil,
	); err != nil {
		return fmt.Errorf("OneDrive 마운트 실패: %v", err)
	}

	fmt.Printf("✅ OneDrive 마운트 완료! (경로: %s)\n", mountPath)
	return nil
}

func init() {
	mountAddCmd.Flags().StringVarP(&mountOnedrivePath, "path", "p", "/root/onedrive", "마운트 경로")
	mountCmd.AddCommand(mountAddCmd)
	rootCmd.AddCommand(mountCmd)
}
