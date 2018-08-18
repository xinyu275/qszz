package battle

import (
	"math"
	"qserver/common"
)

//安全区域逻辑处理
func (self *Room) doSafeArea() {
	if self.safe_end {
		return
	}
	nowTime := common.GetUnixTime()
	if self.safe_index == 0 && nowTime >= self.start_time {
		self.doSafeNext(nowTime)
		return
	} else if nowTime < self.start_time {
		// 还没开始
		return
	}
	//safe_index > 0,那就要处理区域大小了
	self.doResizeSafeArea(nowTime)

}

//开启下一个
func (self *Room) doSafeNext(nowTime int) {
	nextSafeIndex := self.safe_index + 1
	baseSafe := GetBaseSafeArea(nextSafeIndex)
	//没有下一个索引了
	if baseSafe == nil {
		self.safe_end = true
		return
	}
	self.safe_index++
	self.safe_index_start_time = nowTime
}

//调整安全区域
func (self *Room) doResizeSafeArea(nowTime int) {
	baseSafe := GetBaseSafeArea(self.safe_index)
	//如果时间 nowTime >= startTime+ContinueTime,那就要马上调整为目标安全区域大小
	//先计算这一档次最后的目标大小
	objWidth := int(math.Sqrt(baseSafe.ObjAreaSize) * float64(MapWidth))
	objHeight := int(math.Sqrt(baseSafe.ObjAreaSize) * float64(MapHeight))
	objX := MapWidth/2 - objWidth/2
	objY := MapHeight/2 - objHeight/2
	//按时间算比例
	if nowTime >= self.safe_index_start_time+baseSafe.ContinueTime {
		//调整为目标安全区域大小
		self.safe_left_x = objX
		self.safe_left_y = objY
		self.safe_right_x = objX + objWidth
		self.safe_right_y = objY + objHeight
	} else {
		// 按比例调整安全区域
		rate := (nowTime - self.safe_index_start_time) / baseSafe.ContinueTime
		self.safe_left_x += (objX - self.safe_left_x) * rate
		self.safe_left_y += (objY - self.safe_left_y) * rate
		self.safe_right_x -= (self.safe_right_x - (objX + objWidth)) * rate
		self.safe_right_y -= (self.safe_right_y - (objY + objHeight)) * rate
	}

	if nowTime >= self.safe_index_start_time+baseSafe.ContinueTime+baseSafe.WaitTime {
		//开启下一个
		self.doSafeNext(nowTime)
	}
}

//非安全区域玩家扣血处理
func (self *Room) doSafeArenaHurt() {
	if self.safe_index <= 0 {
		return
	}
	baseSafe := GetBaseSafeArea(self.safe_index)
	for _, player := range self.players {
		if player.Hp <= 0 {
			continue
		}
		if self.checkInSafeArea(player) {
			continue
		}
		//扣除玩家的血量
		hurt := int(baseSafe.Val * float64(player.HpMax))
		self.doPlayerHurt(nil, player, hurt)
	}
}

//判断玩家是否在安全区域
//bool(true:在安全区域)
func (self *Room) checkInSafeArea(player *Player) bool {
	return player.X >= self.safe_left_x && player.X <= self.safe_right_x && player.Y >= self.safe_left_y && player.Y <= self.safe_right_y
}
