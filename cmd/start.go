package cmd

import (
	"fmt"

	"github.com/dyl0115/cup/internal/executor"
	"github.com/dyl0115/cup/internal/registry"
)

func Start(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("사용법: cup start <id>")
	}
	id := args[0]

	cf, err := registry.Load(id)
	if err != nil {
		return err
	}

	fmt.Printf("▶️  [%s] 시작 중...\n", id)
	if err := executor.Run(cf.Start); err != nil {
		return err
	}
	fmt.Printf("✅ [%s] 시작 완료!\n", id)
	return nil
}
