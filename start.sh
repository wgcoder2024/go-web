#!/bin/bash

# 启动后端
cd backend && go run main.go &
BACKEND_PID=$!

# 启动前端
cd frontend/go-web-frontend && npm run dev &
FRONTEND_PID=$!

echo "服务已启动。按 Ctrl+C 停止所有服务..."

# 捕获 SIGINT 信号（Ctrl+C）
trap "kill $BACKEND_PID $FRONTEND_PID; exit" INT

# 保持脚本运行
wait 