#!/bin/bash

# 数据库功能测试脚本

echo "🧪 测试 SQLite 数据库功能..."

# 1. 启动应用
echo "1. 启动应用..."
./bin/go-web-template -c config/app.yaml -p 8081 -v &
APP_PID=$!

# 等待应用启动
sleep 3

# 2. 测试管理员登录
echo "2. 测试管理员登录..."
curl -X POST http://localhost:8081/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}' \
  -c cookies.txt

echo -e "\n"

# 3. 测试管理员密码修改
echo "3. 测试管理员密码修改..."
curl -X POST http://localhost:8081/admin/change_password \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{"old_password":"admin","new_password":"newadmin123"}'

echo -e "\n"

# 4. 测试新密码登录
echo "4. 测试新密码登录..."
curl -X POST http://localhost:8081/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"newadmin123"}' \
  -c cookies_new.txt

echo -e "\n"

# 5. 测试 API 登录（使用默认用户）
echo "5. 测试 API 登录..."
curl -X POST http://localhost:8081/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"newadmin123"}'

echo -e "\n"

# 6. 检查数据库文件
echo "6. 检查数据库文件..."
if [ -f "data/app.db" ]; then
    echo "✅ 数据库文件已创建: data/app.db"
    echo "数据库大小: $(du -h data/app.db | cut -f1)"
else
    echo "❌ 数据库文件未找到"
fi

# 7. 停止应用
echo "7. 停止应用..."
kill $APP_PID 2>/dev/null || true

# 清理临时文件
rm -f cookies.txt cookies_new.txt

echo "🎉 数据库功能测试完成！"
