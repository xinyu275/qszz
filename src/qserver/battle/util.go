package battle

import (
	"github.com/liangdas/mqant/log"
	"math/rand"
	"qserver/battle/battleconf"
	"strconv"
	"time"
)

//判断x,y坐标是否可走
func CanWalk(mapId, x, y int) bool {
	if x <= 0 || y <= 0 {
		return false
	}
	cellWidthX := x / battleconf.CellToPx[mapId]
	cellWidthY := y / battleconf.CellToPy[mapId]
	mmap := battleconf.MapWalk[mapId]
	_, ok := mmap[cellWidthX*10000+cellWidthY]
	return ok
}

//获取玩家的随机位置
func RandPlayerPos(mapId int) (int, int) {
	posList := battleconf.BattlePlayerPosConf[mapId]
	length := len(posList)
	rand := rand.Intn(length)
	pos := posList[rand]
	return pos.X, pos.Y
}

//随机N个坐标 物品坐标
func RandItemPosN(mapId, n int) []battleconf.PosXY {
	posList := battleconf.BattleItemPosConf[mapId]
	return RandPosN(posList, n)
}

//随机n个坐标 武器
func RandWeaponPosN(mapId, n int) []battleconf.PosXY {
	posList := battleconf.BattleWeaponPosConf[mapId]
	return RandPosN(posList, n)
}

//在坐标列表中随机N个
func RandPosN(posList []battleconf.PosXY, n int) []battleconf.PosXY {
	if len(posList) <= n {
		return posList
	}
	rand.Seed(time.Now().UnixNano())
	length := len(posList)
	r := make([]int, 0, n)
	p := make([]int, length)
	for i := 0; i < length; i++ {
		p[i] = i
	}
	for len(r) < n {
		lenp := len(p)
		index := rand.Intn(lenp)
		r = append(r, p[index])
		p = append(p[:index], p[index+1:]...)
	}
	result := make([]battleconf.PosXY, n)
	for index, v := range r {
		result[index] = posList[v]
	}
	return result
}

//根据物品baseid获取基础物品属性
func GetBaseItem(baseId int) *battleconf.BaseItem {
	if baseItem, ok := battleconf.BaseItemList[baseId]; ok {
		return baseItem
	}
	log.Error("item:" + strconv.Itoa(baseId) + " not config")
	return nil
}

//根据枪baseid(类型)获取基础属性
func GetBaseWeapon(baseId int) *battleconf.BaseWeapon {
	if baseWeapon, ok := battleconf.BaseWeaponList[baseId]; ok {
		return baseWeapon
	}
	log.Error("weapon:" + strconv.Itoa(baseId) + " not config")
	return nil
}

//根据技能id获取技能基础属性
func GetBaseSkill(skillId int) *battleconf.BaseSkill {
	if baseSkill, ok := battleconf.BaseSkillList[skillId]; ok {
		return baseSkill
	}
	log.Error("skillId:" + strconv.Itoa(skillId) + " not config")
	return nil
}

//根据safeindex获取安全区域信息
func GetBaseSafeArea(safeIndex int) *battleconf.BaseSafeArea {
	if baseSafeArea, ok := battleconf.BaseSafeAreaList[safeIndex]; ok {
		return baseSafeArea
	}
	return nil
}

//随机[min, max]
func Rand(min, max int) int {
	if min >= max {
		return max
	}
	return rand.Intn(max-min+1) + min
}
