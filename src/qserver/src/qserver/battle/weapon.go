package battle

//武器
type Weapon struct {
	Id     int // 武器唯一id
	BaseId int // 武器baseid
	X      int
	Y      int

	Owner     int // 所有者玩家战场唯一id
	BulletNum int //子弹数量
}

//武器子弹-1
func (weapon *Weapon) Fire() {
	weapon.BulletNum -= 1
}

//武器被拾取
func (weapon *Weapon) PickUp(player *Player) {
	weapon.Owner = player.Id
}

//丢弃
func (weapon *Weapon) Discard(player *Player) {
	weapon.Owner = 0
	weapon.X = player.X
	weapon.Y = player.Y
}

//增加子弹数量
func (weapon *Weapon) AddBullet(addNum int) {
	weapon.BulletNum += addNum
}
