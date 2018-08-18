package battle

//子弹
type Bullet struct {
	Id           int // 子弹唯一id
	Owner        int // 发射者玩家战场虚拟唯一id
	WeaponId     int // 发射者枪的id
	WeaponBaseId int // 发射者枪的baseid
	Frame        int // 发射的时候的帧数
	FireDir      int // 发送方向
	StartX       int // 子弹开始位置
	StartY       int
	AddRange     int // 子弹射程增益
	AddSpeed     int // 子弹速度增益
	AddHurt      int // 子弹伤害增益
	HurtObjNum   int // 子弹射中的人数，散弹可以是同一个人受伤害
}

//新建子弹
func NewBullet(room *Room, player *Player, fireDir int) *Bullet {
	return &Bullet{
		Id:           room.bullet_id,
		Owner:        player.Id,
		WeaponId:     player.CurWeaponId,
		WeaponBaseId: player.WeaponType,
		Frame:        room.current_frame,
		FireDir:      fireDir,
		StartX:       player.X,
		StartY:       player.Y,
		AddRange:     player.AddRange,
		AddSpeed:     player.AddBulletSpeed,
		AddHurt:      player.AddHurt,
	}
}

//子弹射中
func (self *Bullet) AddHurtObjNum() {
	self.HurtObjNum++
}
