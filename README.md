# cup

OCI 서버 관리 CLI. nginx 기반 리버스 프록시를 중심으로 각 서비스를 `cupFile.yaml`로 통합 관리한다.

## 설치

```bash
git clone https://github.com/dyl0115/cup-cli.git
cd cup-cli
go mod tidy
go build -o cup .
sudo mv cup /usr/local/bin/
```

## 사용법

### 프록시 서버 초기 설치 (최초 1회)

```bash
sudo cup install <domain> --email <email>
# 예시
sudo cup install dymcp.duckdns.org --email your@email.com
# 단축 플래그
sudo cup install dymcp.duckdns.org -e your@email.com
```

### 서비스 등록

GitHub URL을 바로 넘기면 clone부터 등록까지 자동으로 처리한다.

```bash
# GitHub URL로 바로 등록 (권장)
cup add https://github.com/dyl0115/my-go-server

# 이미 clone된 레포라면 git pull 후 등록
cup add https://github.com/dyl0115/my-go-server

# 로컬 경로로도 가능
cup add /home/ubuntu/my-go-server
cup add .
```

GitHub URL로 실행하면 `/opt/cup-services/<레포명>` 에 자동으로 clone된다.

### 서비스 관리

```bash
cup start <id>      # 서비스 시작
cup stop <id>       # 서비스 중지
cup restart <id>    # 서비스 재시작 (RESTART 섹션 없으면 STOP → START 자동)
cup remove <id>     # 서비스 완전 제거
```

### 전체 상태 확인

```bash
cup status
```

출력 예시:
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  SERVICE              STATUS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  my-go-server         ✅ running
  my-spring-server     ❌ stopped
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## cupFile.yaml

각 서비스 레포 루트에 위치하는 설정 파일. cup 프로젝트 자체에는 없음.

```yaml
id: my-go-server      # 서비스 식별자
port: 8080            # 서비스 포트
path: /api            # nginx 라우팅 경로

INIT:
  - go build -o server .
  - mv server /usr/local/bin/my-go-server

START:
  - systemctl start my-go-server

STOP:
  - systemctl stop my-go-server

RESTART:              # 선택사항. 없으면 cup이 STOP → START 자동 실행
  - systemctl restart my-go-server

PING:
  - systemctl is-active my-go-server

TERMINATE:
  - systemctl disable --now my-go-server
  - rm /usr/local/bin/my-go-server
```

언어마다 INIT이 다름:

```yaml
# Spring
INIT:
  - mvn clean package
  - cp target/app.jar /opt/my-spring-server/

# Python
INIT:
  - pip install -r requirements.txt
```

## GitHub Actions 연동

각 서비스 레포에서 Actions로 자동 빌드/배포 후 cup restart까지 자동화 가능:

```yaml
# .github/workflows/deploy.yml 마지막 step
- name: restart
  run: ssh server "cup restart my-go-server"
```

### 클라우드 스토리지 마운트

`cup mount`로 클라우드 스토리지를 서버에 마운트한다. nginx/서비스 관리와 별개로 동작.

```bash
# 토큰 발급 (로컬 PC에서 먼저 실행)
rclone authorize "onedrive"
# → 브라우저 로그인 후 터미널에 JSON 토큰 출력

# OneDrive 마운트 (기본 경로: /root/onedrive)
sudo cup mount add onedrive '<token_json>'

# 마운트 경로 직접 지정
sudo cup mount add onedrive '<token_json>' --path /root/myonedrive
```

내부적으로 rclone + systemd 서비스로 등록되어 서버 재시작 시 자동 마운트된다.

## registry

`cup add` 시 서비스 경로를 `/etc/cup/registry/<id>` 에 저장.  
이후 cupFile.yaml이 업데이트돼도 항상 원본 경로의 최신 파일을 참조한다.

```
/etc/cup/registry/
├── my-go-server       # 내용: "/opt/cup-services/my-go-server"
└── my-spring-server   # 내용: "/opt/cup-services/my-spring-server"
```
