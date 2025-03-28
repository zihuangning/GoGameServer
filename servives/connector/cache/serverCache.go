package cache

import (
	"encoding/json"

	"github.com/yicaoyimuys/GoGameServer/core/libs/logger"
	"github.com/yicaoyimuys/GoGameServer/core/libs/stack"
	"github.com/yicaoyimuys/GoGameServer/servives/public/redisInstances"
	"github.com/yicaoyimuys/GoGameServer/servives/public/redisKeys"
	"go.uber.org/zap"
)

func SetServerInfo(domainName string, serverPort string, onlineUsersNum int) {
	defer stack.TryError()

	if domainName == "" || serverPort == "" {
		logger.Error("Invalid domainName or serverPort", zap.String("domainName", domainName), zap.String("serverPort", serverPort))
		return
	}

	redisClient := redisInstances.Global()
	if redisClient == nil {
		logger.Error("Redis client is nil")
		return
	}

	oldServerInfo := GetServerInfo(domainName, serverPort)

	//读取最高在线
	var oldMaxOnlineUsersNum = 0
	if oldServerInfo != nil {
		if num, exists := oldServerInfo["maxOnlineUsersNum"]; exists {
			oldMaxOnlineUsersNum = num
		}
	}

	redisKey := redisKeys.ServerInfo
	serverKey := domainName + ":" + serverPort

	serverInfo := make(map[string]int)
	serverInfo["onlineUsersNum"] = onlineUsersNum
	if onlineUsersNum > oldMaxOnlineUsersNum {
		serverInfo["maxOnlineUsersNum"] = onlineUsersNum
	} else {
		serverInfo["maxOnlineUsersNum"] = oldMaxOnlineUsersNum
	}

	byteData, err := json.Marshal(serverInfo)
	if err != nil {
		logger.Error("Failed to marshal server info", zap.Error(err))
		return
	}

	err = redisClient.HSet(redisKey, serverKey, string(byteData)).Err()
	if err != nil {
		logger.Error("Failed to set server info to redis", zap.Error(err))
	}
}

func GetServerInfo(domainName string, serverPort string) map[string]int {
	defer stack.TryError()

	if domainName == "" || serverPort == "" {
		logger.Error("Invalid domainName or serverPort", zap.String("domainName", domainName), zap.String("serverPort", serverPort))
		return nil
	}

	redisClient := redisInstances.Global()
	if redisClient == nil {
		logger.Error("Redis client is nil")
		return nil
	}

	redisKey := redisKeys.ServerInfo
	serverKey := domainName + ":" + serverPort
	val, err := redisClient.HGet(redisKey, serverKey).Result()
	if err != nil {
		return nil
	}

	var serverInfo map[string]int
	err = json.Unmarshal([]byte(val), &serverInfo)
	if err != nil {
		logger.Error("Failed to unmarshal server info", zap.Error(err))
		return nil
	}
	return serverInfo
}
