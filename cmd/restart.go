package cmd

import (
	"fmt"

	"github.com/dyl0115/cup/internal/executor"
	"github.com/dyl0115/cup/internal/registry"
)

func Restart(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("사용법: cup restart <id>")
	}
	id := args[0]

	cf, err := registry.Load(id)
	if err != nil {
		return err
	}

	// RESTART 섹션이 있으면 그걸 쓰고, 없으면 STOP → START
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
}
