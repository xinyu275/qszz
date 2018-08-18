package mgate

import (
	"bufio"
	"io"
	"encoding/binary"
	"qserver/mproto"
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
	return ProtoName, databuf, nil
}


