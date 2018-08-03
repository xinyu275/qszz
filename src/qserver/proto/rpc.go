package proto

var (
	Proto_To_MsgId = map[string]int{
		"c_player_login":101,
		"s_player_login":102,
	}

	MsgId_To_Proto = map[int]string{
		101:"c_player_login",
		102:"s_player_login",
	}
)


func ProtoToMsgId(ProtoName string) int {
	return Proto_To_MsgId[ProtoName]
}

func MsgIdToProto(msgId int) string {
	return MsgId_To_Proto[msgId]
}