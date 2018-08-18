package battle

import "encoding/binary"

//发送包处理
type Msg struct {
	msgId int
	mtype int //0.玩家移动 1.玩家附加 2.物品 3.武器 4.子弹 5.玩家操作
	body  []byte
}

//----------------打包--------------
//0.打包角色移动
func packMove(msgId int, player *Player) *Msg {
	//[id.8b|dir.8b|x.16b|y.16b]
	var w1 uint32 = uint32(player.Id) | uint32(player.Dir)<<8 | uint32(player.X)<<16
	var w2 uint16 = uint16(player.Y)
	return &Msg{
		msgId: msgId,
		mtype: 0,
		body:  append(pack4(w1), pack2(w2)...),
	}
}

//1.打包玩家附加属性
func packUnitExtra(msgId int, player *Player) *Msg {
	//玩家附加属性[id.8b|weapon.8b|wear.8b|shield.8b|hp.10b|hp_max.10b|itemcount.8b|预留.4b]
	var w1 uint32 = uint32(player.Id) | uint32(player.WeaponType)<<8 |
		uint32(player.WeaponType)<<16 | uint32(player.Shield)<<24
	var w2 uint32 = uint32(player.Hp) | uint32(player.HpMax)<<10 |
		uint32(player.BagItemsCount)<<20
	return &Msg{
		msgId: msgId,
		mtype: 1,
		body:  append(pack4(w1), pack4(w2)...),
	}
}

//打包物品信息。全地图广播
func packItem(msgId int, item *Item, isDel bool) *Msg {
	//物品[id.14|owner.8b|baseid.8b|state.2b|x.16b|y.16b]
	//state==1销毁物品  owner > 0说明被人捡起要在场景里消失 owner == 0 物品要重新回到场景中显示出来
	state := 0
	if isDel {
		state = 1
	}
	var w1 uint32 = uint32(item.Id) | uint32(item.Owner)<<14 | uint32(item.BaseId)<<22 |
		uint32(state)<<30
	var w2 uint32 = uint32(item.X) | uint32(item.Y)<<16
	return &Msg{
		msgId: msgId,
		mtype: 2,
		body:  append(pack4(w1), pack4(w2)...),
	}
}

//3.打包武器信息，全地图广播
func packWeapon(msgId int, weapon *Weapon, isDel bool) *Msg {
	//武器[id.14|owner.8b|state:10b|x.16b|y.16b]
	//state==1销毁武器  owner > 0说明被人捡起要在场景里消失 owner == 0 武器要重新回到场景中显示出来
	state := 0
	if isDel {
		state = 1
	}
	var w1 uint32 = uint32(weapon.Id) | uint32(weapon.Owner)<<14 |
		uint32(weapon.BaseId)<<22 | uint32(state)<<30
	var w2 uint32 = uint32(weapon.X) | uint32(weapon.Y)<<16
	return &Msg{
		msgId: msgId,
		mtype: 3,
		body:  append(pack4(w1), pack4(w2)...),
	}
}

//4.打包子弹开火生成一颗子弹
func packBullet(msgId int, bullet *Bullet) *Msg {
	//[id.14|dir.8b|speed.10b|x.16b|y.16b|weaponType.8b]
	var w1 uint32 = uint32(bullet.Id) | uint32(bullet.FireDir)<<14 |
		uint32(bullet.AddSpeed)<<22
	var w2 uint32 = uint32(bullet.StartX) | uint32(bullet.StartY)<<16
	var w3 uint16 = uint16(bullet.WeaponBaseId) | uint16(bullet.AddRange)<<8
	return &Msg{
		msgId: msgId,
		mtype: 4,
		body:  append(append(pack4(w1), pack4(w2)...), pack2(w3)...),
	}
}

//5.1.打包玩家开火操作返回包
func packFireOpReply(msgId int, bullet *Bullet, bulletNum int) *Msg {
	//[op.8b|id.14b|type.2b|val.8b]
	var w uint32 = uint32(OP_FIRE) | uint32(bullet.Id)<<8 | uint32(OBJ_BULLET)<<22 | uint32(bulletNum)<<24
	return &Msg{
		msgId: msgId,
		mtype: 5,
		body:  pack4(w),
	}
}

//5.2打包玩家换枪返回
func packSwitchOpReply(msgId int, player *Player) *Msg {
	var w uint32 = uint32(OP_SWITCHWEAPON) | uint32(player.CurWeaponId)<<8 | uint32(OBJ_WEAPON)<<22
	return &Msg{
		msgId: msgId,
		mtype: 5,
		body:  pack4(w),
	}
}

//5.3打包玩家拾取武器返回
func packPickUpWeaponReply(msgId int, weapon *Weapon) *Msg {
	var w uint32 = uint32(OP_PICKUP) | uint32(weapon.Id)<<8 | uint32(OBJ_WEAPON)<<22 |
		uint32(weapon.BulletNum)<<24
	return &Msg{
		msgId: msgId,
		mtype: 5,
		body:  pack4(w),
	}
}

//5.4打包玩家物品返回
func packPickUpItemReply(msgId int, item *Item, Num int) *Msg {
	//[op.8b|id.14b|type.2b|val.8b]
	var w uint32 = uint32(OP_PICKUP) | uint32(item.BaseId)<<8 | uint32(OBJ_ITEM)<<22 |
		uint32(Num)<<24
	return &Msg{
		msgId: msgId,
		mtype: 5,
		body:  pack4(w),
	}
}

//5.5打包玩家能量值
func packPlayerMp(msgId int, player *Player) *Msg {
	var w uint32 = uint32(OP_MP) | uint32(player.MpMax)<<8 | uint32(player.Mp)<<20
	return &Msg{
		msgId: msgId,
		mtype: 5,
		body:  pack4(w),
	}
}

func pack4(w uint32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, w)
	return buf
}

func pack2(w uint16) []byte {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, w)
	return buf
}
