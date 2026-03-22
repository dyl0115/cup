package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/dyl0115/cup/internal/cupfile"
	"github.com/dyl0115/cup/internal/executor"
	"github.com/dyl0115/cup/internal/proxy"
	"github.com/dyl0115/cup/internal/registry"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <경로>",
	Short: "cupFile.yaml 읽고 서비스 등록",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}

		absDir, err := filepath.Abs(dir)
		if err != nil {
			return fmt.Errorf("경로 변환 실패: %v", err)
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
