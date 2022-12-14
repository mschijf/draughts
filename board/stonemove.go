package board

import "draughts/math/bit64math"

func (bitBoard *BitBoard) generateStoneMoves(colorToMove int) []Move {
	var resultList []Move
	var occupied = bitBoard.kings[0] | bitBoard.kings[1] | bitBoard.stones[0] | bitBoard.stones[1]
	var freeFields = legalBits &^ occupied

	if colorToMove == white {
		for dir := 5; dir <= 6; dir++ {
			var moveToCandidates = (bitBoard.stones[white] >> dir) & freeFields
			for moveToCandidates != 0 {
				var moveTo = moveToCandidates &^ (moveToCandidates - 1)
				resultList = append(resultList, Move{true, moveTo << dir, moveTo, 0, 0})
				moveToCandidates ^= moveTo
			}
		}
	} else {
		for dir := 5; dir <= 6; dir++ {
			var moveToCandidates = (bitBoard.stones[black] << dir) & freeFields
			for moveToCandidates != 0 {
				var moveTo = moveToCandidates &^ (moveToCandidates - 1)
				resultList = append(resultList, Move{true, moveTo >> dir, moveTo, 0, 0})
				moveToCandidates ^= moveTo
			}
		}
	}

	return resultList
}

//-------------------------------------------------------------------------------------------------

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
			var tmpList = bitBoard.generateStoneCapturesPerStone(colorToMove, moveTo>>(2*dir), moveTo, moveTo>>dir, opponent^(moveTo>>dir))
			var currentHitCount = bit64math.BitCount(tmpList[0].stonesHit | tmpList[0].kingsHit)
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
			var tmpList = bitBoard.generateStoneCapturesPerStone(colorToMove, moveTo<<(2*dir), moveTo, moveTo<<dir, opponent^(moveTo<<dir))
			var currentHitCount = bit64math.BitCount(tmpList[0].stonesHit | tmpList[0].kingsHit)
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

	resultList = removeDuplicates(resultList)
	return resultList
}

func (bitBoard *BitBoard) generateStoneCapturesPerStone(colorToMove int, moveFrom, currentPos, piecesHit, opponentsToHit uint64) []Move {
	var resultList = []Move{{true, moveFrom, currentPos, piecesHit & bitBoard.stones[1-colorToMove], piecesHit & bitBoard.kings[1-colorToMove]}}
	var piecesHitCount = bit64math.BitCount(piecesHit)

	var occupied = bitBoard.kings[0] | bitBoard.kings[1] | bitBoard.stones[0] | bitBoard.stones[1]
	var freeFields = (legalBits &^ occupied) | moveFrom

	for dir := 5; dir <= 6; dir++ {
		var nextStep = (((currentPos << dir) & opponentsToHit) << dir) & freeFields
		if nextStep != 0 {
			var hitPiece = nextStep >> dir
			var tmpList = bitBoard.generateStoneCapturesPerStone(colorToMove, moveFrom, nextStep, piecesHit|hitPiece, opponentsToHit^hitPiece)
			var currentHitCount = bit64math.BitCount(tmpList[0].stonesHit | tmpList[0].kingsHit)
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
			var hitPiece = nextStep << dir
			var tmpList = bitBoard.generateStoneCapturesPerStone(colorToMove, moveFrom, nextStep, piecesHit|hitPiece, opponentsToHit^hitPiece)
			var currentHitCount = bit64math.BitCount(tmpList[0].stonesHit | tmpList[0].kingsHit)
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
