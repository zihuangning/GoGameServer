#!/bin/sh

# 定义函数
function func() {
    local target=$1
    echo "Searching for processes matching: $target"

    # 使用 pgrep 精确匹配命令行
    local pids=$(pgrep -fl "$target" | awk '{print $1}')
    if [ -n "$pids" ]; then
        echo "Found processes to stop: $pids"
        
        # 尝试发送 SIGTERM 信号
        echo "Sending SIGTERM to processes..."
        kill $pids || true

        # 检查是否仍有进程存活
        local remaining_pids=$(pgrep -fl "$target" | awk '{print $1}')
        if [ -n "$remaining_pids" ]; then
            echo "Some processes are still running. Sending SIGKILL..."
            kill -9 $remaining_pids || true
        else
            echo "All processes terminated successfully."
        fi
    else
        echo "No running processes found for: $target"
    fi
}

# 默认停止所有服务
if [ $# -eq 0 ]; then
    echo "Stopping all services..."
    func "./bin/connector -e local -s 1"
    func "./bin/game -e local -s 1"
    func "./bin/login -e local -s 1"
    func "./bin/chat -e local -s 1"
    func "./bin/log -e local -s 1"
    func "./bin/api -e local -s 1"
else
    # 停止指定的服务
    echo "Stopping service: $1"
    func "$1"
fi