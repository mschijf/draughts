package board

import (
	"fmt"
	"testing"
	"time"
)

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

// func Test_tBitBoard_perft(t *testing.T) {
// 	hb := InitStartBoard()

// 	tables := []struct {
// 		x int
// 		n int64
// 	}{
// 		{1, 4},
// 		{2, 12},
// 		{3, 56},
// 		{4, 244},
// 		{5, 1396},
// 		{6, 8200},
// 		{7, 55092},
// 		{8, 390216},
// 		{9, 3005288},
// 		{10, 24571284},
// 	}

// 	for _, table := range tables {
// 		nodeCount := hb.bitBoard.perft(table.x, hb.colorToMove, false)
// 		if nodeCount != table.n {
// 			t.Errorf("Perft of %d was incorrect, got: %d, want: %d.", table.x, nodeCount, table.n)
// 		}
// 	}
// }

func Test_bitBoard_perft_print(t *testing.T) {
	var hb = InitStartBoard()

	for i := 1; i < 12; i++ {
		currentTime := time.Now()
		result := hb.bitBoard.perft(i, hb.colorToMove)
		diff := time.Since(currentTime)
		fmt.Printf("depth %3d  : %12.6f ms --> %14d\n", i, diff.Seconds(), result)
	}
}
