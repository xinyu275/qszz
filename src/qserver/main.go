package main

import (
"github.com/liangdas/mqant"
"github.com/liangdas/mqant/module/modules"
"qserver/gate"
)
func main() {
	app := mqant.CreateApp()
	//app.Route("Chat",ChatRoute)
	app.Run(true, //只有是在调试模式下才会在控制台打印日志, 非调试模式下只在日志文件中输出日志
		modules.MasterModule(),
		mgate.Module(),  //这是默认网关模块,是必须的支持 TCP,websocket,MQTT协议
		//tracing.Module(), //很多初学者不会改文件路径，先移除了
	)  //这是聊天模块

}

