#!/bin/bash

# 项目根目录
export PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# 数据库配置
export DB_USER="root"
export DB_PASSWORD="123456"
export DB_CHARSET="utf8mb4"
export DB_COLLATE="utf8mb4_unicode_ci"

# 目录配置
export DBS_DIR="${PROJECT_ROOT}/dbs"
export BIN_DIR="${PROJECT_ROOT}/bin"
export LOGS_DIR="${BIN_DIR}/logs"
export SERVICES_DIR="${PROJECT_ROOT}/servives"

# 服务列表
export SERVICES=(
    "api"
    "log"
    "game"
    "login"
    "chat"
    "connector"
    "test"
) 