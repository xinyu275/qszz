package battle

import "qserver/battle/battleconf"

// 使用技能
func (self *Room) castSkill(player *Player, skillId int) {
	if !self.checkSkill(player, skillId) {
		return
	}
	//技能处理
	baseSkill := GetBaseSkill(skillId)
	if baseSkill == nil {
		return
	}
	//能量值不够
	if player.Mp < baseSkill.CostMp {
		return
	}
	switch skillId {
	case SKILL_UP_SPEED:
		self.skillSpeedUp(player, baseSkill)
	}
}

//技能判断
func (self *Room) checkSkill(player *Player, skillId int) bool {
	if skillId == 0 && player.MoveSpeedUp > 0 {
		self.skillStopSpeedUp(player)
		return false
	}
	//已经在加速中了
	if skillId == SKILL_UP_SPEED && player.MoveSpeedUp > 0 {
		return false
	}
	//如果释放其他技能还在加速，这里要停止加速(作弊可以出现这样的情况)
	if skillId > 0 && player.MoveSpeedUp > 0 {
		self.skillStopSpeedUp(player)
	}

	if skillId == 0 {
		return false
	}
	return true
}

//1.加速技能
func (self *Room) skillSpeedUp(player *Player, baseSkill *battleconf.BaseSkill) {
	//这里不扣能量，系统计时器会根据玩家是否在加速状态 扣或者回复能量
	player.SpeedUp(self.current_frame)
}

//2.停止加速
func (self *Room) skillStopSpeedUp(player *Player) {
	player.StopSpeedUp(self.current_frame)
}
