package bit64math

import "math/bits"

var bitCount [65536]int
var mostRightBitIndex [65536]int

func init() {
	for bitNumber := range bitCount {
		bitCount[bitNumber] = bits.OnesCount(uint(bitNumber))
		for bitPos := 0; bitPos < 16; bitPos++ {
			if (1<<bitPos)&bitNumber != 0 {
				mostRightBitIndex[bitNumber] = bitPos
				break
			}
		}
	}
}

func SmallesBit(bitValue uint64) uint64 {
	return ^(bitValue - 1) & bitValue
}

func BitCount(bitValue uint64) int {
	return bitCount[bitValue&0xffff] +
		bitCount[(bitValue>>16)&0xffff] +
		bitCount[(bitValue>>32)&0xffff] +
		bitCount[(bitValue>>48)&0xffff]
}

func MostRightBitIndex(bitValue uint64) int {
	if bitValue&0xffffffff != 0 {
		if bitValue&0xffff != 0 {
			return mostRightBitIndex[bitValue&0xffff]
		}
		return 16 + mostRightBitIndex[(bitValue>>16)&0xffff]
	} else {
		if bitValue&0xffffffffffff != 0 {
			return 32 + mostRightBitIndex[(bitValue>>32)&0xffff]
		}
		return 48 + mostRightBitIndex[(bitValue>>48)&0xffff]
	}
}
