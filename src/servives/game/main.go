package main

import (
	"core/consts/service"
	. "core/libs"
	"core/messages"
	"core/protos/gameProto"
	"core/service"
	_ "net/http/pprof"
	"servives/game/module"
)

func main() {
	//初始化Service
	newService := service.NewService(Service.Game)
	newService.StartIpcServer()
	newService.StartRpcServer(&module.GameRpcServer{})
	newService.StartRpcClient([]string{Service.Platform, Service.Log})
	newService.StartRedis()
	newService.StartMysql()
	newService.StartDebug()

	//模块初始化
	initModule()

	//保持进程
	Run()
}

func initModule() {
	messages.RegisterIpcServerHandle(gameProto.ID_user_getInfo_c2s, module.GetInfo)
}
