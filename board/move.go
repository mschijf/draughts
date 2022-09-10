package board

type Move struct {
	isStoneMove                   bool
	from, to, stonesHit, kingsHit uint64
}

func (bitBoard *BitBoard) GenerateMoves(colorToMove int) []Move {
	var resultList []Move

	resultList = bitBoard.GenerateCaptureMoves(colorToMove)
	if len(resultList) == 0 {
		resultList = bitBoard.GenerateNonCaptureMoves(colorToMove)
	}

	return resultList
}

func (bitBoard *BitBoard) GenerateNonCaptureMoves(colorToMove int) []Move {
	return append(bitBoard.generateKingMoves(colorToMove), bitBoard.generateStoneMoves(colorToMove)...)
}

func (bitBoard *BitBoard) GenerateCaptureMoves(colorToMove int) []Move {
	return append(bitBoard.generateKingCaptures(colorToMove), bitBoard.generateStoneCaptures(colorToMove)...)
}

//-------------------------------------------------------------------------------------------------

func (bitBoard *BitBoard) doMove(move *Move, colorToMove int) {
	if move.isStoneMove {
		bitBoard.stones[colorToMove] ^= move.from
		if (move.to & kingLine[colorToMove]) != 0 {
			bitBoard.kings[colorToMove] ^= move.to
		} else {
			bitBoard.stones[colorToMove] ^= move.to
		}
	} else {
		bitBoard.kings[colorToMove] ^= move.from
		bitBoard.kings[colorToMove] ^= move.to
	}
	bitBoard.stones[1-colorToMove] ^= move.stonesHit
	bitBoard.kings[1-colorToMove] ^= move.kingsHit
}

func (bitBoard *BitBoard) undoMove(move *Move, colorToMove int) {
	if move.isStoneMove {
		bitBoard.stones[colorToMove] ^= move.from
		if (move.to & kingLine[colorToMove]) != 0 {
			bitBoard.kings[colorToMove] ^= move.to
		} else {
			bitBoard.stones[colorToMove] ^= move.to
		}
	} else {
		bitBoard.kings[colorToMove] ^= move.from
		bitBoard.kings[colorToMove] ^= move.to
	}
	bitBoard.stones[1-colorToMove] ^= move.stonesHit
	bitBoard.kings[1-colorToMove] ^= move.kingsHit
}
