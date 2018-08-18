package lobby

import (
	"github.com/liangdas/mqant/gate"
	"sync"
)

//玩家信息
type Player struct {
	Session         gate.Session
	PlayerId        int
	IsReady			bool

	mu sync.RWMutex
}

//更新玩家sesssion
func (p *Player) updateSession(session gate.Session) {
	p.mu.Lock()
	p.Session = session
	p.mu.Unlock()
}
