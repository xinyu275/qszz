package battleconf

//战斗中的配置
type PosXY struct {
	X int
	Y int
}

//1.道具位置配置
//mapid -> []int
var BattleItemPosConf map[int][]PosXY = map[int][]PosXY{
	2: []PosXY{
		{66 * 64, 90 * 64},
		{68 * 64, 92 * 64},
		{60 * 64, 88 * 64},
		{63 * 64, 86 * 64},
	},
}

//道具列表
var BattleItemConf map[int][]int = map[int][]int{
	2: []int{0, 1, 3, 10, 10, 10},
}

//武器坐标
var BattleWeaponPosConf map[int][]PosXY = map[int][]PosXY{
	2: []PosXY{
		{66 * 64, 86 * 64},
		{68 * 64, 80 * 64},
		{60 * 64, 82 * 64},
		{63 * 64, 84 * 64},
	},
}

//武器id列表
var BattleWeaponConf map[int][]int = map[int][]int{
	2: []int{0, 3, 1, 2, 4},
}

//角色坐标
var BattlePlayerPosConf map[int][]PosXY = map[int][]PosXY{
	2: []PosXY{
		{59 * 64, 85 * 64},
		{62 * 64, 85 * 64},
	},
}
