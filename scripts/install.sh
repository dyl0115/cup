#!/bin/bash

set -e

if [ "$EUID" -ne 0 ]; then
  echo "root 권한이 필요해요. sudo로 실행해주세요!"
  exit 1
fi

DOMAIN=${1:-""}
if [ -z "$DOMAIN" ]; then
  echo "사용법: cup install <도메인>"
  exit 1
fi

# NGINX_TEMP_CONF, NGINX_CONF 는 Go가 환경변수로 주입

echo "=== nginx 설치 ==="
apt update
apt install -y nginx

echo "=== 방화벽 포트 오픈 ==="
iptables -I INPUT -p tcp --dport 80 -j ACCEPT
iptables -I INPUT -p tcp --dport 443 -j ACCEPT

echo "=== 임시 nginx.conf 적용 (80포트만) ==="
cp "$NGINX_TEMP_CONF" /etc/nginx/nginx.conf
sed -i "s/\$DOMAIN/$DOMAIN/g" /etc/nginx/nginx.conf

echo "=== nginx 시작 ==="
systemctl start nginx
systemctl enable nginx

echo "=== certbot 설치 ==="
apt install -y certbot python3-certbot-nginx

echo "=== Let's Encrypt 인증서 발급 ==="
certbot --nginx -d $DOMAIN

echo "=== 최종 nginx.conf 적용 ==="
cp "$NGINX_CONF" /etc/nginx/nginx.conf
sed -i "s/\$DOMAIN/$DOMAIN/g" /etc/nginx/nginx.conf
systemctl reload nginx

echo "=== 완료! ==="
