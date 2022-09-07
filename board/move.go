package board

import "draughts/math/bit64math"

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
		for dir := 5; dir <= 6; dir++ {
			var moveToCandidates = (bitBoard.stones[white] >> dir) & freeFields
			for moveToCandidates != 0 {
				var moveTo = moveToCandidates &^ (moveToCandidates - 1)
				resultList = append(resultList, Move{moveTo << dir, moveTo, 0})
				moveToCandidates ^= moveTo
			}
		}
	} else {
		for dir := 5; dir <= 6; dir++ {
			var moveToCandidates = (bitBoard.stones[black] << dir) & freeFields
			for moveToCandidates != 0 {
				var moveTo = moveToCandidates &^ (moveToCandidates - 1)
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
	var piecesHitCount = 0

	var colorOpponent = 1 - colorToMove
	var opponent = bitBoard.kings[colorOpponent] | bitBoard.stones[colorOpponent]
	var occupied = bitBoard.kings[colorToMove] | bitBoard.stones[colorToMove] | opponent
	var freeFields = legalBits &^ occupied

	for dir := 5; dir <= 6; dir++ {
		var moveToCandidates = (((bitBoard.stones[colorToMove] << dir) & opponent) << dir) & freeFields
		for moveToCandidates != 0 {
			var moveTo = moveToCandidates &^ (moveToCandidates - 1)
			var tmpList = bitBoard.generateStoneCapturesPerStone(colorToMove, moveTo>>(2*dir), moveTo, moveTo>>dir, opponent ^ (moveTo>>dir) )
			var currentHitCount = bit64math.BitCount(tmpList[0].stonesHit)
			if currentHitCount >= piecesHitCount {
				if currentHitCount > piecesHitCount {
					piecesHitCount = currentHitCount
					resultList = tmpList
				} else {
					resultList = append(resultList, tmpList...)
				}
			}

			moveToCandidates ^= moveTo
		}
		moveToCandidates = (((bitBoard.stones[colorToMove] >> dir) & opponent) >> dir) & freeFields
		for moveToCandidates != 0 {
			var moveTo = moveToCandidates &^ (moveToCandidates - 1)
			var tmpList = bitBoard.generateStoneCapturesPerStone(colorToMove, moveTo<<(2*dir), moveTo, moveTo<<dir, opponent ^ (moveTo<<dir))
			var currentHitCount = bit64math.BitCount(tmpList[0].stonesHit)
			if currentHitCount >= piecesHitCount {
				if currentHitCount > piecesHitCount {
					piecesHitCount = currentHitCount
					resultList = tmpList
				} else {
					resultList = append(resultList, tmpList...)
				}
			}
			moveToCandidates ^= moveTo
		}
	}

	for i:=0; i<len(resultList)-1; i++ {
		for j:=i+1; j<len(resultList); j++ {
			if resultList[i].from == resultList[j].from && resultList[i].to == resultList[j].to {
				resultList[j] = resultList[len(resultList)-1]
				resultList = resultList[:len(resultList)-1]
			}
		}	
	}
	return resultList
}

func (bitBoard *BitBoard) generateStoneCapturesPerStone(colorToMove int, moveFrom, currentPos, piecesHit, opponentsToHit uint64) []Move {
	var resultList = []Move{{moveFrom, currentPos, piecesHit}}
	var piecesHitCount = bit64math.BitCount(piecesHit)

	var occupied = bitBoard.kings[0] | bitBoard.kings[1] | bitBoard.stones[0] | bitBoard.stones[1]
	var freeFields = (legalBits &^ occupied) | moveFrom

	for dir := 5; dir <= 6; dir++ {
		var nextStep = (((currentPos << dir) & opponentsToHit) << dir) & freeFields
		if nextStep != 0 {
			var hitPiece = nextStep>>dir
			var tmpList = bitBoard.generateStoneCapturesPerStone(colorToMove, moveFrom, nextStep, piecesHit|hitPiece, opponentsToHit ^ hitPiece)
			var currentHitCount = bit64math.BitCount(tmpList[0].stonesHit)
			if currentHitCount >= piecesHitCount {
				if currentHitCount > piecesHitCount {
					piecesHitCount = currentHitCount
					resultList = tmpList
				} else {
					resultList = append(resultList, tmpList...)
				}
			}
		}
		nextStep = (((currentPos >> dir) & opponentsToHit) >> dir) & freeFields
		if nextStep != 0 {
			var hitPiece = nextStep<<dir
			var tmpList = bitBoard.generateStoneCapturesPerStone(colorToMove, moveFrom, nextStep, piecesHit|hitPiece, opponentsToHit ^ hitPiece)
			var currentHitCount = bit64math.BitCount(tmpList[0].stonesHit)
			if currentHitCount >= piecesHitCount {
				if currentHitCount > piecesHitCount {
					piecesHitCount = currentHitCount
					resultList = tmpList
				} else {
					resultList = append(resultList, tmpList...)
				}
			}
		}
	}
	return resultList
}

//-------------------------------------------------------------------------------------------------

func (bitBoard *BitBoard) doMove(move *Move, colorToMove int) {
	bitBoard.stones[colorToMove] ^= move.from
	bitBoard.stones[colorToMove] ^= move.to
	bitBoard.stones[1-colorToMove] &^= move.stonesHit
	bitBoard.kings[1-colorToMove] &^= move.stonesHit
}
