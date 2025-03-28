#!/bin/bash

# 加载服务配置
source "$(dirname "$0")/service_config.sh"
parse_services

echo "开始停止服务..."

# 按照启动的相反顺序停止服务
stop_services() {
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
}

# 停止单个服务
stop_service() {
    local NAME=$1
    local FOUND=$(is_valid_service "$NAME")
    
    if [ $FOUND -eq 0 ]; then
        echo "错误: 未知服务 '$NAME'"
        exit 1
    fi
    
    local NODES=$(get_service_nodes "$NAME")
    IFS=',' read -ra SERVERS <<< "$NODES"
    
    # 首先尝试正常关闭
    for SERVER in "${SERVERS[@]}"
    do
        PID=$(pgrep -f "./bin/$NAME -e local -s $SERVER")
        if [ ! -z "$PID" ]; then
            echo "正在停止 $NAME-$SERVER 服务 (PID: $PID)"
            kill $PID
        fi
    done

    # 等待5秒
    sleep 5

    # 检查是否还有服务在运行，如果有则强制关闭
    for SERVER in "${SERVERS[@]}"
    do
        PID=$(pgrep -f "./bin/$NAME -e local -s $SERVER")
        if [ ! -z "$PID" ]; then
            echo "强制停止 $NAME-$SERVER 服务 (PID: $PID)"
            kill -9 $PID
        fi
    done
}

# 主逻辑
if [ $# -eq 0 ]; then
    stop_services
else
    stop_service "$1"
fi

echo "服务停止完成" 