package lobby

import (
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/module"
	"github.com/golang/protobuf/proto"
	"github.com/liangdas/mqant/log"
	"math/rand"
	"time"
	"sync"
	
	"qserver/mproto"
)


const Default_WaitMatchingTime time.Duration = 10
const Default_PlayerInRoomMax int = 50

type Lobby struct {
	BaseLobby	
	module        	module.Module
	roomId        	int // 房间id, TODO 创建房间
	mapId         	int

	// 角色
	players       	map[int]*Player //房间 PlayerID - >PlayerInfo
	readyAmount		int // 当前排队的玩家数量

	stoped			bool
	mu      		sync.Mutex
}

//新建房间
func (self *Lobby) NewRoom() ([]byte, error) {
	// 开启一个新的房间
	// TODO 增加房间ID
	code := mproto.ECode_ok
	pb := &mproto.SLobbyMatch{
		Code:     &code,
	}
	b, err := proto.Marshal(pb)
	if err != nil {
		return b, err
	}

	// 将玩家拉入新的房间
	self.mu.Lock()
	defer self.mu.Unlock()
	playerInRoom := 0
	for _, player := range self.players {
		playerInRoom++
		go player.Session.Send("s_lobby_match", b)

		// 移除进入房间的玩家
		player.IsReady = false
		self.readyAmount--
		delete(self.players, player.PlayerId)

		if playerInRoom >= Default_PlayerInRoomMax {
			break
		}
	}

	// 返回房间序列化后的内容
	return b, err
}

//Start: 大厅开启
func (self *Lobby) Start() {
	if self.stoped {
		self.stoped = false
		go func() {
			rand.Seed(time.Now().UnixNano())
			// 每10秒匹配一次房间
			tick := time.NewTicker(Default_WaitMatchingTime * time.Second)
			defer func() {
				tick.Stop()
			}()
			for !self.stoped {
				select {
				case <-tick.C:
					self.Update()
				}
			}
		}()
	}
}

func (self *Lobby) Stop() {
	self.stoped = true
}

// 每间隔Default_WaitMatchingTime秒，开始一次匹配的检测
func (self *Lobby) Update() {
	//处理队列
	//self.ExecuteEvent()

	// 判断排队的玩家数量，是否足够开启房间
	if self.readyAmount > 0 {
		// 开启一个新的房间
		for {
			// 剩余排队的人数大于一个房间的上限，继续开启新的房间
			if (self.readyAmount >= Default_PlayerInRoomMax) {
				// 开启一个新的房间
				_, err := self.NewRoom()
				if err != nil {
					log.Error("[Error] NewRoom is faild.")
				}
			} else {
				break
			}
		}
	}
}

//添加玩家
func (self *Lobby) AddPlayer(playerId int, session gate.Session) bool {
	self.mu.Lock()
	defer self.mu.Unlock()

	_, ok := self.players[playerId]
	if ok == false {
		self.players[playerId] = &Player{
			Session: session,
			PlayerId: playerId,
			IsReady: false,
		}
	}

	return true
}

//移除玩家
func (self *Lobby) DelPlayer(playerId int) bool {
	self.mu.Lock()
	defer self.mu.Unlock()

	delete(self.players, playerId)
	
	return true
}

//排队的玩家
func (self *Lobby) Match(playerId int) bool {
	self.mu.Lock()
	defer self.mu.Unlock()

	player, ok := self.players[playerId]
	if ok && player.IsReady == false {
		player.IsReady = true
		self.readyAmount++
		return true
	}

	return false
}

//取消排队的玩家
func (self *Lobby) UnMatch(playerId int) bool {
	self.mu.Lock()
	defer self.mu.Unlock()

	player, ok := self.players[playerId]
	if ok && player.IsReady == true {
		player.IsReady = false
		self.readyAmount--
		return true
	}
	
	return false
}
