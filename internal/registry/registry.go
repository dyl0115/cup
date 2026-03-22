package registry

import (
	"fmt"
	"os"
	"strings"

	"github.com/dyl0115/cup/internal/cupfile"
)

const registryDir = "/etc/cup/registry"

// Save는 서비스 id와 cupFile 경로를 registry에 저장한다.
func Save(id, dir string) error {
	if err := os.MkdirAll(registryDir, 0755); err != nil {
		return fmt.Errorf("registry 디렉토리 생성 실패: %v", err)
	}
	path := registryDir + "/" + id
	return os.WriteFile(path, []byte(dir), 0644)
}

// Load는 registry에서 id로 cupFile을 찾아 로드한다.
func Load(id string) (*cupfile.CupFile, error) {
	path := registryDir + "/" + id
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("[%s] 등록되지 않은 서비스입니다. 먼저 cup add를 실행하세요.", id)
	}
	dir := strings.TrimSpace(string(data))
	return cupfile.Load(dir)
}

// Delete는 registry에서 id를 제거한다.
func Delete(id string) error {
	path := registryDir + "/" + id
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("[%s] registry 제거 실패: %v", id, err)
	}
	return nil
}

// ListAll은 등록된 모든 서비스 id 목록을 반환한다.
func ListAll() ([]string, error) {
	entries, err := os.ReadDir(registryDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("registry 읽기 실패: %v", err)
	}
	var ids []string
	for _, e := range entries {
		if !e.IsDir() {
			ids = append(ids, e.Name())
		}
	}
	return ids, nil
}
