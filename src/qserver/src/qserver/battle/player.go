package battle

import (
	"github.com/liangdas/mqant/gate"
	"math"
	"strconv"
)

type Player struct {
	Session        gate.Session
	PlayerId       int
	RoleType       int
	Id             int // 房间唯一id
	X              int // 像素
	Y              int
	Dir            int // 方向
	Hp             int // 玩家血量
	HpMax          int // 玩家最大血量（默认血量+增益物品）
	Def            int // 玩家防御（增益物品）
	AddRange       int // 子弹射程增益，子弹的射程=子弹基础射程+AddRange
	AddHurt        int // 子弹伤害增益 子弹伤害=子弹基础伤害+AddHurt
	AddBulletSpeed int // 子弹移动速度增益 子弹速度=子弹基础速度+AddBulletSpeed
	Mp             int // 能量
	MpMax          int // 能量上限
	Weight         int // 重量(玩家+武器+增益物品)
	Shield         int // 玩家护盾
	Wear           int // 衣服
	WeaponType     int // 武器类型
	LX             int //本帧开始时的x
	LY             int //本帧开始时的y

	PlayerWeapons []int // 玩家枪的唯一id列表（默认1， >=2为场景生成的枪）
	CurWeaponId   int   // 玩家当前使用的枪
	FireLastTime  int64 // 枪上次射击时间(毫秒),换枪要重置时间(会有bug:玩家不断换枪再开枪，但是应该没事有子弹限制)

	BagItems      map[int]int // 玩家的背包物品 BaseId -> Num
	BagItemsCount int         // 玩家的背包物品数量统计

	Kill        int // 玩家击杀数
	MoveSpeedUp int // 玩家技能加速万分比
	MpFrame     int // 能量帧数(加速/回复能量/加速扣除能量/结束加速 会设置为当前帧)

	Ops    []*Msg //玩家操作包
	Loaded bool   // 前段是否已经加载完成
}

//房间内新建玩家
func NewPlayer(session gate.Session, room *Room) *Player {
	randomX, randomY := RandPlayerPos(room.mapId)
	playerId, _ := strconv.Atoi(session.Get("PlayerId"))
	return &Player{
		Session:    session,
		PlayerId:   playerId,
		RoleType:   ROLE_PLATER,
		Id:         room.player_objid,
		X:          randomX,
		Y:          randomY,
		Hp:         DEFAULT_HP,
		HpMax:      DEFAULT_HP,
		Def:        DEFAULT_DEF,
		Weight:     DEFAULT_WEIGHT,
		Mp:         DEFAULT_MP,
		MpMax:      DEFAULT_MP,
		Shield:     0,
		Wear:       1,
		WeaponType: 0,
		LX:         0,
		LY:         0,

		PlayerWeapons: []int{},
		CurWeaponId:   0,
		BagItems:      make(map[int]int),

		Ops:    nil,
		Loaded: false,
	}
}

//
func (self *Player) OnRequest(session gate.Session) {
	self.Session = session
}

//改变方向/停止
func (self *Player) ChangeDir(session gate.Session, dir int) {
	self.Session = session
	self.Dir = dir
}

//玩家每帧处理移动
func (self *Player) Move(mapId int) {
	dirAngle := (self.Dir - 1) * 6
	//if self.Dir == 0 {
	//	dirAngle = 360
	//}
	//加速
	SpeedUp := float64(10000)
	if self.MoveSpeedUp > 0 {
		SpeedUp = float64(self.MoveSpeedUp)
	}
	x := self.X + int(math.Cos(float64(180-dirAngle)*(math.Pi/180.0))*0.03*float64(DEFAULT_SPEED)*
		(SpeedUp/10000)/(1+float64(self.Weight)/200))
	y := self.Y + int(math.Sin(float64(180-dirAngle)*(math.Pi/180.0))*0.03*float64(DEFAULT_SPEED)*
		(SpeedUp/10000)/(1+float64(self.Weight)/200))
	if CanWalk(mapId, x, y) {
		self.X = x
		self.Y = y
	} else if CanWalk(mapId, self.X, y) {
		self.Y = y
	} else if CanWalk(mapId, x, self.Y) {
		self.X = x
	}
}

//玩家换枪
func (self *Player) SwitchWeapon(weaponId int, weaponType int) {
	self.CurWeaponId = weaponId
	self.WeaponType = weaponType
	//重新计算重量属性
	self.CalcWeight()
	//重置上次开火时间
	self.SetFireTime(0)
}

//玩家拾取枪
func (self *Player) PickUpWeapon(weaponId int, weaponType int) {
	//如果只有1吧枪，那拾取到的枪进被被用枪，如果玩家没有或者已经有2吧，那替换当前的枪
	if len(self.PlayerWeapons) == 1 {
		//那拾取的枪进备用枪
		self.PlayerWeapons = append(self.PlayerWeapons, weaponId)
		return
	}
	//玩家枪的列表要替换
	newWeapons := make([]int, 0, 2)
	for _, weaponId := range self.PlayerWeapons {
		if weaponId != self.CurWeaponId {
			newWeapons = append(newWeapons, weaponId)
		}
	}
	self.PlayerWeapons = append(newWeapons, weaponId)
	self.SwitchWeapon(weaponId, weaponType)
}

//判断玩家是否拥有这把枪
func (self *Player) CheckHaveWeapon(weaponId int) bool {
	have := false
	for _, v := range self.PlayerWeapons {
		if v == weaponId {
			have = true
		}
	}
	return have
}

//玩家扣血
func (self *Player) Hurt(hurtHp int) {
	if self.Shield >= hurtHp {
		self.Shield -= hurtHp
		return
	}
	self.Shield = 0
	self.Hp -= (hurtHp - self.Shield)
	if self.Hp < 0 {
		self.Hp = 0
	}
}

//获取某个物品的个数
func (self *Player) GetItemCountByBaseId(baseId int) int {
	num := self.BagItems[baseId]
	return num
}

//获取玩家物品数量
func (self *Player) GetBagItemsCount() int {
	return self.BagItemsCount
}

//获取增益型物品
func (self *Player) AddItem(baseId int, addNum int) {
	oldv := self.BagItems[baseId]
	self.BagItems[baseId] = oldv + addNum
	self.BagItemsCount += addNum
	//增加加强
	self.ItemCalcAttr(baseId, addNum)
}

//玩家死亡处理
func (self *Player) Die() {
	//清空增益
	self.BagItems = map[int]int{}
	self.BagItemsCount = 0
	//TODO 死亡玩家属性是否要重置
}

//玩家拾取急救箱加血
func (self *Player) AddHp(addHp int) {
	newHp := self.Hp + addHp
	if newHp > self.HpMax {
		newHp = self.HpMax
	}
	self.Hp = newHp
}

//增加护盾(替换)
func (self *Player) AddShield() {
	baseItem := GetBaseItem(ITEM_SHIELDBUFFER)
	//这里不要随机了，min==max
	self.Shield = baseItem.Min
}

//添加玩家击杀数
func (self *Player) AddKill() {
	self.Kill += 1
}

//添加操作
func (self *Player) AddOps(msg *Msg) {
	self.Ops = append(self.Ops, msg)
}

//增加物品时玩家属性增益会增加
func (self *Player) ItemCalcAttr(baseId, addNum int) {
	baseItem := GetBaseItem(baseId)
	randV := Rand(baseItem.Min, baseItem.Max)
	switch baseId {
	case ITEM_DEFENSEBUFFER:
		self.Def += randV * addNum
	case ITEM_BULLETDISTANCE:
		self.AddRange += randV * addNum
	case ITEM_BULLETHURT:
		self.AddHurt += randV * addNum
	case ITEM_BULLETSPEED:
		self.AddBulletSpeed += randV * addNum
	case ITEM_ADD_ENERGY:
		self.MpMax += randV * addNum
		self.Mp += randV * addNum
	case ITEM_UPMAXHPBUFFER:
		self.HpMax += randV * addNum
		self.Hp += randV * addNum
	default:
	}
	self.Weight += baseItem.Weight * addNum
}

//重新计算weight属性
func (self *Player) CalcWeight() {
	weight := DEFAULT_WEIGHT
	//增益物品重要
	for baseId, num := range self.BagItems {
		baseItem := GetBaseItem(baseId)
		weight += baseItem.Weight * num
	}
	//武器重量
	baseWeapon := GetBaseWeapon(self.WeaponType)
	if baseWeapon == nil {
		self.Weight = weight
	}
	self.Weight = weight + baseWeapon.Weight
}

//玩家使用加速技能
func (self *Player) SpeedUp(frame int) {
	baseSkill := GetBaseSkill(SKILL_UP_SPEED)
	self.MoveSpeedUp = baseSkill.Val
	self.SetMpFrame(frame)
}

//玩家停止加速
func (self *Player) StopSpeedUp(curFrame int) {
	self.MoveSpeedUp = 0
	self.SetMpFrame(curFrame)
}

//设置能量帧数
func (self *Player) SetMpFrame(curFrame int) {
	self.MpFrame = curFrame
}

//玩家能量扣除（加速状态）和回复
func (self *Player) UpdateMp(curFrame int) {
	if self.Hp == 0 || !self.Loaded {
		return
	}
	if curFrame-self.MpFrame < MP_CHANGE_INTERVAL_FRAME {
		return
	}
	if self.MoveSpeedUp > 0 {
		//扣除魔法
		self.Mp -= GetBaseSkill(SKILL_UP_SPEED).CostMp
		if self.Mp <= 0 {
			self.Mp = 0
			self.StopSpeedUp(curFrame)
		}
	} else if self.Mp < self.MpMax {
		//回复mp
		self.Mp += RECOVER_MP
		if self.Mp > self.MpMax {
			self.Mp = self.MpMax
		}
	}
	self.SetMpFrame(curFrame)
}

//设置上次开火时间(毫秒)
func (self *Player) SetFireTime(curMillsSecond int64) {
	self.FireLastTime = curMillsSecond
}
