# 启动后端
Start-Process -NoNewWindow powershell -ArgumentList "cd backend; go run main.go"

# 启动前端
Start-Process -NoNewWindow powershell -ArgumentList "cd frontend/go-web-frontend; npm run dev"

# 保持脚本运行
Write-Host "服务已启动。按 Ctrl+C 停止所有服务..."
while ($true) { Start-Sleep -Seconds 1 } 