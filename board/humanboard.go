package board

import (
	"fmt"
	"draughts/collection"
	"draughts/math/bit64math"
	"strconv"
	"strings"
)

//           BitBoard                               Human (0-based) bitFields
//
//  --- --- --- --- --- --- --- ---    RIJ --- --- --- --- --- --- --- ---
// |63 |62 |61 |60 |59 |58 |57 |56 |     0|   |   |   |   |   |   |   |   |
//  --- --- --- --- --- --- --- ---        --- --- --- --- --- --- --- ---
// |55 |   |   |   |   |   |   |48 |     1|   |   |   |   |   |   |   |   |
//  --- --- --- --- --- --- --- ---        --- --- --- --- --- --- --- ---
// |47 |   |   |   |   |   |   |40 |     2|   |   |   |   |   |   |   |   |
//  --- --- --- --- --- --- --- ---        --- --- --- --- --- --- --- ---
// |39 |   |   |   |   |   |   |32 |     3|   |   |   |   |   |   |   |   |
//  --- --- --- --- --- --- --- ---        --- --- --- --- --- --- --- ---
// |31 |   |   |   |   |   |   |24 |     4|   |   |   |   |   |   |   |   |
//  --- --- --- --- --- --- --- ---        --- --- --- --- --- --- --- ---
// |23 |   |   |   |   |   |   |16 |     5|   |   |   |   |   |   |   |   |
//  --- --- --- --- --- --- --- ---        --- --- --- --- --- --- --- ---
// |15 |   |   |   |   |   |   | 8 |     6|   |   |   |   |   |   |   |   |
//  --- --- --- --- --- --- --- ---        --- --- --- --- --- --- --- ---
// | 7 | 6 | 5 | 4 | 3 | 2 | 1 | 0 |     7|   |   |   |   |   |   |   |   |
//  --- --- --- --- --- --- --- ---        --- --- --- --- --- --- --- ---
//                                          0   1   2   3   4   5   6   7   KOLOM

type HumanBoard struct {
	bitBoard    BitBoard
	colorToMove int
	stack       collection.Stack[Move]
}

const delimiter = ":"

func colRowToBit(col, row int) uint64 {
	return 1 << ((BoardSize-1-row)*BoardSize + (BoardSize - 1 - col))
}

func bitToColRow(bit uint64) (int, int) {
	mrb := bit64math.MostRightBitIndex(bit)
	return (BoardSize - 1) - (mrb % BoardSize), (BoardSize - 1) - (mrb / BoardSize)
}

func StringToBitBoard(boardString string) HumanBoard {
	if boardString == "" {
		return InitStartBoard()
	}
	var boardStringParts = strings.Split(boardString, ":")
	bbWhite, _ := strconv.ParseUint(boardStringParts[0], 10, 64)
	bbBlack, _ := strconv.ParseUint(boardStringParts[1], 10, 64)
	colorToMove, _ := strconv.Atoi(boardStringParts[2])

	humanBoard := HumanBoard{bitBoard: InitBoard(bbWhite, bbBlack), colorToMove: colorToMove}

	if len(boardStringParts) == 4 {
		var moveList = boardStringParts[3]
		for i := 0; i < len(moveList); i += 2 {
			if moveList[i] == 'x' && moveList[i+1] == 'x' {
				humanBoard.DoPassMove()
			} else {
				var col = int(moveList[i] - '0')
				var row = int(moveList[i+1] - '0')
				humanBoard.DoColRowMove(col, row)
			}
		}
	}
	return humanBoard
}

func InitStartBoard() HumanBoard {
	var start = BoardSize / 2
	var bbWhite = colRowToBit(start-1, start-1) | colRowToBit(start, start) //0x00_00_00_10_08_00_00_00
	var bbBlack = colRowToBit(start-1, start) | colRowToBit(start, start-1) //0x00_00_00_08_10_00_00_00
	return HumanBoard{bitBoard: InitBoard(bbWhite, bbBlack), colorToMove: black}
}

func (hb *HumanBoard) IsBlackToMove() bool {
	return hb.colorToMove == black
}

func (hb *HumanBoard) IsWhiteDisc(col, row int) bool {
	return colRowToBit(col, row)&hb.bitBoard[white] != 0
}

func (hb *HumanBoard) IsBlackDisc(col, row int) bool {
	return colRowToBit(col, row)&hb.bitBoard[black] != 0
}

func (hb *HumanBoard) IsPlayable(col, row int) bool {
	return colRowToBit(col, row)&hb.bitBoard.GetAllCandidateMoves(hb.colorToMove) != 0
}

func (hb *HumanBoard) MustPass() bool {
	return hb.bitBoard.GetAllCandidateMoves(hb.colorToMove) == 0
}

func (hb *HumanBoard) HasHistory() bool {
	return !hb.stack.IsEmpty()
}

func (hb *HumanBoard) IsEndOfGame() bool {
	if hb.bitBoard.AllFieldsPlayed() {
		return true
	}

	return hb.stack.Size() > 1 && hb.stack.FromTop(0).isPass() && hb.stack.FromTop(1).isPass()
}

func (hb *HumanBoard) doBitBoardMove(moveBit uint64) {
	var moves = hb.bitBoard.GenerateMoves(hb.colorToMove)
	for _, move := range moves {
		if move.discPlayed == moveBit {
			hb.bitBoard.DoMove(&move, hb.colorToMove)
			hb.colorToMove = 1 - hb.colorToMove
			hb.stack.Push(&move)
			return
		}
	}
	panic("Move from UI is not correct")
}

func (hb *HumanBoard) DoColRowMove(col, row int) {
	hb.doBitBoardMove(colRowToBit(col, row))
}

func (hb *HumanBoard) DoPassMove() {
	hb.doBitBoardMove(0)
}

func (hb *HumanBoard) TakeBack() {
	move := hb.stack.Pop()
	hb.colorToMove = 1 - hb.colorToMove
	hb.bitBoard.UndoMove(move, hb.colorToMove)
}

func (hb *HumanBoard) CountDiscs() (whiteCount, blackCount int) {
	return bit64math.BitCount(hb.bitBoard[white]), bit64math.BitCount(hb.bitBoard[black])
}

func (hb *HumanBoard) WhiteHasWon() bool {
	whiteCount, blackCount := hb.CountDiscs()
	return hb.IsEndOfGame() && whiteCount > blackCount
}

func (hb *HumanBoard) BlackHasWon() bool {
	whiteCount, blackCount := hb.CountDiscs()
	return hb.IsEndOfGame() && blackCount > whiteCount
}

func (hb *HumanBoard) ToBoardString() string {
	return fmt.Sprintf("%d%s%d%s%d", hb.bitBoard[0], delimiter, hb.bitBoard[1], delimiter, hb.colorToMove)
}

func (hb *HumanBoard) ToBoardStatusString() string {
	var movesPlayedString = ""
	var tmpStack collection.Stack[Move]
	for !hb.stack.IsEmpty() {
		move := hb.stack.Top()
		tmpStack.Push(move)
		if move.isPass() {
			movesPlayedString = "xx" + movesPlayedString
		} else {
			col, row := bitToColRow(move.discPlayed)
			movesPlayedString = fmt.Sprintf("%d%d", col, row) + movesPlayedString
		}
		hb.TakeBack()
	}

	initialBoardString := hb.ToBoardString()

	for !tmpStack.IsEmpty() {
		move := tmpStack.Pop()
		hb.doBitBoardMove(move.discPlayed)
	}

	return initialBoardString + delimiter + movesPlayedString
}

func (hb *HumanBoard) GetBitBoard() BitBoard {
	return hb.bitBoard
}

func (hb *HumanBoard) GetColorToMove() int {
	return hb.colorToMove
}
