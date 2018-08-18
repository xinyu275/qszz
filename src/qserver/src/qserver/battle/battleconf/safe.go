package battleconf

//安全区域

//意义：N秒内要把安全区域面积缩小为总面积的百分比,在等待N秒，继续一下个
type BaseSafeArea struct {
	Index        int     // 索引
	ContinueTime int     // 持续时间（秒）
	ObjAreaSize  float64 // 目标区域大小
	WaitTime     int     // 等待时间（秒）
	Val          float64 // 扣除伤害比例
}

var BaseSafeAreaList map[int]*BaseSafeArea = map[int]*BaseSafeArea{
	1: &BaseSafeArea{1, 30, 0.8, 10, 0.002},
	2: &BaseSafeArea{2, 30, 0.6, 10, 0.002},
	3: &BaseSafeArea{3, 30, 0.4, 10, 0.002},
	4: &BaseSafeArea{4, 30, 0.2, 10, 0.002},
	5: &BaseSafeArea{5, 30, 0.05, 10, 0.002},
}
