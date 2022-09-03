package search

import (
	"fmt"
	"draughts/board"
	"draughts/math/intmath"
	"math"
)

type Node struct {
	pn, dpn   int
	isMaxNode bool
	parent    *Node
	childList []*Node
	bitBoard  board.BitBoard
}

const MaxNodesInTree = 100_000_000
const Infinite = 999_999_999

var GlobalString = ""

func initPnDpn(position board.BitBoard, endPosition bool, maxNodeColor int) (pn int, dpn int) {
	if !endPosition {
		return 1, 1
	}

	if position.ColorHasWon(maxNodeColor) {
		return 0, Infinite
	}
	return Infinite, 0
}

func (thisNode *Node) expand(colorToMove int, maxNodeColor int) {
	newPositions := thisNode.bitBoard.GeneratePositions(colorToMove)
	if len(newPositions) == 0 {
		doublePass := thisNode.parent != nil && thisNode.bitBoard == thisNode.parent.bitBoard
		newPosition := thisNode.bitBoard
		pn, dpn := initPnDpn(newPosition, thisNode.bitBoard.AllFieldsPlayed() || doublePass, maxNodeColor)
		thisNode.childList = make([]*Node, 1)
		thisNode.childList[0] = &Node{pn, dpn, !thisNode.isMaxNode, thisNode, nil, newPosition}
	} else {
		thisNode.childList = make([]*Node, len(newPositions))
		for index, newPosition := range newPositions {
			pn, dpn := initPnDpn(newPosition, newPosition.AllFieldsPlayed(), maxNodeColor)
			thisNode.childList[index] = &Node{pn, dpn, !thisNode.isMaxNode, thisNode, nil, newPosition}
		}
	}
}

func (thisNode *Node) getMostProvingChild() *Node {
	if thisNode.isMaxNode {
		for _, child := range thisNode.childList {
			if child.pn == thisNode.pn {
				return child
			}
		}
	} else {
		for _, child := range thisNode.childList {
			if child.dpn == thisNode.dpn {
				return child
			}
		}
	}
	panic("Cannot find mpn")
}

func (thisNode *Node) updatePnDpn() {
	if thisNode.isMaxNode {
		sum := 0
		min := math.MaxInt
		for _, child := range thisNode.childList {
			min = intmath.Min(min, child.pn)
			sum += child.dpn
		}
		thisNode.pn = min
		thisNode.dpn = intmath.Min(sum, Infinite)
	} else {
		sum := 0
		min := math.MaxInt
		for _, child := range thisNode.childList {
			min = intmath.Min(min, child.dpn)
			sum += child.pn
		}
		thisNode.pn = intmath.Min(sum, Infinite)
		thisNode.dpn = min
	}
}

func (thisNode *Node) findMostProvingNode() *Node {
	mpn := thisNode
	for mpn.childList != nil {
		mpn = mpn.getMostProvingChild()
	}
	return mpn
}

func (thisNode *Node) updateTree() {
	for currentNode := thisNode; currentNode != nil; currentNode = currentNode.parent {
		currentNode.updatePnDpn()
	}
}

func PnSearch(bitBoard board.BitBoard, colorToMove int) board.Move {
	nodeCount := 1
	root := Node{pn: 1, dpn: 1, isMaxNode: true, parent: nil, childList: nil, bitBoard: bitBoard}
	for root.pn != 0 && root.dpn != 0 && nodeCount < MaxNodesInTree {
		if nodeCount%1_000_000 == 0 {
			GlobalString = fmt.Sprintf("(pn, dpn) = (%d,%d). Nodes visited: %d\n", root.pn, root.dpn, nodeCount)
			for _, child := range root.childList {
				childMove := board.MoveBetween(root.bitBoard, child.bitBoard)
				var col, row = childMove.ToColRow()
				GlobalString += fmt.Sprintf("   %c%c: (pn, dpn) = (%d,%d) \n", 'A'+col, '1'+row, child.pn, child.dpn)
			}
		}
		mpn := root.findMostProvingNode()
		if mpn.isMaxNode {
			mpn.expand(colorToMove, colorToMove)
		} else {
			mpn.expand(1-colorToMove, colorToMove)
		}
		nodeCount += len(mpn.childList)
		mpn.updateTree()
	}
	GlobalString = fmt.Sprintf("(pn, dpn) = (%d,%d). Nodes visited: %d\n", root.pn, root.dpn, nodeCount)
	for _, child := range root.childList {
		childMove := board.MoveBetween(root.bitBoard, child.bitBoard)
		var col, row = childMove.ToColRow()
		GlobalString += fmt.Sprintf("   %c%c: (pn, dpn) = (%d,%d) \n", 'A'+col, '1'+row, child.pn, child.dpn)
	}
	return root.suggestedMove()
}

func (thisNode *Node) suggestedMove() board.Move {
	return board.MoveBetween(thisNode.bitBoard, thisNode.getMostProvingChild().bitBoard)
}
