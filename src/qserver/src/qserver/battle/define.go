package battle

//定义战斗的常量

// 按钮的功能类型
//开火
var BtnFire uint32 = 1 << 6

//换枪
var BtnSwitchWeapon uint32 = 1 << 7

//拾取
var BtnPickUp uint32 = 1 << 8

//丢弃
var BtnDrop uint32 = 1 << 9

//开启倍镜
var BtnXTimesMirror uint32 = 1 << 10

//使用道具
var BtnUseItem uint32 = 1 << 11

//释放技能
var BtnCast uint32 = 1 << 12

//驾驶车、船、摩托
var BtnDrive uint32 = 1 << 13

//销毁
var BtnDestroy uint32 = 1 << 14

//只有移动，没有其它操作
var BtnDir uint32 = 1 << 15

var (
	//九宫格一格占多少像素
	GridWidth  = 800
	GridHeight = 400

	//地图像素大小
	MapWidth  = 12800
	MapHeight = 12800
)

//帧时间设置
var (
	//每帧毫秒数
	PER_FRAME_MILLISECOND = 30
	//1s能量回复换算帧数
	MP_CHANGE_INTERVAL_FRAME = 1000/PER_FRAME_MILLISECOND + 1
)

var (
	//真实玩家
	ROLE_PLATER = 1
	//机器人
	ROLE_ROBOT = 2
)

//物品类型
var (
	OBJ_UNIT   = 0 // 玩家/机器人
	OBJ_WEAPON = 1 // 武器
	OBJ_ITEM   = 2 //物品
	OBJ_BULLET = 3 //子弹
)

//武器类型
var (
	WEAPON_HANDGUN      = 0 // 手枪
	WEAPON_SNIPER_RIFLE = 1 // 狙击枪
	WEAPON_SHOTGUN      = 2 // 散弹枪
	WEAPON_MACHINE_GUN  = 3 // 机关枪
	WEAPON_ELECTRIC_GUN = 4 // 电磁炮
)

//默认武器
var (
	DEFAULT_WEAPON = WEAPON_HANDGUN //默认武器是手枪
)

//返回操作类型
var (
	OP_FIRE         = 1 // 开枪
	OP_SWITCHWEAPON = 2 // 换枪
	OP_PICKUP       = 3 // 拾取
	OP_DROP         = 4 // 丢弃(不需要)
	OP_DESTROY      = 5 // 销毁(暂时不需要)
	OP_USEITEM      = 6 // 使用道具(不需要)
	OP_CAST         = 7 // 释放技能
	OP_DIR          = 8 // 只有移动，没有其它操作
	OP_MP           = 9 // 能量值上限改变
)

//子弹最大id
var (
	//子弹id最多占14位
	MAX_BULLET_ID = 1<<14 - 1
)

//道具类型
var (
	ITEM_BOLLD          = 0  // 血包(急救箱)
	ITEM_BULLET         = 1  // 弹夹包
	ITEM_DEFENSEBUFFER  = 2  // 防御
	ITEM_SHIELDBUFFER   = 3  // 护盾
	ITEM_BULLETDISTANCE = 4  // 武器的射程
	ITEM_BULLETHURT     = 5  // 武器的伤害
	ITEM_BULLETSPEED    = 6  // 武器的弹速
	ITEM_WEAPONCD       = 7  // 武器的冷却
	ITEM_ADD_ENERGY     = 8  // 增加能量值
	ITEM_UPMAXHPBUFFER  = 9  // 增加血量上限buff
	ITEM_RANDOM         = 10 // 随机包
)

var (
	//角色属性
	DEFAULT_HP     = 100 // 默认血量
	DEFAULT_WEIGHT = 50  // 默认重量
	DEFAULT_SPEED  = 200 // 默认移动速度
	DEFAULT_MP     = 100 // 默认能量值
	DEFAULT_DEF    = 0   // 默认防御
	RECOVER_MP     = 20  // 每秒回复能量
	//角色状态
	//ROLE_STATE_IDLE     = 0 // 空闲状态
	//ROLE_STATE_FIGHTING = 1 // 作战状态
	//ROLE_STATE_HIDING   = 2 // 躲藏（隐身）
	//ROLE_STATE_DIE      = 3 // 死亡
)

//玩家死亡增益型buff散落最大数量，程序自动判断堆叠掉落
var (
	DROP_MAX_PLUSITEM = 30
)

//技能
var (
	SKILL_UP_SPEED = 1 // 加速
)
