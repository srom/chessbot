package play

import (
	"log"
	"sort"
	"sync"

	"github.com/notnil/chess"
	"github.com/srom/chessbot/common"
)

const MAX_SCORE int64 = 1e6

type MoveResult struct {
	Move  *chess.Move
	Score int64
}

func Negamax(model *Model, game *chess.Game, depth uint8, alpha, beta int64, player int64) ([]*MoveNode, int64) {
	moves := getMovesWithEval(model, game, player)

	moveResults := make(chan MoveResult, len(moves))
	var wg sync.WaitGroup
	for _, move := range moves {
		wg.Add(1)
		go func(move *chess.Move) {
			defer wg.Done()
			gameCopy, _ := common.CopyGame(game)
			err := gameCopy.Move(move)
			if err != nil {
				log.Fatal(err)
			}
			score := negamaxSync(model, gameCopy, depth - 1, alpha, beta, player, -player)
			moveResults <- MoveResult{
				Move: move,
				Score: score,
			}
		}(move)
	}
	wg.Wait()
	close(moveResults)

	bestScore := alpha
	bestMoveNodes := []*MoveNode{}
	for moveResult := range moveResults {
		score := moveResult.Score
		if len(bestMoveNodes) == 0 || score >= bestScore {
			if score > bestScore {
				bestMoveNodes = []*MoveNode{}
			}
			bestScore = score
			bestMoveNodes = append(bestMoveNodes, &MoveNode{
				Move: moveResult.Move,
				Score: score,
			})
		}
	}

	return bestMoveNodes, bestScore
}

func negamaxSync(model *Model, game *chess.Game, depth uint8, alpha, beta int64, player, currentPlayer int64) int64 {
	if isGameOver(game) {
		return player * endScore(game)
	} else if depth <= 0 {
		boardInput := common.ParseBoard(game.Position().Board())
		score, err := model.Evaluate(boardInput)
		if err != nil {
			log.Fatal(err)
		}
		return player * score
	}

	bestScore := alpha
	for _, move := range getMoves(game) {
		gameCopy, _ := common.CopyGame(game)
		err := gameCopy.Move(move)
		if err != nil {
			log.Fatal(err)
		}

		score := negamaxSync(model, gameCopy, depth - 1,  -beta, -alpha, player, -player)

		if score >= bestScore {
			bestScore = score
		}

		if score > alpha {
			alpha = score
		}
		if alpha >= beta {
			break
		}
	}

	return bestScore
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
		return -MAX_SCORE
	}
}

func getMovesWithEval(model *Model, game *chess.Game, player int64) []*chess.Move {
	goodMoves := []*chess.Move{}
	bestScore := -MAX_SCORE
	for _, move := range getMoves(game) {
		gameCopy, _ := common.CopyGame(game)
		gameCopy.Move(move)
		boardInput := common.ParseBoard(gameCopy.Position().Board())
		score, err := model.Evaluate(boardInput)
		if err != nil {
			log.Fatalf("Error evaluating move %v: %v", move.String(), err)
		}
		score = player * score
		if score >= bestScore {
			if score > bestScore {
				goodMoves = []*chess.Move{}
			}
			bestScore = score
			goodMoves = append(goodMoves, move)
		}
	}
	return goodMoves
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
