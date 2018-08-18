package battleconf

//技能基础属性配置
type BaseSkill struct {
	SkillId int     // 技能id(参见define.go中定义)
	CostMp  int     // 消耗能量（加速是每秒消耗能量值，其他技能是单次消耗能量值）
	Cd      float64 // cd冷却时间(s)
	Val     int     // 每个技能属性(加速是加速万分比(10000就是基础值，>10000才会加速))
}

var BaseSkillList map[int]*BaseSkill = map[int]*BaseSkill{
	1: &BaseSkill{1, 20, 0, 20000}, // 加速
}
