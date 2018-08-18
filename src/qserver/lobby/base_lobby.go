package lobby

import (
	"fmt"
	"github.com/liangdas/mqant/log"
	"github.com/yireyun/go-queue"
	"reflect"
	"runtime"
	"sync"
)

type QueueMsg struct {
	Func   string
	Params []interface{}
}
type BaseLobby struct {
	functions       map[string]interface{}
	queue0          *queue.EsQueue
	queue1          *queue.EsQueue
	current_w_queue int
	mu              sync.RWMutex
}

func (self *BaseLobby) Init() {
	self.functions = map[string]interface{}{}
	self.queue0 = queue.NewQueue(256)
	self.queue1 = queue.NewQueue(256)
	self.current_w_queue = 0
}

//Register : 注册函数
func (self *BaseLobby) Register(id string, f interface{}) {
	if _, ok := self.functions[id]; ok {
		panic(fmt.Sprintf("function id %v: already registered", id))
	}

	self.functions[id] = f
}

/**
协成安全,任意协成可调用
*/
func (self *BaseLobby) PutQueue(_func string, params ...interface{}) error {
	ok, quantity := self.wqueue().Put(&QueueMsg{
		Func:   _func,
		Params: params,
	})
	if ok {
		return fmt.Errorf("Put Fail, quantity:%v\n", quantity)
	}
	return nil
}

/**
切换并且返回读的队列
*/
func (self *BaseLobby) switchqueue() *queue.EsQueue {
	self.mu.Lock()
	if self.current_w_queue == 0 {
		self.current_w_queue = 1
		self.mu.Unlock()
		return self.queue0
	}
	self.current_w_queue = 0
	self.mu.Unlock()
	return self.queue1
}

//获取当前队列
func (self *BaseLobby) wqueue() *queue.EsQueue {
	self.mu.Lock()
	if self.current_w_queue == 0 {
		self.mu.Unlock()
		return self.queue0
	} else {
		self.mu.Unlock()
		return self.queue1
	}
}

/**
【每帧调用】执行队列中的所有事件
*/
func (self *BaseLobby) ExecuteEvent() {
	ok := true
	queue := self.switchqueue()
	index := 0
	for ok {
		val, _ok, _ := queue.Get()
		index++
		if _ok {
			msg := val.(*QueueMsg)
			function, ok := self.functions[msg.Func]
			if !ok {
				log.Error("Remote function(%s) not found", msg.Func)
				continue
			}
			f := reflect.ValueOf(function)
			in := make([]reflect.Value, len(msg.Params))
			for k, _ := range in {
				in[k] = reflect.ValueOf(msg.Params[k])
			}
			_runFunc := func() {
				defer func() {
					if r := recover(); r != nil {
						var rn = ""
						switch r.(type) {

						case string:
							rn = r.(string)
						case error:
							rn = r.(error).Error()
						}
						buf := make([]byte, 1024)
						l := runtime.Stack(buf, false)
						errstr := string(buf[:l])
						log.Error("table qeueu event(%s) exec fail error:%s \n ----Stack----\n %s", msg.Func, rn, errstr)
					}
				}()
				f.Call(in)
			}
			_runFunc()
		}
		ok = _ok
	}
}
