package service

import (
	"draughts/board"
)

func GetNewBoard() (BoardModel, string) {
	initialBoard := board.InitStartBoard()
	return ToBoardModel(&initialBoard), initialBoard.ToBoardStatusString()
}

func GetBoard(boardStatusString string) (BoardModel, string) {
	currentBoard := board.BoardStatusStringToBitBoard(boardStatusString)
	return ToBoardModel(&currentBoard), currentBoard.ToBoardStatusString()
}

func DoMove(boardStatusString string, from, to int) (BoardModel, string) {
	currentBoard := board.BoardStatusStringToBitBoard(boardStatusString)
	currentBoard.DoMove(from, to)
	return ToBoardModel(&currentBoard), currentBoard.ToBoardStatusString()
}

func TakeBackLastMove(boardStatusString string) (BoardModel, string) {
	currentBoard := board.BoardStatusStringToBitBoard(boardStatusString)
	currentBoard.TakeBack()
	return ToBoardModel(&currentBoard), currentBoard.ToBoardStatusString()
}
