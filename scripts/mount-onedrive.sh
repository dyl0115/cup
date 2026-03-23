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

echo "=== rclone.conf 생성 ==="
mkdir -p ~/.config/rclone
cat > ~/.config/rclone/rclone.conf << EOF
[onedrive]
type = onedrive
token = $ONEDRIVE_TOKEN
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
