package cupfile

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type CupFile struct {
	ID        string   `yaml:"id"`
	Port      int      `yaml:"port"`
	Path      string   `yaml:"path"`
	Init      []string `yaml:"INIT"`
	Start     []string `yaml:"START"`
	Stop      []string `yaml:"STOP"`
	Restart   []string `yaml:"RESTART"`
	Ping      []string `yaml:"PING"`
	Terminate []string `yaml:"TERMINATE"`
}

func Load(dir string) (*CupFile, error) {
	path := dir + "/cupFile.yaml"
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cupFile.yaml을 찾을 수 없습니다: %s", path)
	}

	var cf CupFile
	if err := yaml.Unmarshal(data, &cf); err != nil {
		return nil, fmt.Errorf("cupFile.yaml 파싱 오류: %v", err)
	}

	if err := cf.validate(); err != nil {
		return nil, err
	}

	return &cf, nil
}

func (cf *CupFile) validate() error {
	if cf.ID == "" {
		return fmt.Errorf("cupFile.yaml: id가 없습니다")
	}
	if cf.Port == 0 {
		return fmt.Errorf("cupFile.yaml: port가 없습니다")
	}
	if cf.Path == "" {
		return fmt.Errorf("cupFile.yaml: path가 없습니다")
	}
	return nil
}
