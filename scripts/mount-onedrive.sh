#!/bin/bash

set -e

if [ "$EUID" -ne 0 ]; then
  echo "root 권한이 필요해요. sudo로 실행해주세요!"
  exit 1
fi

ONEDRIVE_TOKEN=${1:-""}
MOUNT_PATH=${2:-"/root/onedrive"}

if [ -z "$ONEDRIVE_TOKEN" ]; then
  echo "사용법: mount-onedrive.sh '<token_json>' [mount_path]"
  exit 1
fi

echo "=== rclone 설치 ==="
if command -v rclone &> /dev/null; then
  echo "rclone 이미 설치되어 있음. 스킵."
else
  apt update
  apt install -y rclone
fi

# 이미 마운트 중이면 재시작만 하고 종료
if systemctl is-active --quiet rclone-onedrive; then
  echo "✅ OneDrive 이미 마운트 중. 재시작합니다..."
  systemctl restart rclone-onedrive
  echo "✅ 재시작 완료! (경로: $MOUNT_PATH)"
  exit 0
fi

echo "=== 임시 rclone.conf 생성 ==="
mkdir -p ~/.config/rclone
cat > ~/.config/rclone/rclone.conf << EOF
[onedrive]
type = onedrive
token = $ONEDRIVE_TOKEN
EOF

echo "=== drive_id 자동 조회 (Microsoft Graph API) ==="
ACCESS_TOKEN=$(echo "$ONEDRIVE_TOKEN" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
DRIVE_ID=$(curl -s -H "Authorization: Bearer $ACCESS_TOKEN" "https://graph.microsoft.com/v1.0/me/drive" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

echo "=== drive_id 자동 조회 ==="
DRIVES_JSON=$(rclone backend drives onedrive: 2>/dev/null)
DRIVE_ID=$(echo "$DRIVES_JSON" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -z "$DRIVE_ID" ]; then
  echo "❌ drive_id 조회 실패. 토큰이 올바른지 확인해주세요."
  exit 1
fi

echo "✅ drive_id 조회 완료: $DRIVE_ID"

echo "=== 최종 rclone.conf 생성 ==="
cat > ~/.config/rclone/rclone.conf << EOF
[onedrive]
type = onedrive
token = $ONEDRIVE_TOKEN
drive_id = $DRIVE_ID
drive_type = personal
EOF

echo "=== 마운트 경로 생성 ==="
mkdir -p "$MOUNT_PATH"

echo "=== systemd 서비스 등록 ==="
cat > /etc/systemd/system/rclone-onedrive.service << EOF
[Unit]
Description=rclone OneDrive mount
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=rclone mount onedrive: $MOUNT_PATH --vfs-cache-mode writes --allow-non-empty
ExecStop=/bin/fusermount -u $MOUNT_PATH
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable --now rclone-onedrive

echo "=== OneDrive 마운트 완료! (경로: $MOUNT_PATH) ==="