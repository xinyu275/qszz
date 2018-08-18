package battleconf

//枪的基础属性配置
type BaseWeapon struct {
	BaseId       int     // 枪baseid(就是现在得类型id)
	BulletMaxNum int     // 子弹最大数量
	Cd           float64 // 冷却cd（s）
	Weight       int     // 重量
	Speed        int     // 子弹速度
	Range        int     // 射程
	PerBulletNum int     // 单次射击发射子弹数(最多可以打中多少个人)
	Hurt         int     // 单发子弹伤害
}

var BaseWeaponList map[int]*BaseWeapon = map[int]*BaseWeapon{
	0: &BaseWeapon{0, 254, 0.3, 10, 800, 1000, 1, 10},  // 手枪
	1: &BaseWeapon{1, 20, 1.5, 20, 3000, 1000, 1, 80},  // 狙击枪
	2: &BaseWeapon{2, 50, 1, 30, 1200, 500, 5, 15},     // 散弹枪
	3: &BaseWeapon{3, 200, 0.01, 50, 1500, 800, 1, 10}, // 机关枪
	4: &BaseWeapon{4, 15, 2.5, 40, 1000, 1000, 1, 120}, // 电磁炮
}
