package battleconf

//物品基础属性配置
type BaseItem struct {
	BaseId  int // 物品baseid(就是现在得类型id)
	Overlap int // 物品叠加上限
	Min     int //每种物品意义不一样（急救箱：加血区间[min, max]）
	Max     int
	Weight  int //重量
}

var BaseItemList map[int]*BaseItem = map[int]*BaseItem{
	0:  &BaseItem{0, 0, 50, 50, 0},     // 血包(急救箱)
	1:  &BaseItem{1, 0, 1000, 1000, 0}, // 弹夹包
	2:  &BaseItem{2, 15, 16, 16, 1},    // 防御
	3:  &BaseItem{3, 0, 200, 200, 0},   // 护盾
	4:  &BaseItem{4, 12, 20, 20, 1},    // 武器的射程
	5:  &BaseItem{5, 10, 3, 3, 1},      // 武器的伤害
	6:  &BaseItem{6, 12, 30, 30, 1},    // 武器的弹速
	7:  &BaseItem{7, 12, 400, 400, 1},  // 武器的冷却（万分比）
	8:  &BaseItem{8, 15, 20, 20, 1},    // 增加能量值
	9:  &BaseItem{9, 15, 30, 30, 1},    // 增加血量上限
	10: &BaseItem{10, 0, 0, 0, 0},      // 随机包
}
