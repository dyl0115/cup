package cmd

import (
	"fmt"
	"strings"

	"github.com/dyl0115/cup/internal/executor"
	"github.com/dyl0115/cup/internal/registry"
)

func Status(args []string) error {
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

		// PING 커맨드들 모두 실행해서 결과 확인
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
}
