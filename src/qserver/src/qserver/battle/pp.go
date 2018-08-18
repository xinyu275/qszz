package battle

import (
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"github.com/liangdas/mqant/gate"
	"qserver/mproto"
)

//每帧数据
func (b *Battle) cPlayerFrame(session gate.Session, body []byte) ([]byte, string) {
	msg := &mproto.CPlayerFrame{}
	err := proto.Unmarshal(body, msg)
	if err != nil {
		return nil, err.Error()
	}
	//1.低6位为方向
	//2.7-16位为操作
	//高16位自定义解析
	Bytes := msg.GetOp()
	var frameInt uint32 = 0
	if len(Bytes) == 1 {
		//只有方向
		frameInt = uint32(Bytes[0])
	} else if len(Bytes) == 2 {
		frameInt = uint32(binary.BigEndian.Uint16(Bytes))
	} else if len(Bytes) == 4 {
		frameInt = binary.BigEndian.Uint32(Bytes)
	}
	//移动
	b.changeDir(session, frameInt)

	//没有op
	if frameInt&0xFFC0 == 0 {
		return nil, ""
	}
	if frameInt&BtnFire != 0 { //开火
		b.fire(session, frameInt)
	} else if frameInt&BtnSwitchWeapon != 0 { //换枪
		b.switch_weapon(session, frameInt)
	} else if frameInt&BtnPickUp != 0 { //拾取
		b.pickup(session, frameInt)
	} else if frameInt&BtnCast != 0 {
		b.cast(session, frameInt)
	}
	//丢弃
	//if frameInt&BtnDrop != 0 {
	//	b.drop(session, frameInt)
	//}
	//开启倍镜
	if frameInt&BtnXTimesMirror != 0 {
		b.timesMirror(session, frameInt)
	}
	//使用道具
	//if frameInt&BtnUseItem != 0 {
	//	b.useItem(session, frameInt)
	//}
	return nil, ""
}

//玩家请求开始
func (b *Battle) cBattleStart(session gate.Session, body []byte) ([]byte, string) {
	b.join(session)
	return nil, ""
}

//场景加载完成
func (b *Battle) cBattleLoaded(session gate.Session, body []byte) ([]byte, string) {
	b.room.PutQueue("Loaded", session)
	return nil, ""
}

//发射者通知服务器谁被击中了
func (b *Battle) cBattleHurt(session gate.Session, body []byte) ([]byte, string) {
	msg := &mproto.CBattleHurt{}
	err := proto.Unmarshal(body, msg)
	if err != nil {
		return nil, err.Error()
	}
	bulletId := int(msg.GetBulletId())
	playerBattleId := int(msg.GetPlayerBattleId())
	b.room.PutQueue("Hurt", session, bulletId, playerBattleId)
	return nil, ""
}

//------------分解协议的操作--------------------------
//玩家加入
func (b *Battle) join(session gate.Session) {
	b.room.PutQueue("Join", session)
}

//玩家改变改变移动方向/停止移动
func (b *Battle) changeDir(session gate.Session, frameInt uint32) {
	//低6位为方向
	dir := int(frameInt & 0x3f)
	b.room.PutQueue("ChangeDir", session, dir)
}

//开火
func (b *Battle) fire(session gate.Session, frameInt uint32) {
	//高16位为方向
	fireDir := int(frameInt >> 16)
	b.room.PutQueue("Fire", session, fireDir)
}

//换枪
func (b *Battle) switch_weapon(session gate.Session, frameInt uint32) {
	//高16位为枪的id
	weaponId := int(frameInt >> 16)
	b.room.PutQueue("SwitchWeapon", session, weaponId)
}

//拾取
func (b *Battle) pickup(session gate.Session, frameInt uint32) {
	//拾取id
	pickId := int((frameInt << 2) >> 18)
	pickType := int(frameInt >> 30)
	if pickId > 0 && (pickType == OBJ_WEAPON || pickType == OBJ_ITEM) {
		b.room.PutQueue("Pickup", session, pickId, pickType)
	}

}

//释放技能
func (b *Battle) cast(session gate.Session, frameInt uint32) {
	//高16位为技能id
	skillId := int(frameInt >> 16)
	b.room.PutQueue("Cast", session, skillId)
}

//丢弃
//func (b *Battle) drop(session gate.Session, frameInt uint32) {
//
//}

//开启倍镜
func (b *Battle) timesMirror(session gate.Session, frameInt uint32) {

}

//使用道具
//func (b *Battle) useItem(session gate.Session, frameInt uint32) {
//
//}
