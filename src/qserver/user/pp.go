package user

//协议处理模块

import (
	"github.com/golang/protobuf/proto"
	"github.com/liangdas/mqant/gate"
	"qserver/mproto"
	"strconv"
)

//CPlayerInfo ： 前端获取玩家信息协议
func (self *User) cPlayerInfo(session gate.Session, msg []byte) ([]byte, string) {
	playerId, _ := strconv.Atoi(session.Get("PlayerId"))
	player, err := self.getPlayer(playerId)
	if err != nil {
		return nil, err.Error()
	}
	pb := &mproto.SPlayerInfo{
		PlayerId:   proto.Uint32(uint32(player.PlayerId)),
		PlayerName: proto.String(player.Nickname),
		Gender:     mproto.EShinGender(0).Enum(),
		Avatar:     proto.Uint32(0),
		Gold:       proto.Uint32(0),
		Gems:       proto.Uint32(0),
	}
	//更新玩家session
	player.updateSession(session)
	return mproto.PackPBReply("s_player_info", pb)
}
