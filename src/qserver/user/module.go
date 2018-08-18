/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package user

import (
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/module/base"
	"sync"
)

var Module = func() module.Module {
	user := new(User)
	return user
}

type User struct {
	basemodule.BaseModule
	players map[int]*Player
	mu      sync.Mutex
}

func (self *User) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "User"
}
func (self *User) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (self *User) OnInit(app module.App, settings *conf.ModuleSettings) {
	self.BaseModule.OnInit(self, app, settings)
	self.players = make(map[int]*Player)

	self.GetServer().RegisterGO("c_player_info", self.cPlayerInfo)
}

func (self *User) Run(closeSig chan bool) {
}

func (self *User) OnDestroy() {
	//一定别忘了关闭RPC
	self.GetServer().OnDestroy()
}

//根据玩家id获取玩家信息
func (self *User) getPlayer(playerId int) (*Player, error) {
	self.mu.Lock()
	defer self.mu.Unlock()

	if player, ok := self.players[playerId]; ok {
		return player, nil
	}
	player, err := getPlayerInfoFromDb(playerId)
	if err != nil {
		return nil, err
	}
	self.players[playerId] = player
	return player, nil
}

//TODO更新玩家session
