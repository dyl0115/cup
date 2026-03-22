package executor

import (
	"fmt"
	"os/exec"
)

// Run은 shell 커맨드 목록을 순서대로 실행한다.
// 하나라도 실패하면 즉시 중단하고 에러를 반환한다.
func Run(commands []string) error {
	for _, command := range commands {
		fmt.Printf("  $ %s\n", command)
		cmd := exec.Command("bash", "-c", command)
		cmd.Stdout = nil // 실시간 출력
		cmd.Stderr = nil

		// 실시간으로 출력 보여주기
		cmd.Stdout = newPrefixWriter("  ")
		cmd.Stderr = newPrefixWriter("  ")

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("커맨드 실패 [%s]: %v", command, err)
		}
	}
	return nil
}

// RunCapture는 커맨드를 실행하고 stdout을 문자열로 반환한다. (PING 등에 사용)
func RunCapture(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.Output()
	return string(out), err
}
