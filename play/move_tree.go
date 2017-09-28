package play

import "github.com/notnil/chess"

type MoveNode struct {
	Move     chess.Move
	Children []*MoveNode
	Score    int32
}
