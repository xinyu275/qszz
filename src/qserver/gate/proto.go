package mgate

import (
	"bufio"
	"io"
	"encoding/binary"
	mproto "qserver/proto"

	"github.com/golang/protobuf/proto"
)

//proto解析协议
func ReadPack(r *bufio.Reader) (ProtoName string, data []byte, err error) {
	//2字节msgid + 4字节长度 + protodata

	//msgid
	msgbuf := make([]byte, 2)
	_, err = io.ReadFull(r, msgbuf)
	if err != nil {
		return
	}
	msgid := binary.BigEndian.Uint16(msgbuf)

	ProtoName = mproto.MsgIdToProto(int(msgid))

	//len
	lenbuf := make([]byte, 4)
	_, err = io.ReadFull(r, lenbuf)
	if err != nil {
		return
	}
	len := binary.BigEndian.Uint32(lenbuf)

	//data
	databuf := make([]byte, len)
	_, err = io.ReadFull(r, databuf)
	if err != nil {
		return
	}
	return
}

func PackSPlayerLogin() []byte{
	result := mproto.EShinResult_LOGIN_SUCCESS
	SLogin := &mproto.SPlayerLogin{
		Result:&result,
		PlayerId:proto.Uint32(1),
	}
	b, _ := proto.Marshal(SLogin)

	buf := make([]byte, 6 + len(b))
	msgid := mproto.ProtoToMsgId("s_player_login")
	binary.BigEndian.PutUint16(buf, (uint16)(msgid))
	binary.BigEndian.PutUint32(buf[2:], (uint32)(len(b)))
	copy(buf[6:], b)
	return buf
}
