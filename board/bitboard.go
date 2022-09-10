package board

/*

BitBoard representation                                           Human Board representation

+----+----+----+----+----+----+----+----+----+----+----+          +----+----+----+----+----+----+----+----+----+----+
| ** |    |    |    |    |    |    |    |    |    |    |          |    |    |    |    |    |    |    |    |    |    |
| 00 |    | 01 |    | 02 |    | 03 |    | 04 |    | 05 |          |    | 01 |    | 02 |    | 03 |    | 04 |    | 05 |
+----+----+----+----+----+----+----+----+----+----+----+----+     +----+----+----+----+----+----+----+----+----+----+
     |    |    |    |    |    |    |    |    |    |    | ** |     |    |    |    |    |    |    |    |    |    |    |
     | 06 |    | 07 |    | 08 |    | 09 |    | 10 |    | 11 |     | 06 |    | 07 |    | 08 |    | 09 |    | 10 |    |
+----+----+----+----+----+----+----+----+----+----+----+----+     +----+----+----+----+----+----+----+----+----+----+
| ** |    |    |    |    |    |    |    |    |    |    |          |    |    |    |    |    |    |    |    |    |    |
| 11 |    | 12 |    | 13 |    | 14 |    | 15 |    | 16 |          |    | 11 |    | 12 |    | 13 |    | 14 |    | 15 |
+----+----+----+----+----+----+----+----+----+----+----+----+     +----+----+----+----+----+----+----+----+----+----+
     |    |    |    |    |    |    |    |    |    |    | ** |     |    |    |    |    |    |    |    |    |    |    |
     | 17 |    | 18 |    | 19 |    | 20 |    | 21 |    | 22 |     | 16 |    | 17 |    | 18 |    | 19 |    | 20 |    |
+----+----+----+----+----+----+----+----+----+----+----+----+     +----+----+----+----+----+----+----+----+----+----+
| ** |    |    |    |    |    |    |    |    |    |    |          |    |    |    |    |    |    |    |    |    |    |
| 22 |    | 23 |    | 24 |    | 25 |    | 26 |    | 27 |          |    | 21 |    | 22 |    | 23 |    | 24 |    | 25 |
+----+----+----+----+----+----+----+----+----+----+----+----+     +----+----+----+----+----+----+----+----+----+----+
     |    |    |    |    |    |    |    |    |    |    | ** |     |    |    |    |    |    |    |    |    |    |    |
     | 28 |    | 29 |    | 30 |    | 31 |    | 32 |    | 33 |     | 26 |    | 27 |    | 28 |    | 29 |    | 30 |    |
+----+----+----+----+----+----+----+----+----+----+----+----+     +----+----+----+----+----+----+----+----+----+----+
| ** |    |    |    |    |    |    |    |    |    |    |          |    |    |    |    |    |    |    |    |    |    |
| 33 |    | 34 |    | 35 |    | 36 |    | 37 |    | 38 |          |    | 31 |    | 32 |    | 33 |    | 34 |    | 35 |
+----+----+----+----+----+----+----+----+----+----+----+----+     +----+----+----+----+----+----+----+----+----+----+
     |    |    |    |    |    |    |    |    |    |    | ** |     |    |    |    |    |    |    |    |    |    |    |
     | 39 |    | 40 |    | 41 |    | 42 |    | 43 |    | 44 |     | 36 |    | 37 |    | 38 |    | 39 |    | 40 |    |
+----+----+----+----+----+----+----+----+----+----+----+----+     +----+----+----+----+----+----+----+----+----+----+
| ** |    |    |    |    |    |    |    |    |    |    |          |    |    |    |    |    |    |    |    |    |    |
| 44 |    | 45 |    | 46 |    | 47 |    | 48 |    | 49 |          |    | 41 |    | 42 |    | 43 |    | 44 |    | 45 |
+----+----+----+----+----+----+----+----+----+----+----+----+     +----+----+----+----+----+----+----+----+----+----+
     |    |    |    |    |    |    |    |    |    |    | ** |     |    |    |    |    |    |    |    |    |    |    |
     | 50 |    | 51 |    | 52 |    | 53 |    | 54 |    | 55 |     | 46 |    | 47 |    | 48 |    | 49 |    | 50 |    |
+----+----+----+----+----+----+----+----+----+----+----+----+     +----+----+----+----+----+----+----+----+----+----+
| ** |    | ** |    | ** |    | ** |    | ** |    | ** |
| 55 |    | 56 |    | 57 |    | 58 |    | 59 |    | 60 |
+----+----+----+----+----+----+----+----+----+----+----+
     | ** |    | ** |    | ** |
     | 61 |    | 62 |    | 63 |
     +----+----+----+----+----+


The bits with a '**' are illegal fields, but the representation has been chosen so, that a diagonal move is always
a shift of 5 or a shift of 6 positions. The constant LEGAL_FIELDS describes all legal fields, thus excluding the added bits.

So illegal bits are:

1111 1111 1000 0000 0001 0000 0000 0010 0000 0000 0100 0000 0000 1000 0000 0001

So the legal fields are the complement of the above.

This means that the fields on the human board not always 1-on-1 corresponds to the bits oin the bitboard.
Field 1 corresponds to bit 1, but field 22 corresponds to bit 23:

                       5             4            3             2             1
REAL BOARD  ---- ---- -098 7654 321- 0987 6543 21-0 9876 5432 1-09 8765 4321 -098 7654 321-

BIT INDEX   3210 9876 5432 1098 7654 3210 9876 6432 1098 7654 3210 9876 5432 1098 7654 3210
               6            5           4            3           2            1

WHITE-START 0000 0000 0111 1111 1110 1111 1111 1100 0000 0000 0000 0000 0000 0000 0000 0000
BLACK-START 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0011 1111 1111 0111 1111 1110

*/

type BitBoard struct {
	stones [2]uint64
	kings  [2]uint64
}

const white = 0
const black = 1

const illegalBits uint64 = 0xFF_80_10_02_00_40_08_01
const legalBits uint64 = ^illegalBits

const whiteStonesStartFields uint64 = 0x00_7F_EF_FC_00_00_00_00
const blackStonesStartFields uint64 = 0x00_00_00_00_00_3F_F7_FE
const whiteKingsStartFields uint64 = 0x0
const blackKingsStartFields uint64 = 0x0

var kingLine = []uint64{0x00_00_00_00_00_00_00_3E, 0x00_7C_00_00_00_00_00_00}

func InitBoard(whitePieces, blackPieces, whiteKings, blackKings uint64) BitBoard {
	var bb = BitBoard{}
	bb.stones[white] = whitePieces
	bb.stones[black] = blackPieces
	bb.kings[white] = whiteKings
	bb.kings[black] = blackKings
	return bb
}

func GetStartBoard() BitBoard {
	var bb = BitBoard{}
	bb.stones[white] = whiteStonesStartFields
	bb.stones[black] = blackStonesStartFields
	bb.kings[white] = whiteKingsStartFields
	bb.kings[black] = blackKingsStartFields
	return bb
}

// Int64 FINDPIECE (Int64 freeFields, Int64 color, Int64 fromField, int direction)
// {
// 	Int64 last, stillFree;

//     do {
//     	last = fromField;
//         fromField = SGNSHR64 (fromField, direction);
//         stillFree = AND64(fromField, freeFields);
//     } while (NOTNULL(stillFree));
//     last = SGNSHR64 (last, direction);

//     return (AND64(last, color));
// }

func FirstPieceOfColorByShiftLeft(freeFields, color, fromField uint64, shiftNumber int) uint64 {
     var last uint64
     for {
          last = fromField
          fromField <<= shiftNumber
          if fromField & freeFields == 0 {
               return last & color
          }                
     }     
}

func FirstPieceOfColorByShiftRight(freeFields, color, fromField uint64, shiftNumber int) uint64 {
     var last uint64
     for {
          last = fromField
          fromField >>= shiftNumber
          if fromField & freeFields == 0 {
               return last & color
          }                
     }     
}
