package user

import (
	"github.com/liangdas/mqant/gate"
	"sync"
)

//玩家信息
type Player struct {
	Session         gate.Session
	PlayerId        int
	Account         string
	Nickname        string
	Lv              int
	Exp             int
	RegTime         int
	LastLoginTime   int
	LastOfflineTime int
	LastLoginIp     string
	Gold            int

	mu sync.RWMutex
}

//UpdateMap: 根据map信息更新玩家数据
func (p *Player) updateByMap(m map[string]interface{}) {
}

//更新玩家sesssion
func (p *Player) updateSession(session gate.Session) {
	p.mu.Lock()
	p.Session = session
	p.mu.Unlock()
}
