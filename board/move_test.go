package board

import (
	"testing"
)

func TestMoveBetween(t *testing.T) {
	hb := InitStartBoard()
	startBitBoard := hb.bitBoard
	bb := startBitBoard
	moveList := bb.GenerateMoves(black)
	for _, move := range moveList {

		bb.DoMove(&move, black)
		calcMove := MoveBetween(startBitBoard, bb)
		bb.UndoMove(&move, black)

		if calcMove.discPlayed != move.discPlayed || calcMove.discsFlipped != move.discsFlipped {
			t.Errorf("MoveBetween for move %d incorrect", move.discPlayed)
		}
	}

}
