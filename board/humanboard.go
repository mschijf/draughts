package board

import (
	"draughts/math/bit64math"
	"fmt"

	// "draughts/collection"
	// "draughts/math/bit64math"
	"strconv"
	"strings"
)

type HumanBoard struct {
	bitBoard     BitBoard
	colorToMove  int
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

func bitToField(bit uint64) int {
	var bitIndex = bit64math.MostRightBitIndex(bit)
	switch {
	case 1 <= bitIndex && bitIndex <= 10:
		return bitIndex
	case 12 <= bitIndex && bitIndex <= 21:
		return bitIndex - 1
	case 23 <= bitIndex && bitIndex <= 32:
		return bitIndex - 2
	case 34 <= bitIndex && bitIndex <= 43:
		return bitIndex - 3
	case 45 <= bitIndex && bitIndex <= 54:
		return bitIndex - 4
	default:
		panic(fmt.Sprintf("Bitindex for bit cannot be mapoped to field. Bit is %x, bitindex = %d", bit, bitIndex))
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
	kings, _ := strconv.ParseUint(boardStringParts[3], 16, 64)

	humanBoard := HumanBoard{bitBoard: InitBoard(whitePieces&^kings, blackPieces&^kings, whitePieces&kings, blackPieces&kings), colorToMove: colorToMove}

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
	var whitePieces = hb.bitBoard.stones[1] | hb.bitBoard.kings[1]
	var blackPieces = hb.bitBoard.stones[0] | hb.bitBoard.kings[0]
	var kings = hb.bitBoard.kings[0] | hb.bitBoard.kings[1]
	return fmt.Sprintf("%d:%x:%x:%x", hb.colorToMove, blackPieces, whitePieces, kings)
}

func (hb *HumanBoard) ToBoardStatusString() string {
	return hb.ToBoardString()
}

func (hb *HumanBoard) OpponentColor() int {
	return 1 - hb.colorToMove
}

func (hb *HumanBoard) GetColorToMove() int {
	return hb.colorToMove
}

func (hb *HumanBoard) IsPlayableFromField(field int) bool {
	var moveList []Move = hb.bitBoard.GeneratePositions(hb.colorToMove)
	for _, move := range moveList {
		if fieldToBit(field) == move.from {
			return true
		}
	}
	return false
}

func (hb *HumanBoard) GetToFields(field int) []int {
	var result []int

	var moveList []Move = hb.bitBoard.GeneratePositions(hb.colorToMove)
	for _, move := range moveList {
		if fieldToBit(field) == move.from {
			result = append(result, bitToField(move.to))
		}
	}
	return result
}

//-------------------------------------------------------------------------------------------------

func (hb *HumanBoard) DoMove(fromField, toField int) {
	//todo verify if there is a field that is touched, corresponding with fromField
	var moveList []Move = hb.bitBoard.GeneratePositions(hb.colorToMove)
	for _, move := range moveList {
		if fieldToBit(fromField) == move.from && fieldToBit(toField) == move.to {
			hb.bitBoard.doMove(&move, hb.colorToMove)
			hb.colorToMove = hb.OpponentColor()
			return
		}
	}
	panic("move not found")
}
