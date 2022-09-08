package service

import (
	"draughts/board"
)

type FieldModel struct {
	FieldNumber         int    `json:"fieldNumber"`
	Color               string `json:"color"`
	IsKing              bool   `json:"isKing"`
	IsPlayableFromField bool   `json:"isPlayableFromField"`
	IsPieceToBeCaptured bool   `json:"isPieceToBeCaptured"`
	MoveToFields        []int  `json:"toFields"`
}

const maxFields = 50

type BoardModel struct {
	Fields           [maxFields]FieldModel `json:"fields"`
	ColorToMove      string                `json:"colorToMove"`
	TakeBackPossible bool                  `json:"takeBackPossible"`
	GameFinished     bool                  `json:"gameFinished"`
	BoardString      string                `json:"boardString"`
	ColorHasWon      string                `json:"colorHasWon"`
}

const whiteColor = "white"
const blackColor = "black"
const noneColor = "none"

func ToBoardModel(humanBoard *board.HumanBoard) BoardModel {
	var bm = BoardModel{}
	for field := 0; field < maxFields; field++ {
		bm.Fields[field] = getFieldModel(humanBoard, field+1)
	}
	if humanBoard.IsBlackToMove() {
		bm.ColorToMove = blackColor
	} else {
		bm.ColorToMove = whiteColor
	}
	bm.TakeBackPossible = humanBoard.HasHistory()
	bm.GameFinished = humanBoard.IsEndOfGame()
	bm.BoardString = humanBoard.ToBoardString()
	switch {
	case humanBoard.WhiteHasWon():
		bm.ColorHasWon = whiteColor
	case humanBoard.BlackHasWon():
		bm.ColorHasWon = blackColor
	default:
		bm.ColorHasWon = noneColor
	}
	return bm
}

func getFieldModel(bb *board.HumanBoard, field int) FieldModel {
	var toFields = bb.GetToFields(field)
	switch {
	case bb.IsBlackStone(field):
		return FieldModel{field, blackColor, false, bb.IsPlayableFromField(field), false, toFields}
	case bb.IsBlackKing(field):
		return FieldModel{field, blackColor, true, bb.IsPlayableFromField(field), false, toFields}
	case bb.IsWhiteStone(field):
		return FieldModel{field, whiteColor, false, bb.IsPlayableFromField(field), false, toFields}
	case bb.IsWhiteKing(field):
		return FieldModel{field, whiteColor, true, bb.IsPlayableFromField(field), false, toFields}
	default:
		return FieldModel{field, noneColor, false, false, false, toFields}
	}
}
