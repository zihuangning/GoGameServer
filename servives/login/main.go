package main

import (
	"github.com/yicaoyimuys/GoGameServer/core/consts"
	. "github.com/yicaoyimuys/GoGameServer/core/libs"
	"github.com/yicaoyimuys/GoGameServer/core/messages"
	"github.com/yicaoyimuys/GoGameServer/core/service"
	"github.com/yicaoyimuys/GoGameServer/servives/login/module"
	"github.com/yicaoyimuys/GoGameServer/servives/public/gameProto"
)

func main() {
	//初始化Service
	newService := service.NewService(consts.Service_Login)
	newService.StartIpcServer()
	newService.StartRpcServer()
	newService.StartRpcClient([]string{consts.Service_Log})
	newService.StartRedis()
	newService.StartMysql()

	//消息初始化
	initMessage()

	//模块初始化
	initModule()

	//保持进程
	Run()
}

func initMessage() {
	messages.RegisterIpcServerHandle(gameProto.ID_user_login_c2s, module.Login)
}

func initModule() {
	//sessions.StartSessionCleanup()
}
