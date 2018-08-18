package battle

import (
	"github.com/liangdas/mqant/log"
	"math"
	"qserver/battle/battleconf"
)

//如果是机器人嵌入到player中
type Robot struct {
	robot_state   int // 机器人状态
	robot_att_obj int // 上次攻击目标
	//上1s的坐标,如果在dir>0，1s坐标还没改变，那说明没可走的点了，要改变方向了
	last_x int
	last_y int
}

var (
	ROBOT_STATE_IDLE      = 0
	ROBOT_STATE_MOVE_DRAG = 1 // 移出毒圈
	ROBOT_STATE_FIGHTING  = 2 // 攻击/移动攻击中
)

//机器人逻辑
func (self *Room) doRobot() {
	for _, player := range self.players {
		//if player.RoleType == ROLE_ROBOT && player.Hp > 0 {
		//	self.doRobotPlayer(player)
		//}
		self.doRobotPlayer(player)
	}
}

//机器人处理
func (self *Room) doRobotPlayer(player *Player) {
	//1.判断是否在毒圈中，如果在，那就往中心移动,如果不是往毒圈外移动，那就要改变方向
	//2.如果玩家是攻击状态（前1s在攻击玩家），那判断目标是否还在自己的视野范围内
	//如果在并且在攻击范围内，那就攻击，如果不在攻击范围，那就朝着目标移动再攻击
	//3.如果不是攻击状态，那看视野里面有没有玩家，优先选择真实玩家pk
	//4.如果视野里没玩家那就随机移动

	//1.毒圈判断
	if !self.checkInSafeArea(player) {
		self.robotMoveOutDrag(player)
		return
	}

	//2.攻击状态下的判断
	if player.robot_state == ROBOT_STATE_FIGHTING && player.robot_att_obj > 0 {
		objPlayer := self.GetPlayerById(player.robot_att_obj)
		if objPlayer != nil && player.Hp > 0 {
			if self.robotBattle(player, objPlayer) {
				return
			}
		}
	}

	//3.找视野里的玩家
	var viewPlayers []*Player = self.getViewPlayers(player)
	for _, objPlayer := range viewPlayers {
		if objPlayer.Id != player.Id && objPlayer.Hp > 0 {
			if self.robotBattle(player, objPlayer) {
				return
			}
		}
	}

	//4.随机移动吧
	self.robotRandMove(player)
	player.robot_state = ROBOT_STATE_IDLE
	player.last_x = player.X
	player.last_y = player.Y
}

//在毒圈中，往中心移动
func (self *Room) robotMoveOutDrag(player *Player) {
	//目标在移出毒圈中
	if player.robot_state == ROBOT_STATE_MOVE_DRAG {
		return
	}
	//计算中心点角度
	centerX := MapWidth / 2
	centerY := MapHeight / 2
	self.MoveToTarget(player, centerX, centerY)
	player.robot_state = ROBOT_STATE_MOVE_DRAG
}

//如果机器人在攻击状态，如果目标在，那就攻击
func (self *Room) robotBattle(player, objPlayer *Player) bool {
	//如果不在视野范围，那就返回false
	if math.Abs(float64(player.X)-float64(objPlayer.X)) >= float64(GridWidth)*1.2 ||
		math.Abs(float64(player.Y)-float64(objPlayer.Y)) >= float64(GridHeight)*1.2 {
		return false
	}
	//TODO 判断是否在攻击范围
	checkAttArea := true
	if !checkAttArea {
		//移动取攻击
		self.MoveToTarget(player, objPlayer.X, objPlayer.Y)
	} else {
		//如果在攻击范围内，那就停止移动，并且攻击
		player.Dir = 0
		log.Info("发起攻击.....")
		//TODO 攻击
	}
	player.robot_state = ROBOT_STATE_FIGHTING
	player.robot_att_obj = objPlayer.Id
	//发射子弹吧
	return true
}

//机器人随机移动
func (self *Room) robotRandMove(player *Player) {
	if player.Dir > 0 && player.robot_state == ROBOT_STATE_IDLE &&
		player.last_x == player.X && player.last_y == player.Y {
		player.Dir += 45
		if player.Dir >= 360 {
			player.Dir = 10
		}
		return
	}
	if player.Dir == 0 {
		player.Dir = Rand(120, 240)
	}
}

//找玩家视野内的玩家，这里简化一下，不判断九宫格了，判断5个自己+上下左右
func (self *Room) getViewPlayers(player *Player) []*Player {
	var gridxy [5]battleconf.PosXY = [5]battleconf.PosXY{
		{player.X, player.Y},
		{player.X - GridWidth, player.Y},
		{player.X + GridWidth, player.Y},
		{player.X, player.Y - GridHeight},
		{player.X, player.Y + GridHeight},
	}
	var viewPlayers []*Player = make([]*Player, 0, 5)
	for _, pos := range gridxy {
		viewPlayers = append(viewPlayers, self.view.GetGridPlayers(pos.X, pos.Y)...)
	}
	return viewPlayers
}

//玩家移向目标
func (self *Room) MoveToTarget(player *Player, objx, objy int) {
	dx := float64(objx - player.X)
	dy := float64(objy - player.Y)
	dir := 180 - int(math.Atan2(dy, dx)/math.Pi*180)
	player.Dir = dir + 1
	if player.Dir >= 360 {
		player.Dir = 360
	}
}
