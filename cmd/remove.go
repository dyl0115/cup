package cmd

import (
	"fmt"

	"github.com/dyl0115/cup/internal/executor"
	"github.com/dyl0115/cup/internal/proxy"
	"github.com/dyl0115/cup/internal/registry"
)

func Remove(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("사용법: cup remove <id>")
	}
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

	// registry에서 제거
	if err := registry.Delete(id); err != nil {
		return err
	}

	fmt.Printf("✅ [%s] 제거 완료!\n", id)
	return nil
}
