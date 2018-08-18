package battle

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"github.com/liangdas/mqant/log"
	"qserver/mproto"
	"sort"
)

//视野 管理九宫格
type View struct {
	GridList map[int]*Grid
	MsgId    int //消息唯一id
}

func NewView() *View {
	return &View{
		GridList: make(map[int]*Grid),
		MsgId:    1,
	}
}

//GetGrid 获取格子的*Grid信息
func (self *View) GetGrid(gridId int) *Grid {
	if g, ok := self.GridList[gridId]; ok {
		return g
	}
	g := NewGrid(self, gridId)
	self.GridList[gridId] = g
	return g
}

//初始化玩家
func (self *View) AddPlayer(player *Player) {
	gridId := CalcGridByXY(player.X, player.Y)
	g := self.GetGrid(gridId)
	g.addPlayer(player)
}

//玩家移动
func (self *View) PlayerMove(player *Player, ox, oy, nx, ny int) {
	oGrid := CalcGridByXY(ox, oy)
	nGrid := CalcGridByXY(nx, ny)
	if oGrid == nGrid {
		//只发移动包
		g := self.GetGrid(nGrid)
		g.playerMove(player)
		return
	}
	og := self.GetGrid(oGrid)
	ng := self.GetGrid(nGrid)
	//这里会导致如果离开的格子和当前格子都在一个玩家视野，那会发2个相同的包，不过对游戏没影响
	og.playerLeave(player)
	ng.addPlayer(player)
}

//玩家停止移动
func (self *View) PlayerStopMove(player *Player) {
	gridId := CalcGridByXY(player.X, player.Y)
	g := self.GetGrid(gridId)
	g.playerStopMove(player)
}

//玩家发射子弹
func (self *View) Fire(bullet *Bullet) {
	gridId := CalcGridByXY(bullet.StartX, bullet.StartY)
	g := self.GetGrid(gridId)
	g.fire(bullet)
}

//玩家换枪
func (self *View) SwitchWeapon(player *Player) {
	gridId := CalcGridByXY(player.X, player.Y)
	g := self.GetGrid(gridId)
	g.switchWeapon(player)
}

//玩家额外属性改变
func (self *View) ChangePlayerExtraAttr(player *Player) {
	gridId := CalcGridByXY(player.X, player.Y)
	g := self.GetGrid(gridId)
	g.changePlayerExtraAttr(player)
}

//清空所有的发送数据
func (self *View) Clear() {
	for _, g := range self.GridList {
		g.clear()
	}
}

//消息id自增
func (self *View) GetMsgId() int {
	MsgId := self.MsgId
	self.MsgId += 1
	return MsgId
}

//打包发送玩家视野数据
//TODO 这里可以优化，在同一个九宫格小格子中，打包数据是一样的
func (self *View) Send(player *Player, commonMsgs []*Msg) {
	msgs := make([]*Msg, 0, 16)
	//1.根据lx,ly,x,y判断跨格，吧新玩家进视野（不管玩家是否走动）
	addGridList := CalcAddGrid(player.LX, player.LY, player.X, player.Y)
	if len(addGridList) > 0 {
		for _, gridId := range addGridList {
			if g, ok := self.GridList[gridId]; ok {
				for _, player := range g.inPlayers {
					//移动包
					msg := packMove(0, player)
					msgs = append(msgs, msg)
					//属性包
					extraUnitMsg := packUnitExtra(0, player)
					msgs = append(msgs, extraUnitMsg)
				}
			}
		}
	}
	//2.吧玩家视野里的九宫格数据排序累加
	gridList := CalcGridList(player.X, player.Y)
	for _, gridId := range gridList {
		if g, ok := self.GridList[gridId]; ok {
			msgs = append(msgs, g.getMsgs()...)
		}
	}
	//玩家自己私有的操作
	if len(player.Ops) > 0 {
		msgs = append(msgs, player.Ops...)
	}
	//共有消息
	if len(commonMsgs) > 0 {
		msgs = append(msgs, commonMsgs...)
	}
	if len(msgs) > 0 {
		//msgs排序
		sort.Sort(SortMsg(msgs))
		reply := &mproto.SPlayerFrame{}
		//封装数据
		dataUnits := make([][]byte, 0)
		dataExtras := make([][]byte, 0)
		dataItems := make([][]byte, 0)
		dataWeapons := make([][]byte, 0)
		dataBullets := make([][]byte, 0)
		dataOps := make([][]byte, 0)
		for _, msg := range msgs {
			switch msg.mtype {
			case 0:
				dataUnits = append(dataUnits, msg.body)
			case 1:
				dataExtras = append(dataExtras, msg.body)
			case 2:
				dataItems = append(dataItems, msg.body)
			case 3:
				dataWeapons = append(dataWeapons, msg.body)
			case 4:
				dataBullets = append(dataBullets, msg.body)
			case 5:
				dataOps = append(dataOps, msg.body)
			}
		}
		if len(dataUnits) > 0 {
			reply.Units = bytes.Join(dataUnits, []byte(""))
		}
		if len(dataExtras) > 0 {
			reply.ExtraUnits = bytes.Join(dataExtras, []byte(""))
		}
		if len(dataItems) > 0 {
			reply.Items = bytes.Join(dataItems, []byte(""))
		}
		if len(dataWeapons) > 0 {
			reply.Weapons = bytes.Join(dataWeapons, []byte(""))
		}
		if len(dataBullets) > 0 {
			reply.Bullets = bytes.Join(dataBullets, []byte(""))
		}
		if len(dataOps) > 0 {
			reply.Ops = bytes.Join(dataOps, []byte(""))
		}
		b, err := proto.Marshal(reply)
		if err != nil {
			log.Error("send fail:%s", err.Error())
			return
		}
		//log.Info("s_player_frame success: bytes:%d", (len(b)))
		player.Session.Send("s_player_frame", b)
	}
}

//根据像素计算Grid
func CalcGridByXY(x, y int) int {
	return GenGridId(x/GridWidth, y/GridHeight)
}

//根据GrixX,GridY生成GridId
func GenGridId(gridX int, gridY int) int {
	return gridX*10000 + gridY
}

//排序
type SortMsg []*Msg

func (s SortMsg) Len() int { return len(s) }
func (s SortMsg) Less(i, j int) bool {
	return s[i].msgId < s[j].msgId
}
func (s SortMsg) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

//计算新加的九宫格视野
func CalcAddGrid(ox, oy, nx, ny int) []int {
	result := make([]int, 0)
	m := make(map[int]int)
	if CalcGridByXY(ox, oy) == CalcGridByXY(nx, ny) {
		return nil
	}
	oGridList := CalcGridList(ox, oy)
	for _, v := range oGridList {
		m[v] = 1
	}
	nGridList := CalcGridList(nx, ny)
	for _, gridId := range nGridList {
		if _, ok := m[gridId]; !ok {
			result = append(result, gridId)
		}
	}
	return result
}

//获取九宫格的gridid列表
func CalcGridList(x, y int) []int {
	r := make([]int, 0)
	gx := x / GridWidth
	gy := y / GridHeight
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			r = append(r, GenGridId(gx+i, gy+j))
		}
	}
	return r
}
