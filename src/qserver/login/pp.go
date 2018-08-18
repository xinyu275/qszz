package login

import (
	"github.com/golang/protobuf/proto"
	"github.com/liangdas/mqant/gate"
	"qserver/common"
	"qserver/mproto"
)

//协议处理模块

//CPlayerLogin ： 玩家登陆 如果成功返回玩家id
func (m *Login) cPlayerLogin(msg map[string]interface{}) (int32, string) {
	accname := msg["token"].(string)
	playerId, err := getPlayerIdByAccname(accname)
	if err != nil {
		return 0, err.Error()
	}
	if playerId == 0 {
		playerId, err = registerUser(accname)
		if err != nil {
			return 0, err.Error()
		}
	}
	return int32(playerId), ""
}

//查询服务器时间
func (m *Login) cServerTime(session gate.Session, body []byte) ([]byte, string) {
	pb := &mproto.SServerTime{
		ServerTime: proto.Uint64(uint64(common.GetMillUnixTime())),
	}
	//更新玩家session
	return mproto.PackPBReply("s_server_time", pb)
}
