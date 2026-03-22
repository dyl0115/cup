package cmd

import (
	"fmt"
	"strings"

	"github.com/dyl0115/cup/internal/executor"
	"github.com/dyl0115/cup/internal/registry"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "전체 서비스 상태 확인",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ids, err := registry.ListAll()
		if err != nil {
			return err
		}

		if len(ids) == 0 {
			fmt.Println("등록된 서비스가 없습니다.")
			return nil
		}

		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Printf("  %-20s %s\n", "SERVICE", "STATUS")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		for _, id := range ids {
			cf, err := registry.Load(id)
			if err != nil {
				fmt.Printf("  %-20s ❓ (cupFile 읽기 실패)\n", id)
				continue
			}

			if len(cf.Ping) == 0 {
				fmt.Printf("  %-20s ❓ (PING 없음)\n", id)
				continue
			}

			allOk := true
			for _, pingCmd := range cf.Ping {
				out, err := executor.RunCapture(pingCmd)
				if err != nil || strings.TrimSpace(out) == "inactive" || strings.TrimSpace(out) == "failed" {
					allOk = false
					break
				}
			}

			if allOk {
				fmt.Printf("  %-20s ✅ running\n", id)
			} else {
				fmt.Printf("  %-20s ❌ stopped\n", id)
			}
		}

		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		return nil
	},
}
