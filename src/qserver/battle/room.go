package battle

import (
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/module"
	"math/rand"
	"qserver/battle/battleconf"
	"qserver/common"
	"strconv"
	"time"
)

type Room struct {
	BaseRoom
	module        module.Module
	roomId        int // 房间id
	current_frame int // 当前帧
	start_time    int // 开始时间(秒时间抽)
	stoped        bool
	mapId         int
	view          *View // 九宫格视野管理

	//角色
	players           map[int]*Player //房间objid -> 玩家信息
	playerid_to_objid map[int]int     //玩家id -> objid
	player_objid      int

	//道具
	items   map[int]*Item
	item_id int

	//武器
	weapons   map[int]*Weapon
	weapon_id int

	//子弹唯一id
	bullets   map[int]*Bullet
	bullet_id int

	//安全区域
	safe_left_x           int  // 安全区域左下角位置x
	safe_left_y           int  // 安全区域左下角位置y
	safe_right_x          int  // 安全区域右上角位置x
	safe_right_y          int  // 安全区域右上角位置y
	safe_index            int  // 安全区域索引
	safe_index_start_time int  // 开始时间抽(秒)
	safe_end              bool // 调整是否结束
	//全局消息
	msgs []*Msg
}

//新建房间
func NewRoom(module module.Module, mapId int, roomId int) *Room {
	room := &Room{
		module:        module,
		roomId:        roomId,
		current_frame: 0,
		start_time:    common.GetUnixTime(),
		stoped:        true,
		mapId:         mapId,
		view:          NewView(),

		player_objid:      1,
		players:           make(map[int]*Player),
		playerid_to_objid: make(map[int]int),
		msgs:              make([]*Msg, 0),

		items:     make(map[int]*Item),
		item_id:   1,
		weapons:   make(map[int]*Weapon),
		weapon_id: 2,
		bullets:   make(map[int]*Bullet),
		bullet_id: 1,

		//默认全部都是安全区域
		safe_right_x: MapWidth,
		safe_right_y: MapHeight,
	}
	room.Init()
	room.initItems(5)
	room.initWeapon(5)

	room.Register("ChangeDir", room.changeDir)
	room.Register("Fire", room.fire)
	room.Register("SwitchWeapon", room.switchWeapon)
	room.Register("Pickup", room.pickup)
	room.Register("Hurt", room.hurt)
	room.Register("Cast", room.cast)
	room.Register("Join", room.join)
	room.Register("Loaded", room.loaded)
	//注册函数
	return room
}

//Start: 房间开启
func (self *Room) Start() {
	if self.stoped {
		self.stoped = false
		go func() {
			rand.Seed(time.Now().UnixNano())
			tick := time.NewTicker(30 * time.Millisecond)
			secondTick := time.NewTicker(time.Second)
			defer func() {
				tick.Stop()
				secondTick.Stop()
			}()
			for !self.stoped {
				select {
				case <-tick.C:
					self.Update()
				case <-secondTick.C:
					self.UpdateSec()
				}
			}
		}()
	}
}

func (self *Room) Stop() {
	self.stoped = true
}

//每帧执行
func (self *Room) Update() {
	self.current_frame++
	//处理队列
	self.ExecuteEvent()

	//位置计算
	for _, player := range self.players {
		if player.Dir != 0 && player.Hp > 0 {
			ox := player.X
			oy := player.Y
			player.Move(self.mapId)
			self.view.PlayerMove(player, ox, oy, player.X, player.Y)
		}
		//玩家能量扣除（加速状态）和回复
		player.UpdateMp(self.current_frame)
	}
	//统一发给客户端
	self.executeCallBackMsg()
}

//每秒执行
func (self *Room) UpdateSec() {
	//1.安全区域范围处理
	self.doSafeArea()
	//2.在非安全区域玩家扣血
	self.doSafeArenaHurt()
}

/**
【每帧调用】统一发送所有消息给各个客户端
*/
func (self *Room) executeCallBackMsg() {
	for _, player := range self.players {
		if player.Loaded {
			//打包玩家发送数据
			self.view.Send(player, self.msgs)
			//清除ops数据
			player.Ops = nil
			//这里设置下一帧的开始坐标
			player.LX = player.X
			player.LY = player.Y
		}
	}
	//清空九宫格管理器的发送数据
	self.view.Clear()
	self.msgs = nil

	//删除过期子弹(判断规则是子弹已经射出30*3帧)
	for bulletId, b := range self.bullets {
		if self.current_frame-b.Frame >= 30*3 {
			delete(self.bullets, bulletId)
		}
	}
}

//初始化物品信息
func (self *Room) initItems(num int) {
	//随机n个不重复的坐标
	PosList := RandItemPosN(self.mapId, num)
	rand.Seed(time.Now().UnixNano())
	for _, pos := range PosList {
		//随机一件物品
		itemList := battleconf.BattleItemConf[self.mapId]
		index := rand.Intn(len(itemList))
		itemBaseId := itemList[index]
		self.AddItem(itemBaseId, 1, pos.X, pos.Y)
	}
}

//添加物品
func (self *Room) AddItem(itemBaseId, num int, x, y int) *Item {
	item := &Item{
		Id:     self.item_id,
		BaseId: itemBaseId,
		X:      x,
		Y:      y,
		Num:    num,
	}
	self.items[item.Id] = item
	self.item_id++
	return item
}

//随机武器
func (self *Room) initWeapon(num int) {
	//随机n个不重复的坐标
	PosList := RandWeaponPosN(self.mapId, num)
	rand.Seed(time.Now().UnixNano())
	for _, pos := range PosList {
		//随机一件物品
		weaponList := battleconf.BattleWeaponConf[self.mapId]
		index := rand.Intn(len(weaponList))
		weaponBaseId := weaponList[index]
		self.AddWeapon(weaponBaseId, pos.X, pos.Y)
	}
}

//新建一把武器
func (self *Room) AddWeapon(weaponBaseId, x, y int) *Weapon {
	BaseWeapon := GetBaseWeapon(weaponBaseId)
	weapon := &Weapon{
		Id:        self.weapon_id,
		BaseId:    weaponBaseId,
		X:         x,
		Y:         y,
		BulletNum: BaseWeapon.BulletMaxNum,
	}
	self.weapons[weapon.Id] = weapon
	self.weapon_id++
	return weapon
}

//根据玩家session获取玩家信息
func (self *Room) GetPlayer(session gate.Session) *Player {
	playerId, _ := strconv.Atoi(session.Get("PlayerId"))
	if objId, ok := self.playerid_to_objid[playerId]; ok {
		player, _ := self.players[objId]
		return player
	}
	return nil
}

//根据唯一id获取玩家信息
func (self *Room) GetPlayerById(objId int) *Player {
	if player, ok := self.players[objId]; ok {
		return player
	}
	return nil
}

//获取玩家唯一场景id
func (self *Room) GetPlayerBattleId(playerId int) int {
	if objId, ok := self.playerid_to_objid[playerId]; ok {
		return objId
	}
	return 0
}

//根据枪的唯一id获取枪的数据
func (self *Room) GetWeapon(weaponId int) *Weapon {
	if weapon, ok := self.weapons[weaponId]; ok {
		return weapon
	}
	return nil
}

//删除武器
func (self *Room) DelWeapon(weaponId int) {
	delete(self.weapons, weaponId)
}

//根据物品id获取物品信息
func (self *Room) GetItem(itemId int) *Item {
	if item, ok := self.items[itemId]; ok {
		return item
	}
	return nil
}

//删除物品
func (self *Room) DelItem(itemId int) {
	delete(self.items, itemId)
}

//获取子弹信息
func (self *Room) GetBullet(bulletId int) *Bullet {
	if bullet, ok := self.bullets[bulletId]; ok {
		return bullet
	}
	return nil
}
