package main

import (
	"fmt"
	"os"

	"github.com/dyl0115/cup/cmd"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	var err error
	switch command {
	case "install":
		err = cmd.Install(args)
	case "add":
		err = cmd.Add(args)
	case "start":
		err = cmd.Start(args)
	case "stop":
		err = cmd.Stop(args)
	case "restart":
		err = cmd.Restart(args)
	case "remove":
		err = cmd.Remove(args)
	case "status":
		err = cmd.Status(args)
	default:
		fmt.Printf("알 수 없는 커맨드: %s\n", command)
		printHelp()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "오류: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println(`cup - OCI 서버 관리 CLI

사용법:
  cup install <domain>   nginx + HTTPS 세팅
  cup add <경로>          cupFile.yaml 읽고 서비스 등록
  cup start <id>         서비스 시작
  cup stop <id>          서비스 중지
  cup restart <id>       서비스 재시작
  cup remove <id>        서비스 제거
  cup status             전체 서비스 상태 확인`)
}
