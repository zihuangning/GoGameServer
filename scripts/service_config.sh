#!/bin/bash

# 解析run.sh中的服务配置
parse_services() {
    # 单节点服务(不带节点号的日志文件)
    declare -g -A SINGLE_NODE_SERVICES=(
        ["api"]="1"
        ["log"]="1"
        ["login"]="1"
        ["chat"]="1"
        ["test"]="1"
    )

    # 按照run.sh中的启动顺序定义服务（反序，因为停止时需要反向顺序）
    declare -g -a SERVICES=(
        "test:1"
        "api:1"
        "connector:1,2"
        "chat:1"
        "login:1"
        "game:1,2"
        "log:1"
    )
}

# 获取服务的所有节点
get_service_nodes() {
    local NAME=$1
    local NODES=""
    
    for SERVICE in "${SERVICES[@]}"; do
        if [ "${SERVICE%%:*}" = "$NAME" ]; then
            NODES=${SERVICE#*:}
            break
        fi
    done
    
    echo "$NODES"
}

# 检查服务是否存在
is_valid_service() {
    local NAME=$1
    local FOUND=0
    
    for SERVICE in "${SERVICES[@]}"; do
        if [ "${SERVICE%%:*}" = "$NAME" ]; then
            FOUND=1
            break
        fi
    done
    
    echo $FOUND
} 