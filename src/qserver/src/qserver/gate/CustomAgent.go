package mgate

import (
	"bufio"
	"github.com/golang/protobuf/proto"
	"github.com/juju/errors"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/gate/base"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/network"
	"github.com/liangdas/mqant/rpc/util"
	"github.com/liangdas/mqant/utils"
	"qserver/mproto"
	"strconv"
	"time"
)

func NewAgent(module module.RPCModule) *CustomAgent {
	a := &CustomAgent{
		module: module,
	}
	return a
}

type CustomAgent struct {
	gate.Agent
	module                           module.RPCModule
	session                          gate.Session
	conn                             network.Conn
	r                                *bufio.Reader
	w                                *bufio.Writer
	gate                             gate.Gate
	rev_num                          int64
	send_num                         int64
	last_storage_heartbeat_data_time int64 //上一次发送存储心跳时间
	isclose                          bool
}

func (this *CustomAgent) OnInit(gate gate.Gate, conn network.Conn) error {
	log.Info("CustomAgent", "OnInit")
	this.conn = conn
	this.gate = gate
	this.r = bufio.NewReader(conn)
	this.w = bufio.NewWriter(conn)
	this.isclose = false
	this.rev_num = 0
	this.send_num = 0
	return nil
}

/**
给客户端发送消息
*/
func (this *CustomAgent) WriteMsg(protoName string, body []byte) error {
	this.send_num++
	//粘包完成后调下面的语句发送数据
	b := mproto.Pack(protoName, body)
	this.w.Write(b)
	this.w.Flush()
	//this.w.Write()
	return nil
}

//给客户端发消息
func (this *CustomAgent) Write(msg []byte) error {
	this.send_num++
	this.w.Write(msg)
	this.w.Flush()
	return nil
}

func (this *CustomAgent) Run() (err error) {
	log.Info("CustomAgent", "开始读数据了")
	protoName, body, err := ReadPack(this.r)
	if err != nil {
		log.Error("gate read data error：%s", err.Error())
		return
	}
	//处理登录包
	if protoName != "c_player_login" {
		err = errors.New("first pack is not login package")
		return
	}
	log.Info("sucess read data:%s, %v", protoName, body)

	//到登陆系统登陆，如果玩家不存在就自动注册
	playerId, userId, err := this.cPlayerLogin(protoName, body)
	if err != nil {
		log.Error("login error", err.Error())
		return
	}
	this.session, err = this.gate.NewSessionByMap(map[string]interface{}{
		"Sessionid": utils.GenerateID().String(),
		"Userid":    userId,
		"Network":   this.conn.RemoteAddr().Network(),
		"IP":        this.conn.RemoteAddr().String(),
		"Serverid":  this.module.GetServerId(),
		"Settings": map[string]string{
			"PlayerId": strconv.Itoa(int(playerId)),
		},
	})
	if err != nil {
		log.Error("gate create agent fail", err.Error())
		return
	}
	this.gate.GetAgentLearner().Connect(this)

	//返回客户端登陆成功
	this.login_success_reply(mproto.ECode_ok)

	//这里可以循环读取客户端的数据
	for {
		protoName, body, err := ReadPack(this.r)
		if err != nil {
			log.Error("recv data error", err.Error())
			break
		}
		this.OnRecover(protoName, body)
	}
	//这个函数返回后连接就会被关闭
	return nil
}

/**
接收到一个数据包
*/
func (this *CustomAgent) OnRecover(protoName string, body []byte) {
	moduleType, err := mproto.ProtoToModule(protoName)
	//如果moduleType == Battle,那要根据session的房间调用到哪个battle上，例如：Battle001@Battle
	if err != nil {
		log.Error("ProtoToModule error:protoName %s, err %s", protoName, err.Error())
		return
	}
	_func := protoName

	//如果要对这个请求进行分布式跟踪调试,就执行下面这行语句
	//a.session.CreateRootSpan("gate")
	var ArgsType []string = make([]string, 2)
	var args [][]byte = make([][]byte, 2)
	//封装session
	ArgsType[0] = basegate.RPC_PARAM_SESSION_TYPE
	b, err := this.GetSession().Serializable()
	if err != nil {
		log.Error(err.Error())
		return
	}
	args[0] = b

	if err != nil {
		log.Error("mproto Deode error ", err.Error())
		return
	}
	ArgsType[1] = argsutil.BYTES
	args[1] = body
	result, e := this.module.RpcInvokeArgs(moduleType, _func, ArgsType, args)
	if e != "" {
		log.Error("RpcInvokeArgs error, ", e)
		return
	}
	reply := result.([]byte)
	if len(reply) > 0 {
		this.Write(reply)
	}

	this.heartbeat()
}

func (this *CustomAgent) heartbeat() {
	//自定义网关需要你自己设计心跳协议
	if this.GetSession().GetUserId() != "" {
		//这个链接已经绑定Userid
		interval := time.Now().UnixNano()/1000000/1000 - this.last_storage_heartbeat_data_time //单位秒
		if interval > this.gate.GetMinStorageHeartbeat() {
			//如果用户信息存储心跳包的时长已经大于一秒
			if this.gate.GetStorageHandler() != nil {
				this.gate.GetStorageHandler().Heartbeat(this.GetSession().GetUserId())
				this.last_storage_heartbeat_data_time = time.Now().UnixNano() / 1000000 / 1000
			}
		}
	}
}

func (this *CustomAgent) Close() {
	log.Info("CustomAgent", "主动断开连接")
	this.conn.Close()
}
func (this *CustomAgent) OnClose() error {
	this.isclose = true
	log.Info("CustomAgent", "连接断开事件")
	//这个一定要调用，不然gate可能注销不了,造成内存溢出
	this.gate.GetAgentLearner().DisConnect(this) //发送连接断开的事件
	return nil
}
func (this *CustomAgent) Destroy() {
	this.conn.Destroy()
}
func (this *CustomAgent) RevNum() int64 {
	return this.rev_num
}
func (this *CustomAgent) SendNum() int64 {
	return this.send_num
}
func (this *CustomAgent) IsClosed() bool {
	return this.isclose
}
func (this *CustomAgent) GetSession() gate.Session {
	return this.session
}

//登陆成功返回信息
func (this *CustomAgent) login_success_reply(code mproto.ECode) {
	SPlayerLogin := &mproto.SPlayerLogin{
		Code: &code,
	}
	body, _ := proto.Marshal(SPlayerLogin)
	this.WriteMsg("s_player_login", body)
}
