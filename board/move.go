package board

type Move struct {
	from, to, stonesHit uint64
}

func (bitBoard *BitBoard) GeneratePositions(colorToMove int) []Move {
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

func (bitBoard *BitBoard) generateKingMoves(colorToMove int) []Move {
	var resultList []Move
	return resultList
}

func (bitBoard *BitBoard) generateStoneMoves(colorToMove int) []Move {
	var resultList []Move
	var occupied = bitBoard.kings[0] | bitBoard.kings[1] | bitBoard.stones[0] | bitBoard.stones[1]
	var freeFields = legalBits &^ occupied

	if colorToMove == white {
		for dir:=5; dir <= 6; dir++ {
			var moveToCandidates = (bitBoard.stones[white] >> dir) & freeFields
			for moveToCandidates != 0 {
				var moveTo = moveToCandidates &^ (moveToCandidates-1)
				resultList = append(resultList, Move{moveTo << dir, moveTo, 0})
				moveToCandidates ^= moveTo
			}
		}
	} else {
		for dir:=5; dir <= 6; dir++ {
			var moveToCandidates = (bitBoard.stones[black] << dir) & freeFields
			for moveToCandidates != 0 {
				var moveTo = moveToCandidates &^ (moveToCandidates-1)
				resultList = append(resultList, Move{moveTo >> dir, moveTo, 0})
				moveToCandidates ^= moveTo
			}
		}
	}

	return resultList
}

//-------------------------------------------------------------------------------------------------

func (bitBoard *BitBoard) GenerateCaptureMoves(colorToMove int) []Move {
	return append(bitBoard.generateKingCaptures(colorToMove), bitBoard.generateStoneCaptures(colorToMove)...)
}

func (bitBoard *BitBoard) generateKingCaptures(colorToMove int) []Move {
	var resultList []Move
	return resultList
}

func (bitBoard *BitBoard) generateStoneCaptures(colorToMove int) []Move {
	var resultList []Move
	return resultList
}
