package main

import (
	"github.com/yicaoyimuys/GoGameServer/core/consts/Service"
	. "github.com/yicaoyimuys/GoGameServer/core/libs"
	"github.com/yicaoyimuys/GoGameServer/core/messages"
	"github.com/yicaoyimuys/GoGameServer/core/service"
	"github.com/yicaoyimuys/GoGameServer/servives/chat/module"
	"github.com/yicaoyimuys/GoGameServer/servives/public/gameProto"
	"github.com/yicaoyimuys/GoGameServer/servives/public/rpcModules"
)

func main() {
	//初始化Service
	newService := service.NewService(Service.Chat)
	newService.StartIpcServer()
	newService.StartRpcServer()
	newService.StartRpcClient([]string{Service.Log})
	newService.StartRedis()
	newService.RegisterRpcModule("Client", &rpcModules.Client{})

	//消息初始化
	initMessage()

	//模块初始化
	initModule()

	//保持进程
	Run()
}

func initMessage() {
	messages.RegisterIpcServerHandle(gameProto.ID_user_joinChat_c2s, module.JoinChat)
	messages.RegisterIpcServerHandle(gameProto.ID_user_chat_c2s, module.Chat)
}

func initModule() {

}
