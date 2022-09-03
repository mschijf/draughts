package service

import (
	"draughts/board"
)

type FieldModel struct {
	Col      int    `json:"col"`
	Row      int    `json:"row"`
	Color    string `json:"color"`
	Playable bool   `json:"playable"`
}

type BoardModel struct {
	Fields           [board.BoardSize][board.BoardSize]FieldModel `json:"fields"`
	ColorToMove      string                                       `json:"colorToMove"`
	TakeBackPossible bool                                         `json:"takeBackPossible"`
	GameFinished     bool                                         `json:"gameFinished"`
	MustPass         bool                                         `json:"mustPass"`
	BoardString      string                                       `json:"boardString"`
	WhiteCount       int                                          `json:"whiteCount"`
	BlackCount       int                                          `json:"blackCount"`
	ColorHasWon      string                                       `json:"colorHasWon"`
}

const whiteDiscColor = "white"
const blackDiscColor = "black"
const noneColor = "none"

func ToBoardModel(humanBoard *board.HumanBoard) BoardModel {
	var bm = BoardModel{}
	for row := 0; row < board.BoardSize; row++ {
		for col := 0; col < board.BoardSize; col++ {
			bm.Fields[row][col] = getFieldModel(humanBoard, col, row)
		}
	}
	if humanBoard.IsBlackToMove() {
		bm.ColorToMove = blackDiscColor
	} else {
		bm.ColorToMove = whiteDiscColor
	}
	bm.TakeBackPossible = humanBoard.HasHistory()
	bm.GameFinished = humanBoard.IsEndOfGame()
	bm.MustPass = humanBoard.MustPass() && !humanBoard.IsEndOfGame()
	bm.BoardString = humanBoard.ToBoardString()
	bm.WhiteCount, bm.BlackCount = humanBoard.CountDiscs()
	switch {
	case humanBoard.WhiteHasWon():
		bm.ColorHasWon = whiteDiscColor
	case humanBoard.BlackHasWon():
		bm.ColorHasWon = blackDiscColor
	default:
		bm.ColorHasWon = noneColor
	}
	return bm
}

func getFieldModel(bb *board.HumanBoard, col, row int) FieldModel {
	var discColor string
	if bb.IsBlackDisc(col, row) {
		discColor = blackDiscColor
	} else if bb.IsWhiteDisc(col, row) {
		discColor = whiteDiscColor
	} else {
		discColor = noneColor
	}
	return FieldModel{col, row, discColor, bb.IsPlayable(col, row)}
}
