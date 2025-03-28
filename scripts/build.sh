#!/bin/bash

# 加载配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/config.sh"

# 执行SQL文件
execute_sql() {
    local db_name=$1
    local sql_file=$2
    
    echo "Creating database $db_name and executing $sql_file..."
    
    # 创建数据库并执行SQL文件
    mysql -u "${DB_USER}" -p"${DB_PASSWORD}" -e "CREATE DATABASE IF NOT EXISTS $db_name CHARACTER SET ${DB_CHARSET} COLLATE ${DB_COLLATE};" && \
    mysql -u "${DB_USER}" -p"${DB_PASSWORD}" "$db_name" < "$sql_file" || {
        echo "Error: Failed to execute $sql_file"
        return 1
    }
    
    echo "Successfully executed $sql_file"
}

# 初始化目录
init_directories() {
    mkdir -p "${BIN_DIR}" "${LOGS_DIR}"
    
    if [ ! -d "${DBS_DIR}" ]; then
        echo "Error: dbs directory not found: ${DBS_DIR}"
        exit 1
    fi
}

# 初始化数据库
init_databases() {
    echo "Initializing databases..."
    
    for sql_file in "${DBS_DIR}"/*.sql; do
        if [ -f "$sql_file" ]; then
            db_name=$(basename "$sql_file" .sql | cut -d'_' -f1)
            execute_sql "$db_name" "$sql_file"
        fi
    done
}

# 编译服务
build_services() {
    echo "Building services..."
    
    for service in "${SERVICES[@]}"; do
        echo "Building $service service..."
        go build -o "${BIN_DIR}/${service}" "${SERVICES_DIR}/${service}/main.go" || {
            echo "Error: Failed to build $service service"
            exit 1
        }
    done
    
    echo "All services compiled successfully"
}

# 主函数
main() {
    init_directories
    init_databases
    build_services
}

# 执行主函数
main
