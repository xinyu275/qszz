package main

import (
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"qserver/mproto"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", ":3854")
	if err != nil {
		panic(err)
	}
	go run(conn)

	time.Sleep(10 * time.Second)
}

func run(conn net.Conn) {
	fmt.Println("start")
	login := &mproto.CPlayerLogin{
		Token: proto.String("aaaa"),
		Pf:    proto.String("pftest"),
	}
	loginBuf, err := proto.Marshal(login)
	fmt.Println(loginBuf)
	if err != nil {
		panic(err)
	}
	b := pack("c_player_login", loginBuf)
	conn.Write(b)

	//获取玩家信息
	playerProto := &mproto.CPlayerInfo{}

	playerInfo, err := proto.Marshal(playerProto)
	if err != nil {
		panic(err)
	}
	b = pack("c_player_info", playerInfo)
	conn.Write(b)

	//帧信息
	frame := &mproto.CPlayerFrame{
		Op: []byte{1},
	}
	frameInfo, err := proto.Marshal(frame)
	b = pack("c_player_frame", frameInfo)
	conn.Write(b)

	sProto := &mproto.CBattleStart{}

	startInfo, err := proto.Marshal(sProto)
	if err != nil {
		panic(err)
	}
	b = pack("c_battle_start", startInfo)
	conn.Write(b)

}

func pack(protoName string, body []byte) []byte {
	buf := make([]byte, 6+len(body))
	msgid := mproto.ProtoToMsgId(protoName)
	binary.BigEndian.PutUint16(buf, (uint16)(msgid))
	binary.BigEndian.PutUint32(buf[2:], (uint32)(len(body)))
	copy(buf[6:], body)
	return buf
}
