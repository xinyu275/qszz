package battle

//九宫格里面的一个单元格
type Grid struct {
	//当前九宫格玩家信息
	inPlayers []*Player
	//格子管理器
	view *View

	//信息
	msgs []*Msg
}

//新建格子
func NewGrid(view *View, grid int) *Grid {
	return &Grid{
		inPlayers: make([]*Player, 0),
		msgs:      make([]*Msg, 0),
		view:      view,
	}
}

//添加一个玩家
func (g *Grid) addPlayer(player *Player) {
	//如果已经存在，就不要添加了
	exists := false
	for _, p := range g.inPlayers {
		if p.Id == player.Id {
			exists = true
			break
		}
	}
	if !exists {
		g.inPlayers = append(g.inPlayers, player)
	}
	//添加移动的消息（即使不移动）
	g.packMove(player)
	//玩家属性包
	g.packPlayerExtraAttr(player)
}

//玩家离开本格
func (g *Grid) playerLeave(player *Player) {
	//删除玩家
	for index, p := range g.inPlayers {
		if p.Id == player.Id {
			g.inPlayers = append(g.inPlayers[:index], g.inPlayers[index+1:]...)
			break
		}
	}
	//添加移动消息（如果离开玩家视野，那前段要删除对象）
	g.packMove(player)
}

//玩家移动
func (g *Grid) playerMove(player *Player) {
	g.packMove(player)
}

//玩家停止移动
func (g *Grid) playerStopMove(player *Player) {
	//只需要打包移动包(因为移动包里有玩家属性dir)
	g.packMove(player)
}

//开火后九宫格封装
func (g *Grid) fire(bullet *Bullet) {
	//打包子弹信息报，九宫格广播
	msg := packBullet(g.view.GetMsgId(), bullet)
	g.msgs = append(g.msgs, msg)
}

//玩家换枪
func (g *Grid) switchWeapon(player *Player) {
	//换枪 只要打包玩家当前属性就行了
	g.packPlayerExtraAttr(player)
}

//玩家额外属性改变
func (g *Grid) changePlayerExtraAttr(player *Player) {
	g.packPlayerExtraAttr(player)
}

//清空发送数据
func (g *Grid) clear() {
	g.msgs = nil
}

//获取opsmsg
func (g *Grid) getMsgs() []*Msg {
	return g.msgs
}

//玩家移动
func (g *Grid) packMove(player *Player) {
	Msg := packMove(g.view.GetMsgId(), player)
	g.msgs = append(g.msgs, Msg)
}

//玩家属性,玩家血量/武器/穿着/等
func (g *Grid) packPlayerExtraAttr(player *Player) {
	Msg := packUnitExtra(g.view.GetMsgId(), player)
	g.msgs = append(g.msgs, Msg)
}
