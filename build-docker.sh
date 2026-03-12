#!/bin/bash

# 构建 Docker 镜像
echo "构建 Docker 镜像..."
docker build -t webscreen:latest .

echo ""
echo "构建完成！"
echo ""
echo "运行方式："
echo "1. 使用 docker-compose: docker-compose up -d"
echo "2. 直接运行: docker run -d --privileged --network host -p 8079:8079 -e PIN=123456 webscreen:latest"
echo ""
echo "访问地址: http://localhost:8079"
echo "PIN码: 123456"
