package service

import (
	"draughts/board"
)

func GetNewBoard() (BoardModel, string) {
	initialBoard := board.InitStartBoard()
	return ToBoardModel(&initialBoard), initialBoard.ToBoardStatusString()
}

func GetBoard(boardStatusString string) (BoardModel, string) {
	currentBoard := board.StringToBitBoard(boardStatusString)
	return ToBoardModel(&currentBoard), currentBoard.ToBoardStatusString()
}

// func DoMove(boardStatusString string, col, row int) (BoardModel, string) {
// 	currentBoard := board.StringToBitBoard(boardStatusString)
// 	currentBoard.DoColRowMove(col, row)
// 	return ToBoardModel(&currentBoard), currentBoard.ToBoardStatusString()
// }

// func DoPassMove(boardStatusString string) (BoardModel, string) {
// 	currentBoard := board.StringToBitBoard(boardStatusString)
// 	currentBoard.DoPassMove()
// 	return ToBoardModel(&currentBoard), currentBoard.ToBoardStatusString()
// }

// func TakeBackLastMove(boardStatusString string) (BoardModel, string) {
// 	currentBoard := board.StringToBitBoard(boardStatusString)
// 	currentBoard.TakeBack()
// 	return ToBoardModel(&currentBoard), currentBoard.ToBoardStatusString()
// }