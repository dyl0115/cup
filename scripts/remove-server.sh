#!/bin/bash

set -e

if [ "$EUID" -ne 0 ]; then
  echo "root 권한이 필요해요. sudo로 실행해주세요!"
  exit 1
fi

NGINX_PATH=${1:-""}

if [ -z "$NGINX_PATH" ]; then
  echo "사용법: remove-server.sh [path]"
  exit 1
fi

NGINX_CONF="/etc/nginx/nginx.conf"

sed -i "/location $NGINX_PATH/,/}/d" $NGINX_CONF

systemctl reload nginx

echo "=== $NGINX_PATH 제거 완료! ==="
