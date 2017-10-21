package play

import (
	"github.com/notnil/chess"
	"math/rand"
	"time"
)

type MoveNode struct {
	Move     *chess.Move
	Children []*MoveNode
	Score    int64
}

func PickRandomMove(moveNodes []*MoveNode) *MoveNode {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(moveNodes))
	return moveNodes[index]
}
