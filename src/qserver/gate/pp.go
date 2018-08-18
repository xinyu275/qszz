package mgate

import (
	"github.com/golang/protobuf/proto"
	"qserver/mproto"
	"github.com/liangdas/mqant/log"
	"errors"
)

//协议处理模块
func (this *CustomAgent) cPlayerLogin(protoName string, body []byte) (playerId int32, userId string, err error) {
	msg := &mproto.CPlayerLogin{}
	err = proto.Unmarshal(body, msg)
	if err != nil {
		return
	}

	//去登录模块处理
	moduleType, err := mproto.ProtoToModule(protoName)
	if err != nil {
		return
	}
	m := make(map[string]interface{})
	m["token"] = msg.GetToken()
	m["pf"] = msg.GetPf()
	userId = msg.GetToken()

	reply, serr := this.module.RpcInvoke(moduleType, protoName, m)
	if serr != "" {
		return 0, "", errors.New("login call error " + serr)
	}
	playerId = reply.(int32)
	log.Info("login reply %d", playerId)
	return

}