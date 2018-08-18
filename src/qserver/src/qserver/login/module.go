/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package login

import (
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/module/base"
)

var Module = func() module.Module {
	gate := new(Login)
	return gate
}

type Login struct {
	basemodule.BaseModule
}

func (m *Login) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "Login"
}
func (m *Login) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (m *Login) OnInit(app module.App, settings *conf.ModuleSettings) {
	m.BaseModule.OnInit(m, app, settings)
	m.GetServer().RegisterGO("c_player_login", m.cPlayerLogin)
	m.GetServer().RegisterGO("c_server_time", m.cServerTime)
}

func (m *Login) Run(closeSig chan bool) {
}

func (m *Login) OnDestroy() {
	//一定别忘了关闭RPCc
	m.GetServer().OnDestroy()
}
