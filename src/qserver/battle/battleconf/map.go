package battleconf

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strconv"
)

//地图可走格子信息
var MapWalk map[int]map[int]int = make(map[int]map[int]int)

//mapid -> cellx: cell占用像素
var CellToPx map[int]int = make(map[int]int)

//mapid -> celly: cell占用像素
var CellToPy map[int]int = make(map[int]int)

func initMap() {
	doInitMap(2)
}
func doInitMap(mapId int) {
	mmap := make(map[int]int)
	path := "./bin/conf/Map" + strconv.Itoa(mapId) + ".txt"
	//path := "D:\\qszz\\Server\\bin\\conf\\Map" + strconv.Itoa(mapId) + ".txt"
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	y := -1
	//第一行记录cellx 和celly像素大小
	for {
		line, err := r.ReadSlice('\n')
		if err != nil && err != io.EOF {
			panic(err)
		}
		if err != nil {
			break
		}
		if y == -1 {
			line = bytes.TrimSpace(line)
			b := bytes.Split(line, []byte(","))
			CellToPx[mapId], _ = strconv.Atoi(string(b[0]))
			CellToPy[mapId], _ = strconv.Atoi(string(b[1]))
			y++
			continue
		}

		x := 0
		for lineIndex := 0; lineIndex < len(line); lineIndex++ {
			if line[lineIndex] == '0' {
				mmap[x*10000+y] = 1
				x++
			} else if line[lineIndex] == '1' {
				x++
			}
		}

		y++
	}
	MapWalk[mapId] = mmap
}
