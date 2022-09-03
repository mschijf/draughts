package board

import (
	"fmt"
	"testing"
	"time"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//depth   1  :     0.000000 ms -->              4
//depth   2  :     0.000001 ms -->             12
//depth   3  :     0.000003 ms -->             56
//depth   4  :     0.000013 ms -->            244
//depth   5  :     0.000075 ms -->           1396
//depth   6  :     0.000459 ms -->           8200
//depth   7  :     0.002425 ms -->          55092
//depth   8  :     0.015443 ms -->         390216
//depth   9  :     0.104869 ms -->        3005288
//depth  10  :     0.827167 ms -->       24571284
//depth  11  :     6.935265 ms -->      212258800
//depth  12  :    62.507371 ms -->     1939886636
//depth  13  :   583.527171 ms -->    18429641748
//
//    highest speed: 31.583.176 per second (go version go1.18.2 darwin/arm64)
//
//   see also http://www.aartbik.com/strategy.php
//
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (bitBoard *BitBoard) perft(depth int, colorToMove int, justPassed bool) int64 {
	if depth == 0 {
		return 1
	}
	if bitBoard.AllFieldsPlayed() {
		return 1
	}
	var nodeCount int64 = 0
	positionList := bitBoard.GeneratePositions(colorToMove)
	if len(positionList) == 0 {
		if justPassed {
			return 1
		} else {
			return bitBoard.perft(depth-1, 1-colorToMove, true)
		}
	} else {
		for _, newPosition := range positionList {
			nodeCount += newPosition.perft(depth-1, 1-colorToMove, false)
		}
	}
	return nodeCount
}

func Test_tBitBoard_perft(t *testing.T) {
	hb := InitStartBoard()

	tables := []struct {
		x int
		n int64
	}{
		{1, 4},
		{2, 12},
		{3, 56},
		{4, 244},
		{5, 1396},
		{6, 8200},
		{7, 55092},
		{8, 390216},
		{9, 3005288},
		{10, 24571284},
	}

	for _, table := range tables {
		nodeCount := hb.bitBoard.perft(table.x, hb.colorToMove, false)
		if nodeCount != table.n {
			t.Errorf("Perft of %d was incorrect, got: %d, want: %d.", table.x, nodeCount, table.n)
		}
	}
}

func Test_bitBoard_perft_print(t *testing.T) {
	var hb = InitStartBoard()

	for i := 1; i < 12; i++ {
		currentTime := time.Now()
		result := hb.bitBoard.perft(i, hb.colorToMove, false)
		diff := time.Now().Sub(currentTime)
		fmt.Printf("depth %3d  : %12.6f ms --> %14d\n", i, diff.Seconds(), result)
	}
}
