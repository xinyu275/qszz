package battle

import (
	"github.com/golang/protobuf/proto"
	"github.com/liangdas/mqant/gate"
	"qserver/battle/battleconf"
	"qserver/common"
	"qserver/mproto"
)

//房间操作处理
//玩家加入(加入返回s_battle_info is_start = 0,为来让前段有足够的时间加载，玩家不加入九宫格)
func (self *Room) join(session gate.Session) {
	var player *Player
	player = self.GetPlayer(session)
	if player == nil {
		//新建玩家
		player = NewPlayer(session, self)
		//为玩家创建一把枪
		self.createPlayerWeapon(player)
		self.player_objid++
		self.playerid_to_objid[player.PlayerId] = player.Id
		self.players[player.Id] = player
	} else {
		player.OnRequest(session)
		//玩家重载九宫格数据
		player.LX = 0
		player.LY = 0
		player.Loaded = false
	}
	//告诉玩家你的虚拟id以及场景元素
	pb := &mproto.SBattleInfo{
		Bid:       proto.Int32(int32(player.Id)),
		List:      self.getRoomObjList(false),
		IsStart:   proto.Int32(0),
		TimeStart: proto.Uint32(uint32(self.start_time)),
	}
	b, err := proto.Marshal(pb)
	if err == nil {
		session.Send("s_battle_info", b)
	}
	//告诉视野管理器
	//self.view.AddPlayer(player)
}

//为每个玩家默认创建一把枪
func (self *Room) createPlayerWeapon(player *Player) {
	//创建一把枪，并且默认被玩家拾取
	weapon := self.AddWeapon(DEFAULT_WEAPON, player.X, player.Y)
	weapon.PickUp(player)
	player.PickUpWeapon(weapon.Id, weapon.BaseId)
	self.msgs = append(self.msgs, packWeapon(self.view.GetMsgId(), weapon, false))
}

//前段通知后端场景数据已经加载完成
func (self *Room) loaded(session gate.Session) {
	player := self.GetPlayer(session)
	if player == nil {
		//异常情况不应该出现
		return
	}
	player.Loaded = true

	//告诉玩家你的虚拟id以及场景元素
	pb := &mproto.SBattleInfo{
		Bid:        proto.Int32(int32(player.Id)),
		List:       self.getRoomObjList(true),
		IsStart:    proto.Int32(1),
		OwnWeapons: self.packPlayerWeapon(player),
		OwnItems:   self.packPlayerItem(player),
		TimeStart:  proto.Uint32(uint32(self.start_time)),
	}
	b, err := proto.Marshal(pb)
	if err == nil {
		session.Send("s_battle_info", b)
	}
	//告诉视野管理器
	self.view.AddPlayer(player)
	//第一帧要把自己的私有属性告诉玩家
	player.AddOps(packPlayerMp(self.view.GetMsgId(), player))

	return
}

//打包s_battle_info
func (self *Room) getRoomObjList(loaded bool) []*mproto.PBattleObj {
	r := make([]*mproto.PBattleObj, 0)
	//物品
	for _, item := range self.items {
		r = append(r, &mproto.PBattleObj{
			Id:     proto.Uint32(uint32(item.Id)),
			BaseId: proto.Uint32(uint32(item.BaseId)),
			Type:   mproto.EBattleObjType(OBJ_ITEM).Enum(),
			X:      proto.Uint32(uint32(item.X)),
			Y:      proto.Uint32(uint32(item.Y)),
			Owner:  proto.Uint32(uint32(item.Owner)),
		})
	}
	//武器
	for _, weapon := range self.weapons {
		r = append(r, &mproto.PBattleObj{
			Id:     proto.Uint32(uint32(weapon.Id)),
			BaseId: proto.Uint32(uint32(weapon.BaseId)),
			Type:   mproto.EBattleObjType(OBJ_WEAPON).Enum(),
			X:      proto.Uint32(uint32(weapon.X)),
			Y:      proto.Uint32(uint32(weapon.Y)),
			Owner:  proto.Uint32(uint32(weapon.Owner)),
		})
	}
	//第二次sbattleinfo不需要传所有玩家信息，因为第一帧会发给玩家当前视野的玩家列表
	if loaded {
		return r
	}
	//玩家
	for _, player := range self.players {
		r = append(r, &mproto.PBattleObj{
			Id:     proto.Uint32(uint32(player.Id)),
			Type:   mproto.EBattleObjType(OBJ_UNIT).Enum(),
			BaseId: proto.Uint32(0),
			X:      proto.Uint32(uint32(player.X)),
			Y:      proto.Uint32(uint32(player.Y)),
		})
	}
	return r
}

//打包玩家枪的信息
func (self *Room) packPlayerWeapon(player *Player) []*mproto.PWeapon {
	b := make([]*mproto.PWeapon, 0)
	for _, weaponId := range player.PlayerWeapons {
		isEquip := 0
		if weaponId == player.CurWeaponId {
			isEquip = 1
		}
		weapon := self.GetWeapon(weaponId)
		if weapon == nil {
			continue
		}
		bulletNum := weapon.BulletNum
		b = append(b, &mproto.PWeapon{
			WeaponId:  proto.Int32(int32(weaponId)),
			IsEquip:   proto.Int32(int32(isEquip)),
			BulletNum: proto.Int32(int32(bulletNum)),
		})
	}
	return b
}

//打包玩家的物品列表
func (self *Room) packPlayerItem(player *Player) []*mproto.PItem {
	b := make([]*mproto.PItem, 0)
	for baseId, num := range player.BagItems {
		b = append(b, &mproto.PItem{
			BaseId: proto.Int32(int32(baseId)),
			Num:    proto.Int32(int32(num)),
		})
	}
	return b
}

//玩家改变方向/停止
func (self *Room) changeDir(session gate.Session, dir int) {
	p := self.GetPlayer(session)
	if p == nil {
		return
	}
	oldDir := p.Dir
	p.ChangeDir(session, dir)

	//old_dir 判断
	if dir == 0 && oldDir != 0 {
		//通知九宫格您停止移动了
		self.view.PlayerStopMove(p)
	}
}

//玩家开火
func (self *Room) fire(session gate.Session, fireDir int) {
	player := self.GetPlayer(session)
	if player == nil {
		return
	}
	//没枪了
	if player.CurWeaponId == 0 {
		return
	}
	//获取枪
	weapon := self.GetWeapon(player.CurWeaponId)
	if weapon == nil {
		return
	}
	if weapon.BulletNum <= 0 {
		return
	}
	baseWeapon := GetBaseWeapon(weapon.BaseId)
	//判断枪的cd时间
	curMilliSecond := common.GetMillUnixTime()
	if curMilliSecond-player.FireLastTime < int64(baseWeapon.Cd*1000) {
		return
	}
	//子弹数量-1
	weapon.Fire()

	//生成子弹，放入九宫格
	bullet := NewBullet(self, player, fireDir)
	self.bullet_id++
	//子弹数量太大就重置为1,子弹id用14位表示
	if self.bullet_id >= MAX_BULLET_ID {
		self.bullet_id = 1
	}
	self.bullets[bullet.Id] = bullet

	//告诉九宫格发射子弹了
	self.view.Fire(bullet)
	//玩家自己操作返回
	player.AddOps(packFireOpReply(self.view.GetMsgId(), bullet, weapon.BulletNum))
	player.SetFireTime(curMilliSecond)
}

//换枪
func (self *Room) switchWeapon(session gate.Session, weaponId int) {
	player := self.GetPlayer(session)
	if player == nil {
		return
	}
	//判断玩家是否拥有这只枪
	if !player.CheckHaveWeapon(weaponId) {
		return
	}
	weapon := self.GetWeapon(weaponId)
	if weapon == nil {
		return
	}
	//前面已经判断了武器已经有了
	weaponType := self.GetWeapon(weaponId).BaseId
	//玩家换枪
	player.SwitchWeapon(weaponId, weaponType)
	//九宫格告诉其他玩家你换枪了
	self.view.SwitchWeapon(player)
	//玩家自己操作返回
	player.AddOps(packSwitchOpReply(self.view.GetMsgId(), player))
}

//拾取
func (self *Room) pickup(session gate.Session, pickId, pickType int) {
	if pickType == OBJ_ITEM {
		self.pickupItem(session, pickId)
	} else {
		self.pickWeapon(session, pickId)
	}
}

//拾取物品
func (self *Room) pickupItem(session gate.Session, itemId int) {
	item := self.GetItem(itemId)
	if item == nil {
		return
	}
	player := self.GetPlayer(session)
	if player == nil {
		return
	}
	self.doPickupItem(player, item)
}

//拾取武器
func (self *Room) pickWeapon(session gate.Session, weaponId int) {
	weapon := self.GetWeapon(weaponId)
	if weapon == nil {
		return
	}
	//已经有玩家捡了
	if weapon.Owner > 0 {
		return
	}
	player := self.GetPlayer(session)
	if player == nil {
		return
	}
	//设置武器所有者
	weapon.PickUp(player)
	//场景广播
	self.msgs = append(self.msgs, packWeapon(self.view.GetMsgId(), weapon, false))

	//记录玩家当前的枪
	curWeaponId := player.CurWeaponId
	//马上装备/如果玩家身上枪<2,那就进备用
	player.PickUpWeapon(weapon.Id, weapon.BaseId)
	//告诉玩家您拾取了某个武器
	player.AddOps(packPickUpWeaponReply(self.view.GetMsgId(), weapon))

	if player.CurWeaponId == curWeaponId {
		//枪进备用了
		return
	}
	//九宫格告诉其他玩家你换枪了
	self.view.SwitchWeapon(player)

	//玩家以前的枪，如果子弹数 >0，那就放入场景中
	oldWeapon := self.GetWeapon(curWeaponId)
	if oldWeapon == nil {
		return
	}
	if oldWeapon.BulletNum <= 0 {
		//删除枪信息
		self.DelWeapon(curWeaponId)
		//前段销毁对象
		self.msgs = append(self.msgs, packWeapon(self.view.GetMsgId(), oldWeapon, true))
		return
	}
	//玩家丢弃以前的枪
	oldWeapon.Discard(player)
	//场景广播,旧枪要洒落在场景中
	self.msgs = append(self.msgs, packWeapon(self.view.GetMsgId(), oldWeapon, false))
}

//掉血
func (self *Room) hurt(session gate.Session, bulletId, objId int) {
	player := self.GetPlayer(session)
	if player == nil {
		return
	}
	if player.Id == objId || objId <= 0 || bulletId <= 0 {
		return
	}
	//判断对方玩家是否存在
	objPlayer := self.GetPlayerById(objId)
	if objPlayer == nil {
		return
	}
	//玩家已经死亡
	if objPlayer.Hp <= 0 {
		return
	}
	//判断子弹是否存在
	bullet := self.GetBullet(bulletId)
	if bullet == nil {
		return
	}
	//只有发射者可以自己子弹射中的协议信息
	if bullet.Owner != player.Id {
		return
	}
	//判断对方是否已经中过
	if !self.checkCanHurt(bullet, objPlayer) {
		return
	}
	self.doHurt(player, objPlayer, bullet)
}

//伤害处理
func (self *Room) doHurt(player, objPlayer *Player, bullet *Bullet) {
	//计算伤害 = 子弹伤害*(1-(1-1/(1+总防御/100)))
	//子弹伤害 = 玩家伤害增益 + 武器基础伤害
	baseWeapon := GetBaseWeapon(bullet.WeaponBaseId)
	hurt := int((bullet.AddHurt + baseWeapon.Hurt) * (1 - (1 - 1/(1+objPlayer.Def/100))))
	//子弹射中人数+1
	bullet.AddHurtObjNum()

	//玩家扣血处理
	self.doPlayerHurt(player, objPlayer, hurt)
}

//玩家扣血处理
func (self *Room) doPlayerHurt(player, objPlayer *Player, hurt int) {
	objPlayer.Hurt(hurt)
	if objPlayer.Hp <= 0 {
		if player != nil {
			player.AddKill()
		}
		self.dieDrop(objPlayer)
		objPlayer.Die()
	}
	//告诉视野管理器你属性改变了(目标对象的九宫格视野)
	self.view.ChangePlayerExtraAttr(objPlayer)
}

//判断子弹是否可以攻击玩家
func (self *Room) checkCanHurt(bullet *Bullet, objPlayer *Player) bool {
	//判断子弹打了多个玩家了，手枪只能打1个玩家，散弹最多5个
	baseWeapon := GetBaseWeapon(bullet.WeaponBaseId)
	if bullet.HurtObjNum >= baseWeapon.PerBulletNum {
		return false
	}
	//TODO 判断玩家是否可以被子弹击中(作弊判断)
	return true
}

//玩家死亡掉落处理
func (self *Room) dieDrop(diePlayer *Player) {
	if diePlayer.BagItemsCount == 0 {
		return
	}
	//如果物品数量>DROP_MAX_PLUSITEM,那就要考虑堆叠掉落
	DropNum := DROP_MAX_PLUSITEM
	if diePlayer.BagItemsCount < DROP_MAX_PLUSITEM {
		DropNum = diePlayer.BagItemsCount
	}
	posList := self.genDropList(diePlayer.X, diePlayer.Y, DropNum)
	lenPosList := len(posList)
	var dropList []*Item = make([]*Item, 0, DropNum)
	curPos := 0
	if diePlayer.BagItemsCount <= DropNum {
		//那一个一个散落
		for baseId, num := range diePlayer.BagItems {
			for i := 0; i < num; i++ {
				pos := posList[curPos]
				item := self.AddItem(baseId, 1, pos.X, pos.Y)
				dropList = append(dropList, item)

				curPos++
				if curPos >= lenPosList {
					break
				}
			}
		}
	} else {
		//部分物品要叠加
		leftOverlapNum := diePlayer.BagItemsCount - DropNum
		//5个为一组扣除，物品不会太多，可以应该扣完
		groupOverlap := 5
		for baseId, num := range diePlayer.BagItems {
			for num > 0 {
				overlapNum := 1
				if leftOverlapNum <= 0 || num == 1 {
					overlapNum = 1
				} else if leftOverlapNum >= groupOverlap && num >= groupOverlap {
					overlapNum = groupOverlap
				} else if leftOverlapNum >= groupOverlap {
					overlapNum = num
				} else if num >= groupOverlap {
					overlapNum = leftOverlapNum + 1
				}
				num -= overlapNum
				leftOverlapNum -= (overlapNum - 1)
				//打包物品
				pos := posList[curPos]
				item := self.AddItem(baseId, overlapNum, pos.X, pos.Y)
				dropList = append(dropList, item)

				curPos++
				if curPos >= lenPosList {
					break
				}
			}
		}
	}
	//告诉场景新增了物品了
	for _, item := range dropList {
		self.msgs = append(self.msgs, packItem(self.view.GetMsgId(), item, false))
	}
}

//生成一份可以掉落的坐标列表,只要 DROP_MAX_PLUSITEM个
func (self *Room) genDropList(x, y, num int) []battleconf.PosXY {
	maxCell64Distance := 50
	//图片大小==前端的显示的一格(64 *64),因此掉落位置一定要是32的基数倍
	clientCellX := 64
	clientCellY := 64
	playerX64 := x / clientCellX
	playery64 := y / clientCellY
	x32 := 0
	y32 := 0
	var posList []battleconf.PosXY = make([]battleconf.PosXY, 0, num)
	for dis := 2; dis <= maxCell64Distance; dis += 2 {
		for xdis := 0; xdis <= dis; xdis++ {
			ydis := dis - xdis
			x32 = (playerX64+xdis)*clientCellX + clientCellX/2
			y32 = (playery64+ydis)*clientCellY + clientCellY/2
			//判断坐标是否可以走
			if CanWalk(self.mapId, x32, y32) {
				posList = append(posList, battleconf.PosXY{x32, y32})
				if len(posList) >= num {
					return posList
				}
			}
		}
	}
	return posList
}

//使用技能cast
func (self *Room) cast(session gate.Session, skillId int) {
	player := self.GetPlayer(session)
	if player == nil {
		return
	}
	self.castSkill(player, skillId)
}
