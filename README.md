# cup

OCI 서버 관리 CLI. nginx 기반 리버스 프록시를 중심으로 각 서비스를 `cupFile.yaml`로 통합 관리한다.

## 설치

```bash
git clone https://github.com/dyl0115/cup.git
cd cup
go mod tidy
go build -o cup .
sudo mv cup /usr/local/bin/
```

## 사용법

```bash
# 프록시 서버 초기 설치 (최초 1회)
sudo cup install your.domain.org

# 서비스 등록 (해당 프로젝트 디렉토리에서)
cup add /home/ubuntu/my-go-server

# 서비스 관리
cup start my-go-server
cup stop my-go-server
cup restart my-go-server
cup remove my-go-server

# 전체 상태 확인
cup status
```

## cupFile.yaml 예시

각 서비스 레포 루트에 위치:

```yaml
id: my-go-server
port: 8080
path: /api

INIT:
  - go build -o server .
  - mv server /usr/local/bin/my-go-server

START:
  - systemctl start my-go-server

STOP:
  - systemctl stop my-go-server

RESTART:           # 선택사항. 없으면 cup이 STOP → START 자동 실행
  - systemctl restart my-go-server

PING:
  - systemctl is-active my-go-server

TERMINATE:
  - systemctl disable --now my-go-server
  - rm /usr/local/bin/my-go-server
```

## registry

`cup add` 시 서비스 경로를 `/etc/cup/registry/<id>` 에 저장.  
이후 cupFile.yaml이 업데이트돼도 항상 원본 경로의 최신 파일을 참조한다.
