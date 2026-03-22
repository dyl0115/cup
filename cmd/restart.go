package cmd

import (
	"fmt"

	"github.com/dyl0115/cup/internal/executor"
	"github.com/dyl0115/cup/internal/registry"
	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart <id>",
	Short: "서비스 재시작",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		cf, err := registry.Load(id)
		if err != nil {
			return err
		}

		if len(cf.Restart) > 0 {
			fmt.Printf("🔄 [%s] 재시작 중...\n", id)
			if err := executor.Run(cf.Restart); err != nil {
				return err
			}
		} else {
			fmt.Printf("⏹️  [%s] 중지 중...\n", id)
			if err := executor.Run(cf.Stop); err != nil {
				return err
			}
			fmt.Printf("▶️  [%s] 시작 중...\n", id)
			if err := executor.Run(cf.Start); err != nil {
				return err
			}
		}

		fmt.Printf("✅ [%s] 재시작 완료!\n", id)
		return nil
	},
}
