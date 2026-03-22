package scripts

import "embed"

// FS에 scripts/ 폴더 전체를 바이너리에 내장
//go:embed *.sh *.conf
var FS embed.FS
