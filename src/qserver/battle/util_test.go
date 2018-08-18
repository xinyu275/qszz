package battle

import (
	"qserver/battle/battleconf"
	"testing"
)

type PosXY struct {
	X int
	Y int
}

func TestRandPosN(t *testing.T) {
	data := []battleconf.PosXY{
		{16, 1},
		{21, 2},
		{10, 20},
		{11, 220},
		{144, 2043},
	}
	if len(RandPosN(data, 3)) != 3 {
		t.Error("TestRandPosN error")
	}
}
