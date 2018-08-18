package battle

//战场物品信息
type Item struct {
	Id     int // 物品唯一id
	BaseId int // 物品baseid
	X      int
	Y      int

	Num   int //拾取会获得Num个baseid物品
	Owner int //物品所有者
}

//物品baseid修改(开宝箱)
func (self *Item) UpdateBaseId(newBaseId int) {
	self.BaseId = newBaseId
}
