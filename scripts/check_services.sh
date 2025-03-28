#!/bin/bash

# 定义服务及其实例，以及是否需要节点编号的日志文件
declare -A SINGLE_NODE_SERVICES=( 
    ["api"]="1"
    ["log"]="1"
    ["login"]="1"
    ["chat"]="1"
)

# 定义服务及其实例
SERVICES=(
    "api:1"
    "connector:1,2"
    "chat:1"
    "login:1"
    "game:1,2"
    "log:1"
)

# 显示使用帮助
show_usage() {
    echo "使用方法:"
    echo "  $0                    # 检查所有服务"
    echo "  $0 game              # 检查指定服务的所有节点"
    echo "  $0 connector-1       # 检查指定服务的特定节点"
    echo "支持的服务:"
    for SERVICE in "${SERVICES[@]}"; do
        NAME=${SERVICE%%:*}
        INSTANCES=${SERVICE#*:}
        echo "  $NAME (节点: $INSTANCES)"
    done
}

# 安全地执行命令
safe_execute() {
    local cmd="$1"
    local output
    output=$(eval "$cmd" 2>/dev/null) || return 1
    echo "$output"
}

# 检查单个服务节点
check_service_node() {
    local NAME=$1
    local SERVER=$2
    
    echo -e "\n检查 $NAME-$SERVER:"
    PID=$(pgrep -f "./bin/$NAME -e local -s $SERVER")
    
    if [ ! -z "$PID" ]; then
        echo "状态: 运行中"
        echo "进程ID: $PID"
        echo "启动时间: $(ps -p $PID -o lstart=)"
        echo "内存使用: $(ps -p $PID -o %mem=)%"
        echo "CPU使用: $(ps -p $PID -o %cpu=)%"
        echo "完整命令: $(ps -p $PID -o cmd=)"
        
        # 检查日志文件
        if [ "${SINGLE_NODE_SERVICES[$NAME]}" = "1" ]; then
            LOG_FILE="./bin/logs/${NAME}.log"
        else
            LOG_FILE="./bin/logs/${NAME}_${SERVER}.log"
        fi
        
        if [ -f "$LOG_FILE" ]; then
            echo "最新日志:"
            tail -n 3 "$LOG_FILE" 2>/dev/null || echo "无法读取日志文件"
        else
            echo "警告: 日志文件不存在 ($LOG_FILE)"
        fi
    else
        echo "状态: 未运行"
        echo "警告: 服务未启动"
    fi
    echo "--------------------------------"
}

# 显示系统资源使用情况
show_system_status() {
    echo -e "\n========== 系统资源使用 =========="
    echo "CPU使用率:"
    top -bn1 | grep "Cpu(s)" | awk '{print $2 + $4}' | awk '{print $1"%"}'
    echo "内存使用率:"
    free -m | awk 'NR==2{printf "%.2f%%\n", $3*100/$2}'
    echo "================================="
}

# 检查单个服务的所有节点
check_service() {
    local NAME=$1
    local FOUND=0
    
    for SERVICE in "${SERVICES[@]}"; do
        if [ "${SERVICE%%:*}" = "$NAME" ]; then
            FOUND=1
            INSTANCES=${SERVICE#*:}
            break
        fi
    done
    
    if [ $FOUND -eq 0 ]; then
        echo "错误: 未知服务 '$NAME'"
        show_usage
        exit 1
    fi
    
    echo "========== 服务状态检查 =========="
    echo "检查时间: $(date '+%Y-%m-%d %H:%M:%S')"
    echo "================================="
    
    IFS=',' read -ra SERVERS <<< "$INSTANCES"
    for SERVER in "${SERVERS[@]}"; do
        check_service_node "$NAME" "$SERVER"
    done
    
    show_system_status
}

# 检查所有服务
check_all_services() {
    echo "========== 服务状态检查 =========="
    echo "检查时间: $(date '+%Y-%m-%d %H:%M:%S')"
    echo "================================="
    
    for SERVICE in "${SERVICES[@]}"; do
        NAME=${SERVICE%%:*}
        INSTANCES=${SERVICE#*:}
        IFS=',' read -ra SERVERS <<< "$INSTANCES"
        for SERVER in "${SERVERS[@]}"; do
            check_service_node "$NAME" "$SERVER"
        done
    done
    
    show_system_status
}

# 主逻辑
if [ $# -eq 0 ]; then
    check_all_services
elif [ $# -eq 1 ]; then
    if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
        show_usage
        exit 0
    fi
    
    # 检查是否是特定节点格式 (service-number)
    if [[ $1 == *-* ]]; then
        SERVICE_NAME=${1%-*}
        NODE_NUMBER=${1#*-}
        FOUND=0
        for SERVICE in "${SERVICES[@]}"; do
            if [ "${SERVICE%%:*}" = "$SERVICE_NAME" ]; then
                FOUND=1
                break
            fi
        done
        if [ $FOUND -eq 0 ]; then
            echo "错误: 未知服务 '$SERVICE_NAME'"
            show_usage
            exit 1
        fi
        check_service_node "$SERVICE_NAME" "$NODE_NUMBER"
        show_system_status
    else
        # 检查整个服务
        check_service "$1"
    fi
else
    echo "错误: 参数过多"
    show_usage
    exit 1
fi 