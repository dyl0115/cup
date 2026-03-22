#!/bin/bash

set -e

if [ "$EUID" -ne 0 ]; then
  echo "root 권한이 필요해요. sudo로 실행해주세요!"
  exit 1
fi

ID=${1:-""}
PORT=${2:-""}
NGINX_PATH=${3:-""}

if [ -z "$ID" ] || [ -z "$PORT" ] || [ -z "$NGINX_PATH" ]; then
  echo "사용법: add-server.sh [id] [port] [path]"
  exit 1
fi

NGINX_CONF="/etc/nginx/nginx.conf"

LOCATION="    location $NGINX_PATH {\n        proxy_pass http://localhost:$PORT;\n    }"

sed -i "/# === SERVERS START ===/a $LOCATION" $NGINX_CONF

systemctl reload nginx

echo "=== $ID 추가 완료! ($NGINX_PATH → $PORT) ==="
