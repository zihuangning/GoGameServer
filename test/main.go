package main

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/spf13/cast"
	"github.com/yicaoyimuys/GoGameServer/core/consts"
	. "github.com/yicaoyimuys/GoGameServer/core/libs"
	"github.com/yicaoyimuys/GoGameServer/core/libs/array"
	"github.com/yicaoyimuys/GoGameServer/core/libs/hash"
	"github.com/yicaoyimuys/GoGameServer/core/libs/protos"
	"github.com/yicaoyimuys/GoGameServer/core/libs/random"
	"github.com/yicaoyimuys/GoGameServer/core/libs/stack"
	"github.com/yicaoyimuys/GoGameServer/core/libs/timer"
	"github.com/yicaoyimuys/GoGameServer/core/service"
	"github.com/yicaoyimuys/GoGameServer/servives/public/gameProto"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var (
	servers      []string
	connectNum   = 1
	userAccounts = []string{}
)

func main() {
	//初始化Service
	newService := service.NewService(consts.Service_Test)
	newService.StartRpcClient([]string{consts.Service_Log})

	//请求服务器连接地址
	resp, err := http.Get("http://127.0.0.1:18881/GetConnector?type=Socket")
	if err != nil {
		ERR("请求服务器连接地址错误", zap.Error(err))
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ERR("请求服务器连接地址错误", zap.Error(err))
		return
	}
	json.Unmarshal(body, &servers)
	for _, v := range servers {
		DEBUG("可用服务器地址", zap.String("Server", v))
	}

	//开始测试
	startTest()

	//保持进程
	Run()
}

func startTest() {
	//初始化使用账号
	initAccount()

	//开启链接
	for i := 0; i < len(userAccounts); i++ {
		go startConnect(userAccounts[i])
	}
}

func initAccount() {
	for {
		account := "ys" + cast.ToString(random.RandIntn(10000))
		if array.InArray(userAccounts, account) {
			continue
		}
		userAccounts = append(userAccounts, account)
		if len(userAccounts) == connectNum {
			break
		}
	}
}

func startConnect(account string) {
	hashCode := hash.GetHash([]byte(account))
	serverIndex := hashCode % uint32(len(servers))
	server := servers[serverIndex]

	//服务器端ip和端口
	addr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		ERR("解析服务器地址失败", zap.String("Account", account), zap.Error(err))
		return
	}

	//申请连接客户端
	conn, err := net.DialTCP("tcp4", nil, addr)
	if err != nil {
		ERR("连接失败", zap.String("Account", account), zap.Error(err))
		return
	}

	INFO("连接成功", zap.String("Account", account))

	client := &clientSession{
		con:     conn,
		account: account,
	}

	go client.receiveMsg()
	client.ping()
	client.login()
}

type clientSession struct {
	con         *net.TCPConn
	account     string
	token       string
	pingTimerId *timer.TimerEvent
	chatTimerId *timer.TimerEvent
	closeFlag   int32

	recvMutex sync.Mutex
	sendMutex sync.Mutex
}

// 关闭
func (this *clientSession) close() {
	if atomic.CompareAndSwapInt32(&this.closeFlag, 0, 1) {
		timer.Remove(this.pingTimerId)
		timer.Remove(this.chatTimerId)
		this.con.Close()
		DEBUG("连接关闭", zap.String("Account", this.account))
	}
}

// 是否关闭
func (this *clientSession) isClose() bool {
	return atomic.LoadInt32(&this.closeFlag) == 1
}

// 平台登录
func (this *clientSession) login() {
	msg := &gameProto.UserLoginC2S{
		Account: protos.String(this.account),
	}
	this.sendMsg(msg)
}

// 获取用户数据
func (this *clientSession) getInfo() {
	msg := &gameProto.UserGetInfoC2S{
		Token: protos.String(this.token),
	}
	this.sendMsg(msg)
}

// 心跳
func (this *clientSession) ping() {
	this.pingTimerId = timer.DoTimer(2000, func() {
		var msg = &gameProto.ClientPingC2S{}
		this.sendMsg(msg)
	})
}

// 加入聊天
func (this *clientSession) joinChat() {
	var msg = &gameProto.UserJoinChatC2S{}
	msg.Token = protos.String(this.token)
	this.sendMsg(msg)
}

// 聊天
func (this *clientSession) chat() {
	delay := random.RandIntRange(10000, 20000)
	this.chatTimerId = timer.DoTimer(uint32(delay), func() {
		var msg = &gameProto.UserChatC2S{}
		msg.Msg = protos.String("你好啊")
		this.sendMsg(msg)
	})
}

// 接收消息
func (this *clientSession) receiveMsg() {
	defer stack.TryError()

	for {
		if this.isClose() {
			break
		}

		//消息头
		msgHead := make([]byte, 2)
		if _, err := io.ReadFull(this.con, msgHead); err != nil {
			ERR("读取消息头失败", zap.String("Account", this.account), zap.Error(err))
			break
		}

		//消息体长度
		msgLen := binary.BigEndian.Uint16(msgHead)
		if msgLen <= 0 {
			ERR("消息体长度错误", zap.String("Account", this.account))
			break
		}

		//消息体
		msgBody := make([]byte, msgLen)
		if _, err := io.ReadFull(this.con, msgBody); err != nil {
			ERR("读取消息体失败", zap.String("Account", this.account), zap.Error(err))
			break
		}

		//消息解析
		protoMsg := protos.UnmarshalProtoMsg(msgBody)
		if protoMsg == protos.NullProtoMsg {
			ERR("收到错误消息ID", zap.String("Account", this.account), zap.Uint16("MsgId", protos.UnmarshalProtoId(msgBody)))
			break
		}
		DEBUG("收到消息ID", zap.String("Account", this.account), zap.Uint16("MsgId", protoMsg.ID))

		//消息处理
		this.handleMsg(protoMsg.ID, protoMsg.Body)
	}
	this.close()
}

func (this *clientSession) handleMsg(msgId uint16, msgData proto.Message) {
	defer stack.TryError()

	if msgId == gameProto.ID_user_login_s2c {
		//登录成功
		data := msgData.(*gameProto.UserLoginS2C)
		this.token = data.GetToken()
		INFO("登录成功", zap.String("Account", this.account), zap.String("Token", this.token))
		//获取用户数据
		this.getInfo()
	} else if msgId == gameProto.ID_user_getInfo_s2c {
		//获取用户信息成功
		data := msgData.(*gameProto.UserGetInfoS2C)
		INFO("获取用户信息成功", zap.String("Account", this.account), zap.Any("Data", data.GetData()))
		//加入聊天
		this.joinChat()
	} else if msgId == gameProto.ID_user_joinChat_s2c {
		//加入聊天成功
		INFO("加入聊天成功", zap.String("Account", this.account))
		//开始聊天
		this.chat()
	} else if msgId == gameProto.ID_user_chat_notice_s2c {
		//收到聊天消息
		data := msgData.(*gameProto.UserChatNoticeS2C)
		INFO("收到聊天消息", zap.String("Account", this.account), zap.String("From", data.GetUserName()), zap.String("Message", data.GetMsg()))
	} else if msgId == gameProto.ID_error_notice_s2c {
		data := msgData.(*gameProto.ErrorNoticeS2C)
		errCode := data.GetErrorCode()
		if errCode == consts.ErrCode_SystemError {
			//系统服务错误
			this.close()
			//重新连接
			timer.SetTimeOut(3000, func() {
				go startConnect(this.account)
			})
		}
		ERR("收到错误消息", zap.String("Account", this.account), zap.Any("Data", data))
	}
}

func (this *clientSession) sendMsg(msg proto.Message) {
	if this.isClose() {
		return
	}

	defer stack.TryError()

	this.sendMutex.Lock()
	defer this.sendMutex.Unlock()

	msgBytes := protos.MarshalProtoMsg(msg)

	msgLen := uint16(len(msgBytes))
	sendMsg := make([]byte, msgLen+2)
	binary.BigEndian.PutUint16(sendMsg[:2], msgLen)
	copy(sendMsg[2:], msgBytes)

	this.con.Write(sendMsg)
}
