package controllers

import (
	"time"

	"github.com/yicaoyimuys/GoGameServer/core/consts"
	. "github.com/yicaoyimuys/GoGameServer/core/libs"
	"github.com/yicaoyimuys/GoGameServer/core/libs/consul"
	"go.uber.org/zap"

	"github.com/astaxie/beego"
)

type ConnectorController struct {
	beego.Controller
}

func init() {

}

func packageServiceName(serviceType string, serviceName string) string {
	return "<" + serviceType + ">" + serviceName
}

func (this *ConnectorController) Get() {
	// 添加重试机制
	maxRetries := 3
	var consulClient *consul.Client
	var err error

	for i := 0; i < maxRetries; i++ {
		consulClient, err = consul.NewClient()
		if err == nil {
			break
		}
		INFO("Waiting for Consul service...", zap.Int("RetryCount", i+1))
		time.Sleep(time.Second)
	}

	if err != nil {
		ERR("Failed to connect to Consul service", zap.Error(err))
		this.Data["json"] = []string{}
		this.ServeJSON()
		return
	}

	serviceName := ""

	typeStr := this.GetString("type")
	if typeStr == "Socket" {
		serviceName = packageServiceName(consts.ServiceType_Socket, consts.Service_Connector)
	} else if typeStr == "WebSocket" {
		serviceName = packageServiceName(consts.ServiceType_WebSocket, consts.Service_Connector)
	} else {
		ERR("Invalid service type", zap.String("type", typeStr))
		this.Data["json"] = []string{}
		this.ServeJSON()
		return
	}

	services := consulClient.GetServices(serviceName)
	if len(services) == 0 {
		WARN("No connector services found", zap.String("serviceName", serviceName))
	}

	this.Data["json"] = services
	this.ServeJSON()
}
