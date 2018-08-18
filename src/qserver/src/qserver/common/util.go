package common

import "time"

//公用函数
//1.获取当前时间抽（毫秒）
func GetMillUnixTime() int64 {
	return time.Now().UnixNano() / 1e6
}

//2.获取当前时间抽(秒)
func GetUnixTime() int {
	return int(time.Now().UnixNano() / 1e9)
}
