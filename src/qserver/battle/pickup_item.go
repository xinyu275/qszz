package battle

import (
	"math/rand"
)

//物品拾取操作
func (self *Room) doPickupItem(player *Player, item *Item) {
	//按物品类型分类
	// 1.血包(急救箱)/弹夹包/护盾等消耗性道具 判断是否可以直接拾取/如果可以拾取马上消耗掉
	// 2.随机包 那就随机一个增益道具，在判断是否可以拾取，如果可以拾取那就进背包
	// 3.拾取增益道具(2,4-10) 判断是否可以拾取，如果可以马上就进背包
	itemBaseId := item.BaseId
	if itemBaseId == ITEM_RANDOM {
		self.pickupRandomItem(player, item)
	} else if checkIsPlusItem(itemBaseId) {
		self.pickupPlusItem(player, item)
	} else if checkIsConsumeItem(itemBaseId) {
		self.pickupConsumeItem(player, item)
	}
}

//拾取随机包
func (self *Room) pickupRandomItem(player *Player, item *Item) {
	//1.打包随机包,随机一个增益型物品
	plusItems := getPlusItems()
	randIndex := rand.Intn(len(plusItems))
	randBaseId := plusItems[randIndex]
	//替换场景内物品的baseid
	item.UpdateBaseId(randBaseId)
	//如果不能捡起，那就要告诉所有人宝箱开出的物品
	if !self.pickupPlusItem(player, item) {
		self.msgs = append(self.msgs, packItem(self.view.GetMsgId(), item, false))
	}
}

//拾取增益型物品
//返回是否可以拾取
func (self *Room) pickupPlusItem(player *Player, item *Item) bool {
	//获取物品的叠加上限
	baseItem := GetBaseItem(item.BaseId)
	overlap := baseItem.Overlap
	//如果玩家自己拥有的达到叠加上限，那就不能拾取了
	haveNum := player.GetItemCountByBaseId(item.BaseId)
	if haveNum >= overlap {
		return false
	}
	addNum := item.Num
	if haveNum+item.Num > overlap {
		addNum = overlap - haveNum
	}
	//拾取物品
	player.AddItem(item.BaseId, addNum)
	//如果是获得能量增益道具，那就要发给前段
	if item.BaseId == ITEM_ADD_ENERGY {
		player.AddOps(packPlayerMp(self.view.GetMsgId(), player))
	}
	//场景里物品消失
	self.msgs = append(self.msgs, packItem(self.view.GetMsgId(), item, true))
	//删除场景内物品
	self.DelItem(item.Id)

	//返回玩家ops,告诉玩家你增加了某种类型物品多少个
	totalNum := player.GetItemCountByBaseId(item.BaseId)
	player.AddOps(packPickUpItemReply(self.view.GetMsgId(), item, totalNum))

	return true
}

//消耗性物品拾取
func (self *Room) pickupConsumeItem(player *Player, item *Item) {
	if self.doPickupConsumeItem(player, item) {
		//场景里物品消失
		self.msgs = append(self.msgs, packItem(self.view.GetMsgId(), item, true))
		//删除场景物品
		self.DelItem(item.Id)

	}
}
func (self *Room) doPickupConsumeItem(player *Player, item *Item) (isSuccess bool) {
	switch item.BaseId {
	case ITEM_BOLLD:
		//血包(急救箱) 加血，如果自己满血，那就不能拾取
		if player.Hp == player.HpMax {
			return false
		}
		baseItem := GetBaseItem(item.BaseId)
		randAddHp := Rand(baseItem.Min, baseItem.Max)
		player.AddHp(randAddHp)
		//告诉场景你属性修改了
		self.view.ChangePlayerExtraAttr(player)
		isSuccess = true
		return
	case ITEM_BULLET:
		//子弹夹 增加玩家当前枪的子弹
		//判断当前枪的子弹是否已满，如果满了那就不能再拾取
		if player.CurWeaponId == 0 {
			return
		}
		weapon := self.GetWeapon(player.CurWeaponId)
		if weapon == nil {
			return
		}
		BaseWeapon := GetBaseWeapon(weapon.BaseId)
		if weapon.BulletNum >= BaseWeapon.BulletMaxNum {
			return
		}
		//随机增加子弹数字
		baseItem := GetBaseItem(item.BaseId)
		AddBulletNum := Rand(baseItem.Min, baseItem.Max)
		if weapon.BulletNum+AddBulletNum > BaseWeapon.BulletMaxNum {
			AddBulletNum = BaseWeapon.BulletMaxNum - weapon.BulletNum
		}
		//增加玩家子弹数
		weapon.AddBullet(AddBulletNum)
		//告诉玩家自己增加了子弹数
		player.AddOps(packPickUpItemReply(self.view.GetMsgId(), item, weapon.BulletNum))

		isSuccess = true
		return
	case ITEM_SHIELDBUFFER:
		baseItem := GetBaseItem(item.BaseId)
		//护盾 增加玩家护盾值 替换
		if player.Shield == baseItem.Min {
			return
		}
		player.AddShield()
		//告诉场景你属性修改了
		self.view.ChangePlayerExtraAttr(player)
		isSuccess = true
		return
	}
	return
}

//判断物品是否是消耗性道具
func checkIsConsumeItem(itemBaseId int) bool {
	return itemBaseId == ITEM_BOLLD || itemBaseId == ITEM_BULLET || itemBaseId == ITEM_SHIELDBUFFER
}

//判断物品是否是增益型的
func checkIsPlusItem(itemBaseId int) bool {
	for _, v := range getPlusItems() {
		if v == itemBaseId {
			return true
		}
	}
	return false
}

//增益物品列表
func getPlusItems() []int {
	return []int{
		ITEM_DEFENSEBUFFER,
		ITEM_BULLETDISTANCE,
		ITEM_BULLETHURT,
		ITEM_BULLETSPEED,
		ITEM_WEAPONCD,
		ITEM_ADD_ENERGY,
		ITEM_UPMAXHPBUFFER,
	}
}
