package board

import "draughts/math/bit64math"

func (bitBoard *BitBoard) generateKingMoves(colorToMove int) []Move {
	var resultList []Move

	var occupied = bitBoard.kings[0] | bitBoard.kings[1] | bitBoard.stones[0] | bitBoard.stones[1]
	var freeFields = legalBits &^ occupied

	for dir := 5; dir <= 6; dir++ {
		var step = 1
		for candidates := (bitBoard.kings[colorToMove] >> dir) & freeFields; candidates != 0; candidates = (candidates >> dir) & freeFields {
			var newCandidate = candidates
			for newCandidate != 0 {
				var moveTo = newCandidate &^ (newCandidate - 1)
				resultList = append(resultList, Move{false, moveTo << (dir * step), moveTo, 0, 0})
				newCandidate ^= moveTo
			}
			step++
		}
		step = 1
		for candidates := (bitBoard.kings[colorToMove] << dir) & freeFields; candidates != 0; candidates = (candidates << dir) & freeFields {
			var newCandidate = candidates
			for newCandidate != 0 {
				var moveTo = newCandidate &^ (newCandidate - 1)
				resultList = append(resultList, Move{false, moveTo >> (dir * step), moveTo, 0, 0})
				newCandidate ^= moveTo
			}
			step++
		}
	}
	return resultList
}

//-------------------------------------------------------------------------------------------------

func (bitBoard *BitBoard) generateKingCaptures(colorToMove int) []Move {
	var resultList []Move
	var piecesHitCount = 0

	var colorOpponent = 1 - colorToMove
	var opponent = bitBoard.kings[colorOpponent] | bitBoard.stones[colorOpponent]

	candidates := bitBoard.kings[colorToMove]
	for candidates != 0 {
		aKing := candidates &^ (candidates - 1)
		for dir := 5; dir <= 6; dir++ {

			tmpList := bitBoard.generateKingCapturesPerKingToLeft(colorToMove, aKing, aKing, 0, opponent, dir)
			resultList, piecesHitCount = upMergeResultList(resultList, piecesHitCount, tmpList)

			tmpList = bitBoard.generateKingCapturesPerKingToRight(colorToMove, aKing, aKing, 0, opponent, dir)
			resultList, piecesHitCount = upMergeResultList(resultList, piecesHitCount, tmpList)

		}
		candidates ^= aKing
	}

	resultList = removeDuplicates(resultList)
	return resultList
}

func (bitBoard *BitBoard) generateKingCapturesPerKingToLeft(colorToMove int, moveFrom, currentPos, piecesHit, opponentsToHit uint64, dir int) []Move {
	var resultList []Move
	var last uint64

	occupied := bitBoard.kings[0] | bitBoard.kings[1] | bitBoard.stones[0] | bitBoard.stones[1]
	freeFields := (legalBits &^ occupied) | moveFrom

	piecesHitCount := bit64math.BitCount(piecesHit)
	localHitCount := piecesHitCount
	hitCandidate := FirstPieceOfColorByShiftLeft(freeFields, opponentsToHit, currentPos, dir)
	for hitCandidate != 0 {
		moveTo := (hitCandidate << dir) & freeFields
		opponentsToHit ^= hitCandidate
		piecesHit |= hitCandidate
		localHitCount++
		for moveTo != 0 {

			if len(resultList) > 0 {
				if localHitCount >= piecesHitCount {
					if localHitCount > piecesHitCount {
						piecesHitCount = localHitCount
						resultList = []Move{{false, moveFrom, moveTo, piecesHit & bitBoard.stones[1-colorToMove], piecesHit & bitBoard.kings[1-colorToMove]}}
					} else {
						resultList = append(resultList, Move{false, moveFrom, moveTo, piecesHit & bitBoard.stones[1-colorToMove], piecesHit & bitBoard.kings[1-colorToMove]})
					}
				}
			} else {
				piecesHitCount = localHitCount
				resultList = append(resultList, Move{false, moveFrom, moveTo, piecesHit & bitBoard.stones[1-colorToMove], piecesHit & bitBoard.kings[1-colorToMove]})
			}

			tmpList := bitBoard.generateKingCapturesPerKingToLeft(colorToMove, moveFrom, moveTo, piecesHit, opponentsToHit, 11-dir)
			resultList, piecesHitCount = upMergeResultList(resultList, piecesHitCount, tmpList)

			tmpList = bitBoard.generateKingCapturesPerKingToRight(colorToMove, moveFrom, moveTo, piecesHit, opponentsToHit, 11-dir)
			resultList, piecesHitCount = upMergeResultList(resultList, piecesHitCount, tmpList)

			last = moveTo
			moveTo = (moveTo << dir) & freeFields
		}
		hitCandidate = (last << dir) & opponentsToHit
	}

	return resultList
}

func (bitBoard *BitBoard) generateKingCapturesPerKingToRight(colorToMove int, moveFrom, currentPos, piecesHit, opponentsToHit uint64, dir int) []Move {
	var resultList []Move
	var last uint64

	occupied := bitBoard.kings[0] | bitBoard.kings[1] | bitBoard.stones[0] | bitBoard.stones[1]
	freeFields := (legalBits &^ occupied) | moveFrom

	piecesHitCount := bit64math.BitCount(piecesHit)
	localHitCount := piecesHitCount
	hitCandidate := FirstPieceOfColorByShiftRight(freeFields, opponentsToHit, currentPos, dir)
	for hitCandidate != 0 {
		moveTo := (hitCandidate >> dir) & freeFields
		opponentsToHit ^= hitCandidate
		piecesHit |= hitCandidate
		localHitCount++
		for moveTo != 0 {

			if len(resultList) > 0 {
				if localHitCount >= piecesHitCount {
					if localHitCount > piecesHitCount {
						piecesHitCount = localHitCount
						resultList = []Move{{false, moveFrom, moveTo, piecesHit & bitBoard.stones[1-colorToMove], piecesHit & bitBoard.kings[1-colorToMove]}}
					} else {
						resultList = append(resultList, Move{false, moveFrom, moveTo, piecesHit & bitBoard.stones[1-colorToMove], piecesHit & bitBoard.kings[1-colorToMove]})
					}
				}
			} else {
				piecesHitCount = localHitCount
				resultList = append(resultList, Move{false, moveFrom, moveTo, piecesHit & bitBoard.stones[1-colorToMove], piecesHit & bitBoard.kings[1-colorToMove]})
			}

			tmpList := bitBoard.generateKingCapturesPerKingToLeft(colorToMove, moveFrom, moveTo, piecesHit, opponentsToHit, 11-dir)
			resultList, piecesHitCount = upMergeResultList(resultList, piecesHitCount, tmpList)

			tmpList = bitBoard.generateKingCapturesPerKingToRight(colorToMove, moveFrom, moveTo, piecesHit, opponentsToHit, 11-dir)
			resultList, piecesHitCount = upMergeResultList(resultList, piecesHitCount, tmpList)

			last = moveTo
			moveTo = (moveTo >> dir) & freeFields
		}
		hitCandidate = (last >> dir) & opponentsToHit
	}

	return resultList
}

func upMergeResultList(resultList []Move, piecesHitCount int, tmpList []Move) ([]Move, int) {
	if len(tmpList) == 0 {
		return resultList, piecesHitCount
	}
	currentHitCount := bit64math.BitCount(tmpList[0].stonesHit | tmpList[0].kingsHit)
	switch {
	case currentHitCount > piecesHitCount:
		return tmpList, currentHitCount
	case currentHitCount >= piecesHitCount:
		return append(resultList, tmpList...), piecesHitCount
	default:
		return resultList, piecesHitCount
	}
}

func removeDuplicates(resultList []Move) []Move {
	processed := map[uint64]struct{}{}
	w := 0
	for _, s := range resultList {
		if _, exists := processed[s.from|s.to]; !exists {
			processed[s.from|s.to] = struct{}{}
			resultList[w] = s
			w++
		}
	}
	return resultList[:w]
}
