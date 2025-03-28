#!/bin/sh

# 停止所有正在运行的服务
sh stop.sh

# 初始化 MySQL 数据库
echo "Initializing MySQL databases..."
mysql -u root -p123456 << EOF
CREATE DATABASE IF NOT EXISTS test_global_1 CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS test_user_1 CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS test_log_1 CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EOF

# 编译所有服务
echo "Building services..."
go build -o ./bin/api ./servives/api/main.go
go build -o ./bin/log ./servives/log/main.go
go build -o ./bin/game ./servives/game/main.go
go build -o ./bin/login ./servives/login/main.go
go build -o ./bin/chat ./servives/chat/main.go
go build -o ./bin/connector ./servives/connector/main.go


# 启动服务并记录日志到 /tmp/log
echo "Starting services..."
LOG_DIR="./bin/logs"

# 删除旧的日志文件，如果它存在的话
if [ -f "$LOG_FILE" ]; then
    rm "$LOG_FILE"
fi
# 检查并创建日志目录
if [ ! -d "$LOG_DIR" ]; then
    mkdir -p "$LOG_DIR"
fi

# 按依赖顺序启动服务
echo "Starting services..."

# 1. 先启动基础服务
./bin/log -e local -s 1 > "$LOG_DIR/log.log" 2>&1 &
echo "log started"
sleep 2

# 2. 启动业务服务
./bin/game -e local -s 1 > "$LOG_DIR/game_1.log" 2>&1 &
echo "game_1 started"
sleep 2

./bin/game -e local -s 2 > "$LOG_DIR/game_2.log" 2>&1 &
echo "game_2 started"
sleep 2

./bin/login -e local -s 1 > "$LOG_DIR/login.log" 2>&1 &
echo "login started"
sleep 2

./bin/chat -e local -s 1 > "$LOG_DIR/chat.log" 2>&1 &
echo "chat started"
sleep 2

./bin/connector -e local -s 1 > "$LOG_DIR/connector_1.log" 2>&1 &
echo "connector_1 started"
sleep 2

./bin/connector -e local -s 2 > "$LOG_DIR/connector_2.log" 2>&1 &
echo "connector_2 started"
sleep 2

# 3. 最后启动 API 服务
./bin/api -e local -s 1 > "$LOG_DIR/api.log" 2>&1 &
echo "api started"

echo "All services started."