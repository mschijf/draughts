package search

import (
	"fmt"
	"draughts/board"
	"time"
)

func ComputeMove(board board.HumanBoard) (col, row int) {
	GlobalString = ""

	currentTime := time.Now()
	computedMove := PnSearch(board.GetBitBoard(), board.GetColorToMove())
	col, row = computedMove.ToColRow()
	diff := time.Since(currentTime)
	GlobalString += fmt.Sprintf("\n\nMove played: %c%c after %12.6f ms", 'A'+col, '1'+row, diff.Seconds())

	ABNodeCount = 0
	currentTime = time.Now()
	value := AlfaBeta(board.GetBitBoard(), board.GetColorToMove(), 99, -Infinite, Infinite, false)
	diff = time.Since(currentTime)
	GlobalString += fmt.Sprintf("\n\nAlfabeta: %d after %d nodes and %12.6f ms", value, ABNodeCount, diff.Seconds())

	return col, row
}
