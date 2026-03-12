FROM ubuntu:22.04

# 设置环境变量
ENV DEBIAN_FRONTEND=noninteractive
ENV PIN=123456
ENV PORT=8079

# 安装必要的依赖
RUN apt-get update && apt-get install -y \
    android-tools-adb \
    curl \
    && rm -rf /var/lib/apt/lists/*

# 创建工作目录
WORKDIR /app

# 复制 webscreen 二进制文件
COPY webscreen /app/webscreen

# 设置权限
RUN chmod +x /app/webscreen

# 暴露端口
EXPOSE 8079

# 启动命令
ENTRYPOINT ["sh", "-c", "./webscreen -port ${PORT} -pin ${PIN}"]
