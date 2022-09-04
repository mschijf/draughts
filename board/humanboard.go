package board

import (
	"fmt"
	// "draughts/collection"
	// "draughts/math/bit64math"
	"strconv"
	"strings"
)

type HumanBoard struct {
	bitBoard    BitBoard
	colorToMove int
	// stack       collection.Stack[Move]
}

/*
                       5             4            3             2             1
REAL BOARD  ---- ---- -098 7654 321- 0987 6543 21-0 9876 5432 1-09 8765 4321 -098 7654 321-

BIT INDEX   3210 9876 5432 1098 7654 3210 9876 6432 1098 7654 3210 9876 5432 1098 7654 3210
               6            5           4            3           2            1
*/

func fieldToBit(field int) uint64 {
	switch {
	case 1 <= field && field <= 10:
		return 1 << field
	case 11 <= field && field <= 20:
		return 1 << (field + 1)
	case 21 <= field && field <= 30:
		return 1 << (field + 2)
	case 31 <= field && field <= 40:
		return 1 << (field + 3)
	case 41 <= field && field <= 50:
		return 1 << (field + 4)
	default:
		panic(fmt.Sprintf("Field in field to bit not between 1 and 50. Field is %d", field))
	}
}

func StringToBitBoard(boardString string) HumanBoard {
	if boardString == "" {
		return InitStartBoard()
	}
	var boardStringParts = strings.Split(boardString, ":")
	colorToMove, _ := strconv.Atoi(boardStringParts[0])
	whitePieces, _ := strconv.ParseUint(boardStringParts[1], 16, 64)
	blackPieces, _ := strconv.ParseUint(boardStringParts[2], 16, 64)
	kings, _       := strconv.ParseUint(boardStringParts[3], 16, 64)

	humanBoard := HumanBoard{bitBoard: InitBoard(whitePieces &^ kings, blackPieces &^ kings, whitePieces & kings, blackPieces & kings), colorToMove: colorToMove}

	return humanBoard
}

func InitStartBoard() HumanBoard {
	return HumanBoard{bitBoard: GetStartBoard(), colorToMove: white}
}

func (hb *HumanBoard) IsBlackToMove() bool {
	return hb.colorToMove == black
}

func (hb *HumanBoard) IsWhiteStone(field int) bool {
	return fieldToBit(field)&hb.bitBoard.stones[white] != 0
}

func (hb *HumanBoard) IsBlackStone(field int) bool {
	return fieldToBit(field)&hb.bitBoard.stones[black] != 0
}

func (hb *HumanBoard) IsWhiteKing(field int) bool {
	return fieldToBit(field)&hb.bitBoard.kings[white] != 0
}

func (hb *HumanBoard) IsBlackKing(field int) bool {
	return fieldToBit(field)&hb.bitBoard.kings[black] != 0
}

func (hb *HumanBoard) HasHistory() bool {
	// return !hb.stack.IsEmpty()
	return false
}

func (hb *HumanBoard) IsEndOfGame() bool {
	return false
}

func (hb *HumanBoard) WhiteHasWon() bool {
	return false
}

func (hb *HumanBoard) BlackHasWon() bool {
	return false
}

func (hb *HumanBoard) ToBoardString() string {
	var whitePieces = hb.bitBoard.stones[1]|hb.bitBoard.kings[1]
	var blackPieces = hb.bitBoard.stones[0]|hb.bitBoard.kings[0]
	var kings = hb.bitBoard.kings[0]|hb.bitBoard.kings[1]
	var colorChar = 'w'
	if hb.IsBlackToMove() {
		colorChar = 'b'
	}
	return fmt.Sprintf("%c:%x:%x:%x", colorChar, blackPieces, whitePieces, kings)
}

func (hb *HumanBoard) ToBoardStatusString() string {
	return hb.ToBoardString()
}

func (hb *HumanBoard) GetColorToMove() int {
	return hb.colorToMove
}

func (hb *HumanBoard) IsPlayableField(field int) bool {
	var moveList []Move = hb.bitBoard.GeneratePositions(hb.colorToMove)

	for _, move := range moveList {
		if fieldToBit(field) == move.from {
			return true
		}
	}
	return false
}

