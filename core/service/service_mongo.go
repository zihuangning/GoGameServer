package service

import (
	"github.com/yicaoyimuys/GoGameServer/core/config"
	. "github.com/yicaoyimuys/GoGameServer/core/libs"
	"github.com/yicaoyimuys/GoGameServer/core/libs/mongo"
)

func (this *Service) StartMongo() {
	this.mongoClients = make(map[string]*mongo.Client)

	mongoConfigs := config.GetMongoConfig()
	for aliasName, mongoConfig := range mongoConfigs {
		client, err := mongo.NewClient(mongoConfig)
		CheckError(err)

		if client != nil {
			this.mongoClients[aliasName] = client
			INFO("mongo_" + aliasName + "连接成功")
		} else {
			ERR("mongo_" + aliasName + "连接失败")
		}
	}
}

func (this *Service) GetMongoClient(dbAliasName string) *mongo.Client {
	client, _ := this.mongoClients[dbAliasName]
	return client
}
