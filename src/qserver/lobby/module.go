/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package lobby

import (
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/module/base"
	"sync"
)

var Module = func() module.Module {
	lobbys := new(Lobbys)
	return lobbys
}

type Lobbys struct {
	basemodule.BaseModule
	lobbys map[int]*Lobby
	mu      sync.Mutex
}

func (self *Lobbys) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "Lobby"
}
func (self *Lobbys) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (self *Lobbys) OnInit(app module.App, settings *conf.ModuleSettings) {
	self.BaseModule.OnInit(self, app, settings)
	self.lobbys = make(map[int]*Lobby)

	self.lobbys[0] = &Lobby{
		mapId: 0,
		readyAmount: 0,
		stoped: true,
	}
	self.lobbys[0].Init()

	self.GetServer().RegisterGO("c_lobby_enter", self.cLobbyEnter)
	self.GetServer().RegisterGO("c_lobby_match", self.cLobbyMatch)
	self.GetServer().RegisterGO("c_lobby_cancel", self.cLobbyCancel)
}

func (self *Lobbys) Run(closeSig chan bool) {
	for _, lobby := range self.lobbys {
		lobby.Start()
	}
}

func (self *Lobbys) OnDestroy() {
	//一定别忘了关闭RPC
	self.GetServer().OnDestroy()
}

// 获取地图id对应的房间
func (self *Lobbys) GetLobby(id int) (*Lobby, bool) {
	self.mu.Lock()
	defer self.mu.Unlock()

	lobby, ok := self.lobbys[id]
	return lobby, ok
}
