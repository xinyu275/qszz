package battle

/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/

import (
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/module/base"
	"qserver/battle/battleconf"
)

var Module = func() module.Module {
	battle := new(Battle)
	return battle
}

type Battle struct {
	basemodule.BaseModule
	room *Room
}

func (b *Battle) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "Battle"
}
func (b *Battle) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (b *Battle) OnInit(app module.App, settings *conf.ModuleSettings) {
	b.BaseModule.OnInit(b, app, settings)
	//初始化map数据
	battleconf.Init()

	b.room = NewRoom(b, 2, 1)

	b.GetServer().Register("c_player_frame", b.cPlayerFrame)
	b.GetServer().Register("c_battle_loaded", b.cBattleLoaded)
	b.GetServer().Register("c_battle_start", b.cBattleStart)
	b.GetServer().Register("c_battle_hurt", b.cBattleHurt)
}

func (b *Battle) Run(closeSig chan bool) {
	b.room.Start()
}

func (b *Battle) OnDestroy() {
	//一定别忘了关闭RPCc
	b.GetServer().OnDestroy()
}
