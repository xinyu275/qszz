package lobby

//协议处理模块

import (
	"github.com/golang/protobuf/proto"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/log"
	"qserver/mproto"
	"strconv"
)

//cLobbyEnter ： 进入游戏大厅
func (self *Lobbys) cLobbyEnter(session gate.Session, msg []byte) ([]byte, string) {
	lobby, ok := self.GetLobby(0)
	if ok {
		playerId, err := strconv.Atoi(session.Get("PlayerId"))
		if err == nil {
			lobby.AddPlayer(playerId, session)

			// 登录大厅成功
			// TODO 需要房间ID
			code := mproto.ECode_ok
			pb := &mproto.SLobbyEnter {
				Code: 		&code,
				RoomId:   	proto.Int32(0),
			}
			return mproto.PackPBReply("s_lobby_enter", pb)
		} else {
			log.Error("[Error]Have no playerId in cLobbyEnter, session:%s", session.GetSessionId())
		}
	}

	// 登录大厅失败
	code := mproto.ECode_fail
	pb := &mproto.SLobbyEnter {
		Code: 		&code,
		RoomId:   	proto.Int32(-2),
	}
	return mproto.PackPBReply("s_lobby_enter", pb)
}

//cLobbyLeave ： 离开游戏大厅
func (self *Lobbys) cLobbyLeave(session gate.Session, msg []byte) ([]byte, string) {
	// lobby, ok := self.GetLobby(0)
	// if ok {
	// 	playerId, err := strconv.Atoi(session.Get("PlayerId"))
	// 	if err == nil {
	// 		lobby.DelPlayer(playerId)

	// 		// 离开大厅成功
	// 		code := mproto.ECode_ok
	// 		pb := &mproto.SLobbyEnter {
	// 			Code: 		&code,
	// 			RoomId:   	proto.Int32(0),
	// 		}
	// 		return mproto.PackPBReply("s_lobby_enter", pb)
	// 	} else {
	// 		log.Error("[Error]Have no playerId in cLobbyEnter, session:%s", session.GetSessionId())
	// 	}
	// }

	// // 离开大厅成功失败
	// code := mproto.ECode_fail
	// pb := &mproto.SLobbyEnter {
	// 	Code: 		&code,
	// 	RoomId:   	proto.Int32(-2),
	// }
	// return mproto.PackPBReply("s_lobby_enter", pb)
	return nil, ""
}

//cLobbyMatch ： 开始匹配游戏
func (self *Lobbys) cLobbyMatch(session gate.Session, msg []byte) ([]byte, string) {
	lobby, ok := self.GetLobby(0)
	if ok {
		playerId, err := strconv.Atoi(session.Get("PlayerId"))
		if err == nil {
			// 匹配成功，等待排队匹配
			if lobby.Match(playerId) {
				return nil, ""
			}
		} else {
			log.Error("[Error]Have no playerId in cLobbyMatch, session:%s", session.GetSessionId())
		}
	}

	// 匹配失败
	code := mproto.ECode_fail
	pb := &mproto.SLobbyMatch {
		Code:     &code,
	}
	return mproto.PackPBReply("s_lobby_match", pb)
}

//cLobbyCancel ： 取消匹配
func (self *Lobbys) cLobbyCancel(session gate.Session, msg []byte) ([]byte, string) {
	lobby, ok := self.GetLobby(0)
	if ok {
		playerId, err := strconv.Atoi(session.Get("PlayerId"))
		if err == nil {
			if lobby.UnMatch(playerId) {
				// 取消匹配成功
				code := mproto.ECode_ok
				pb := &mproto.SLobbyCancel {
					Code:     &code,
				}
				return mproto.PackPBReply("s_lobby_cancel", pb)
			}
		} else {
			log.Error("[Error]Have no playerId in cLobbyCancel, session:%s", session.GetSessionId())
		}
	}

	// 取消匹配失败
	code := mproto.ECode_fail
	pb := &mproto.SLobbyCancel {
		Code:     &code,
	}
	return mproto.PackPBReply("s_lobby_cancel", pb)
}
