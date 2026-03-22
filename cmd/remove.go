package cmd

import (
	"fmt"

	"github.com/dyl0115/cup/internal/executor"
	"github.com/dyl0115/cup/internal/proxy"
	"github.com/dyl0115/cup/internal/registry"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <id>",
	Short: "서비스 제거",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		cf, err := registry.Load(id)
		if err != nil {
			return err
		}

		fmt.Printf("💥 [%s] TERMINATE 실행 중...\n", id)
		if err := executor.Run(cf.Terminate); err != nil {
			return err
		}

		fmt.Printf("🌐 nginx에서 제거 중... (path: %s)\n", cf.Path)
		if err := proxy.Remove(cf.Path); err != nil {
			return err
		}

		if err := registry.Delete(id); err != nil {
			return err
		}

		fmt.Printf("✅ [%s] 제거 완료!\n", id)
		return nil
	},
}
