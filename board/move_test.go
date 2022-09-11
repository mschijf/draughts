package board

import (
	"fmt"
	"testing"
	"time"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// depth   1  :     0.000015 ms -->              9
// depth   2  :     0.000003 ms -->             81
// depth   3  :     0.000038 ms -->            658
// depth   4  :     0.000258 ms -->           4265
// depth   5  :     0.001789 ms -->          27117
// depth   6  :     0.008484 ms -->         167140
// depth   7  :     0.058329 ms -->        1049442
// depth   8  :     0.310662 ms -->        6483961
// depth   9  :     1.938674 ms -->       41022423
// depth  10  :    12.458864 ms -->      258895763
//
// appr. 20.780.045 nodes/second
//
// see also https://damforum.nl/bb3/viewtopic.php?t=2308
//
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (bitBoard *BitBoard) perft(depth int, colorToMove int) int64 {
	if depth == 0 {
		return 1
	}

	moveList := bitBoard.GenerateMoves(colorToMove)
	if len(moveList) == 0 {
		return 1
	}

	var nodeCount int64 = 0
	for _, move := range moveList {
		bitBoard.doMove(&move, colorToMove)
		nodeCount += bitBoard.perft(depth-1, 1-colorToMove)
		bitBoard.undoMove(&move, colorToMove)
	}
	return nodeCount
}

func Test_bitBoard_perft(t *testing.T) {
	hb := InitStartBoard()

	tables := []struct {
		x int
		n int64
	}{
		{1, 9},
		{2, 81},
		{3, 658},
		{4, 4265},
		{5, 27117},
		{6, 167140},
		{7, 1049442},
		{8, 6483961},
		{9, 41022423},
		{10, 258895763},
	}

	for _, table := range tables {
		nodeCount := hb.bitBoard.perft(table.x, hb.colorToMove)
		if nodeCount != table.n {
			t.Errorf("Perft of %d was incorrect, got: %d, want: %d.", table.x, nodeCount, table.n)
		}
	}
}

func Test_bitBoard_perft_print(t *testing.T) {
	var hb = InitStartBoard()
	for i := 1; i < 12; i++ {
		currentTime := time.Now()
		result := hb.bitBoard.perft(i, hb.colorToMove)
		diff := time.Since(currentTime)
		fmt.Printf("depth %3d  : %12.6f ms --> %14d\n", i, diff.Seconds(), result)
	}
}
