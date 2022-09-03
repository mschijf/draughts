package board

type Move struct {
	discsFlipped, discPlayed uint64
}

func (mv *Move) isPass() bool {
	return mv.discPlayed == 0
}

func MoveBetween(bitBoardBefore, bitBoardAfter BitBoard) Move {
	discPlayed := (bitBoardAfter[0] | bitBoardAfter[1]) ^ (bitBoardBefore[0] | bitBoardBefore[1])
	if bitBoardBefore[0] < bitBoardAfter[0] {
		return Move{bitBoardBefore[0] ^ bitBoardAfter[0] ^ discPlayed, discPlayed}
	} else if bitBoardBefore[1] < bitBoardAfter[1] {
		return Move{bitBoardBefore[1] ^ bitBoardAfter[1] ^ discPlayed, discPlayed}
	}
	return Move{0, 0}
}

func (mv *Move) ToColRow() (col, row int) {
	return bitToColRow(mv.discPlayed)
}
