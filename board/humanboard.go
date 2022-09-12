package board

import (
	"draughts/math/bit64math"
	"fmt"
	"log"

	"draughts/collection"
	"strconv"
	"strings"
)

type HumanBoard struct {
	bitBoard    BitBoard
	colorToMove int
	moveStack   collection.Stack[Move]
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

//-------------------------------------------------------------------------------------------------

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
	return !hb.moveStack.IsEmpty()
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

func (hb *HumanBoard) GetColorToMove() int {
	return hb.colorToMove
}

func (hb *HumanBoard) IsPlayableFromField(field int) bool {
	var moveList []Move = hb.bitBoard.GenerateMoves(hb.colorToMove)
	for _, move := range moveList {
		if fieldToBit(field) == move.from {
			return true
		}
	}
	return false
}

func (hb *HumanBoard) GetToFields(field int) []int {
	var result []int

	var moveList []Move = hb.bitBoard.GenerateMoves(hb.colorToMove)
	for _, move := range moveList {
		if fieldToBit(field) == move.from {
			result = append(result, bitToField(move.to))
		}
	}
	return result
}

//-------------------------------------------------------------------------------------------------

func InitStartBoard() HumanBoard {
	return HumanBoard{bitBoard: GetStartBoard(), colorToMove: white}
}

func InitEmptyBoard() HumanBoard {
	return HumanBoard{bitBoard: InitBoard(0, 0, 0, 0), colorToMove: white}
}

func BoardStatusStringToHumanBoard(boardString string) HumanBoard {
	if boardString == "" {
		return InitStartBoard()
	}
	var boardStringParts = strings.Split(boardString, ":")
	colorToMove, _ := strconv.Atoi(boardStringParts[0])
	whitePieces, _ := strconv.ParseUint(boardStringParts[1], 16, 64)
	blackPieces, _ := strconv.ParseUint(boardStringParts[2], 16, 64)
	kings, _ := strconv.ParseUint(boardStringParts[3], 16, 64)

	humanBoard := HumanBoard{bitBoard: InitBoard(whitePieces&^kings, blackPieces&^kings, whitePieces&kings, blackPieces&kings), colorToMove: colorToMove}

	for i := 4; i < len(boardStringParts); i++ {
		moveFrom, _ := strconv.Atoi(boardStringParts[i][0:2])
		moveTo, _ := strconv.Atoi(boardStringParts[i][2:4])
		humanBoard.DoMove(moveFrom, moveTo)
	}

	return humanBoard
}

func (hb *HumanBoard) ToBoardString() string {
	return bitBoardToBoardString(&hb.bitBoard, hb.colorToMove)
}

func bitBoardToBoardString(bb *BitBoard, colorToMove int) string {
	var whitePieces = bb.stones[1] | bb.kings[1]
	var blackPieces = bb.stones[0] | bb.kings[0]
	var kings = bb.kings[0] | bb.kings[1]
	return fmt.Sprintf("%d:%x:%x:%x", colorToMove, blackPieces, whitePieces, kings)
}

func (hb *HumanBoard) ToBoardStatusString() string {
	var tmpStack collection.Stack[Move]
	for !hb.moveStack.IsEmpty() {
		move := hb.moveStack.Top()
		tmpStack.Push(move)
		hb.TakeBack()
	}

	result := hb.ToBoardString()

	for !tmpStack.IsEmpty() {
		move := tmpStack.Pop()
		result = fmt.Sprintf("%s:%02d%02d", result, bitToField(move.from), bitToField(move.to))
		hb.DoMove(bitToField(move.from), bitToField(move.to))
	}

	return result
}

//
// Examples:
//   [FEN "B:W18,24,27,28,K10,K15:B12,16,20,K22,K25,K29"]
//   [FEN "B:W18,19,21,23,24,26,29,30,31,32:B1,2,3,4,6,7,9,10,11,12"]
//   [FEN "W:W31-50:B1-20"]
//
// See: http://pdn.fmjd.org/fen.html#fen-section
//
func FenStringToHumanBoard(fenString string) HumanBoard {
	fenString = strings.TrimSuffix(fenString, ".")
	
	if fenString == "" {
		return InitStartBoard()
	}
	
	var boardStringParts = strings.Split(fenString, ":")
	if len(boardStringParts) != 3 {
		log.Printf("FEN string '%s' does not have the expected 3 parts ", fenString)
		return InitEmptyBoard()
	}

	var whiteStones, blackStones, whiteKings, blackKings uint64
	colorToMove := getColorFromFen(boardStringParts[0])
	if colorToMove != white && colorToMove != black {
		log.Printf("FEN string '%s' incorrect colorToMove ", fenString)
		return InitEmptyBoard()
	}
	color, stones, kings := getPiecesBitStringFromFen(boardStringParts[1])
	if color == white {
		whiteStones = stones
		whiteKings = kings
	} else if color == black {
		blackStones = stones
		blackKings = kings
	} else {
		log.Printf("FEN string '%s' incorrect color in first part ", fenString)
		return InitEmptyBoard()
	}
	color, stones, kings = getPiecesBitStringFromFen(boardStringParts[2])
	if color == white {
		whiteStones = stones
		whiteKings = kings
	} else if color == black {
		blackStones = stones
		blackKings = kings
	} else {
		log.Printf("FEN string '%s' incorrect color in second part ", fenString)
		return InitEmptyBoard()
	}

	return HumanBoard{bitBoard: InitBoard(whiteStones, blackStones, whiteKings, blackKings), colorToMove: colorToMove}
}

func getColorFromFen(fenColorString string) int {
	if len(fenColorString) != 1 {
		return unknownColor
	}
	switch fenColorString[0] {
	case 'W':
		return white
	case 'B':
		return black
	default:
		return unknownColor
	}
}

func getPiecesBitStringFromFen(fenPiecesString string) (int, uint64, uint64) {
	var kings, stones uint64 = 0, 0

	if len(fenPiecesString) < 1 {
		return unknownColor, 0, 0
	}
	color := getColorFromFen(fenPiecesString[0:1])
	parts := strings.Split(fenPiecesString[1:], ",")
	for _, part := range parts {
		if len(part) > 0 {
			if part[0] == 'K' {
				kings |= partBits(part[1:])
			} else {
				stones |= partBits(part)
			}
		}
	}
	return color, stones, kings
}

func partBits(part string) uint64 {
	subParts := strings.Split(part, "-")
	if len(subParts) == 1 {
		value, _ := strconv.Atoi(subParts[0])
		if value >= 1 && value <= 50 {
			return fieldToBit(value)
		}
	} else if len(subParts) == 2 {
		fromValue, _ := strconv.Atoi(subParts[0])
		toValue, _ := strconv.Atoi(subParts[1])
		if fromValue >= 1 && fromValue <= 50 && toValue >= 1 && toValue <= 50 {
			value := uint64(0)
			for i := fromValue; i <= toValue; i++ {
				value |= fieldToBit(i)
			}
			return value
		}
	}
	return 0
}

//-------------------------------------------------------------------------------------------------

func (hb *HumanBoard) DoMove(fromField, toField int) {
	var moveList []Move = hb.bitBoard.GenerateMoves(hb.colorToMove)
	for _, move := range moveList {
		if fieldToBit(fromField) == move.from && fieldToBit(toField) == move.to {
			hb.moveStack.Push(&move)
			hb.bitBoard.doMove(&move, hb.colorToMove)
			hb.colorToMove = 1 - hb.colorToMove
			return
		}
	}
	panic("move not found")
}

func (hb *HumanBoard) TakeBack() {
	move := hb.moveStack.Pop()
	hb.colorToMove = 1 - hb.colorToMove
	hb.bitBoard.undoMove(move, hb.colorToMove)
}
