package play

import (
	"log"
	"sort"

	"github.com/notnil/chess"
)

const MAX_SCORE int64 = 1e6

func Negamax(model *Model, game *chess.Game, depth uint8, alpha, beta int64, player int64) ([]*MoveNode, int64) {
	if isGameOver(game) {
		return []*MoveNode{}, player * endScore(game)
	} else if depth <= 0 {
		boardInput := ParseBoard(game.Position().Board())
		score, err := model.Evaluate(boardInput)
		if err != nil {
			log.Fatal(err)
		}
		return []*MoveNode{}, player * score
	}

	bestScore := alpha
	bestMoveNodes := []*MoveNode{}
	for _, move := range getMoves(game) {
		gameCopy := copyGame(game)
		err := gameCopy.Move(move)
		if err != nil {
			log.Fatal(err)
		}

		moveNodes, score := Negamax(model, gameCopy, depth - 1,  -1 * beta, -1 * alpha, -1 * player)

		score = -1 * score
		if len(bestMoveNodes) == 0 || score >= bestScore {
			bestScore = score
			if score > bestScore {
				bestMoveNodes = []*MoveNode{}
			}
			children := []*MoveNode{}
			for _, moveNode := range moveNodes {
				children = append(children, moveNode)
			}
			bestMoveNodes = append(bestMoveNodes, &MoveNode{
				Move: move,
				Children: children,
				Score: bestScore,
			})
		}

		if score > alpha {
			alpha = score
		}
		if alpha >= beta {
			break
		}
	}

	return bestMoveNodes, bestScore
}

func copyGame(game *chess.Game) *chess.Game {
	moves := game.Moves()
	return chess.NewGame(moveGameForward(moves))
}

func moveGameForward(moves []*chess.Move) func(*chess.Game) {
	return func(g *chess.Game) {
		for _, move := range moves {
			g.Move(move)
		}
	}
}

func isGameOver(game *chess.Game) bool {
	return game.Outcome() != chess.NoOutcome
}

func endScore(game *chess.Game) int64 {
	outcome := game.Outcome()
	if outcome == chess.WhiteWon {
		return MAX_SCORE
	} else if outcome == chess.Draw {
		return 0
	} else {
		return -1 * MAX_SCORE
	}
}

func getMoves(game *chess.Game) []*chess.Move {
	moves := game.ValidMoves()
	sort.SliceStable(moves, func(i, j int) bool {
		return moveScore(moves[i]) < moveScore(moves[j])
	})
	return moves
}


func moveScore(move *chess.Move) uint8 {
	if move.HasTag(chess.Capture) {
		return 100
	} else if move.HasTag(chess.Check) {
		return 50
	} else {
		return 0
	}
}
