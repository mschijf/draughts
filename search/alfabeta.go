package search

import (
	"draughts/board"
	"math/rand"
)

func endPositionValue(bitBoard board.BitBoard, colorToMove int) int {
	// if bitBoard.ColorHasWon(colorToMove) {
	// 	return 1
	// } else {
	// 	return -1
	// }
	return rand.Intn(2)
}

var ABNodeCount int
func AlfaBeta(bitBoard board.BitBoard, colorToMove int, depth int, alfa, beta int, hasPassed bool) int {
	ABNodeCount++
	if depth == 0 {
		return 0
	}

	if bitBoard.AllFieldsPlayed() {
		return rand.Intn(2) //endPositionValue(bitBoard, colorToMove)
	}

	positionList := bitBoard.GeneratePositions(colorToMove)
	if len(positionList) == 0 {
		if hasPassed {
			return rand.Intn(2) //endPositionValue(bitBoard, colorToMove)
		}
		return -AlfaBeta(bitBoard, 1-colorToMove, depth-1, -beta, -alfa, true)
	}

	bestValue := alfa
	for _, pos := range positionList {
		value := -AlfaBeta(pos, 1-colorToMove, depth-1, -beta, -bestValue, false)
		if value > bestValue {
			bestValue = value
			if bestValue >= beta {
				break
			}
		}
	}
	return bestValue
}