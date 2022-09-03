package board

const BoardSize = 6

func init() {
	if BoardSize%2 != 0 {
		panic("BoardSize must be even!!")
	}
}
