package cmd

import (
	"fmt"

	"github.com/dyl0115/cup/internal/executor"
	"github.com/dyl0115/cup/internal/registry"
)

func Stop(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("사용법: cup stop <id>")
	}
	id := args[0]

	cf, err := registry.Load(id)
	if err != nil {
		return err
	}

	fmt.Printf("⏹️  [%s] 중지 중...\n", id)
	if err := executor.Run(cf.Stop); err != nil {
		return err
	}
	fmt.Printf("✅ [%s] 중지 완료!\n", id)
	return nil
}
