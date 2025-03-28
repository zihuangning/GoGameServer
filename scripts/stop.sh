#!/bin/bash

echo "开始停止服务..."

# 按照启动的相反顺序停止服务
SERVICES=(
    "api:1"
    "connector:1,2"
    "chat:1"
    "login:1"
    "game:1,2"
    "log:1"
)

# 首先尝试正常关闭
for SERVICE in "${SERVICES[@]}"
do
    NAME=${SERVICE%%:*}
    INSTANCES=${SERVICE#*:}
    
    IFS=',' read -ra SERVERS <<< "$INSTANCES"
    for SERVER in "${SERVERS[@]}"
    do
        PID=$(pgrep -f "./bin/$NAME -e local -s $SERVER")
        if [ ! -z "$PID" ]; then
            echo "正在停止 $NAME-$SERVER 服务 (PID: $PID)"
            kill $PID
        fi
    done
done

# 等待5秒
sleep 5

# 检查是否还有服务在运行，如果有则强制关闭
for SERVICE in "${SERVICES[@]}"
do
    NAME=${SERVICE%%:*}
    INSTANCES=${SERVICE#*:}
    
    IFS=',' read -ra SERVERS <<< "$INSTANCES"
    for SERVER in "${SERVERS[@]}"
    do
        PID=$(pgrep -f "./bin/$NAME -e local -s $SERVER")
        if [ ! -z "$PID" ]; then
            echo "强制停止 $NAME-$SERVER 服务 (PID: $PID)"
            kill -9 $PID
        fi
    done
done

echo "服务停止完成" 