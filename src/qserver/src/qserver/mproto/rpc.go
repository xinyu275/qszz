package mproto

import (
	"encoding/binary"
	"errors"
	"github.com/golang/protobuf/proto"
)

var (
	//协议号对应msgid
	Proto_To_MsgId = map[string]int{
		"c_player_login": 1001,
		"s_player_login": 1002,
		"c_server_time":  1003,
		"s_server_time":  1004,

		"c_player_info": 2001,
		"s_player_info": 2002,

		"c_player_frame":  3001,
		"s_player_frame":  3002,
		"c_battle_start":  3003,
		"s_battle_info":   3004,
		"c_battle_loaded": 3005,
		"c_battle_hurt":   3006,
	}
)
var (
	// msgid/100 对应模块名
	PerMsgId_To_Module = map[int]string{
		10: "Login",
		20: "User",
		30: "Battle",
	}
)

//根据Proto_To_MsgId生成MsgId_To_Proto
var MsgId_To_Proto map[int]string = make(map[int]string)

func init() {
	for protoName, msgId := range Proto_To_MsgId {
		MsgId_To_Proto[msgId] = protoName
	}
}

func ProtoToMsgId(ProtoName string) int {
	return Proto_To_MsgId[ProtoName]
}

func MsgIdToProto(msgId int) string {
	return MsgId_To_Proto[msgId]
}

//定义协议对应的模块名(确保模块名已经配置)
func ProtoToModule(ProtoName string) (string, error) {
	msgId := ProtoToMsgId(ProtoName)
	if module, ok := PerMsgId_To_Module[msgId/100]; ok {
		return module, nil
	}
	return "", errors.New("PerMsgId_To_Module not register " + ProtoName)
}

//打包协议远程返回数据
func PackPBReply(protoName string, pb proto.Message) ([]byte, string) {
	b, err := proto.Marshal(pb)
	if err != nil {
		return nil, err.Error()
	}
	return Pack(protoName, b), ""
}

//封装包体
func Pack(protoName string, b []byte) []byte {
	buf := make([]byte, 6+len(b))
	msgid := ProtoToMsgId(protoName)
	binary.BigEndian.PutUint16(buf, (uint16)(msgid))
	binary.BigEndian.PutUint32(buf[2:], (uint32)(len(b)))
	copy(buf[6:], b)
	return buf
}
