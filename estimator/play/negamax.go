package play

import (
	"log"
	"sort"
	"sync"

	"github.com/notnil/chess"
	"github.com/srom/chessbot/common"
	"fmt"
)

const MAX_SCORE float32 = 1e6

type MoveResult struct {
	Move  *chess.Move
	Score float32
	Moves []*chess.Move
}

func Negamax(
	model *ChessBotModel,
	game *chess.Game,
	depth uint8,
	alpha,
	beta,
	player float32,
) (*MoveResult) {
	moves := game.ValidMoves()

	boards := [][]float32{}
	for _, move := range moves {
		gameCopy, err := common.CopyGame(game)
		if err != nil {
			log.Fatal(err)
		}
		err = gameCopy.Move(move)
		if err != nil {
			log.Fatal(err)
		}
		boardInput := ParseBoard(gameCopy.Position().Board())
		boards = append(boards, boardInput)
	}
	scores, err := model.Evaluate(boards)
	if err != nil {
		log.Fatal(err)
	}
	results := []MoveResult{}
	for i, move := range moves {
		results = append(results, MoveResult{
			Move: move,
			Score: scores[i],
		})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score < results[j].Score
	})
	for _, res := range results {
		fmt.Printf("%v %f\n", res.Move, res.Score)
	}

	moveResults := make(chan *MoveResult, len(moves))
	var wg sync.WaitGroup
	count := 0
	for _, move := range moves {
		count++
		if count > 4 {
			wg.Wait()
			count = 0
		}
		wg.Add(1)
		go func(move *chess.Move) {
			defer wg.Done()
			gameCopy, err := common.CopyGame(game)
			if err != nil {
				log.Fatal(err)
			}
			err = gameCopy.Move(move)
			if err != nil {
				log.Fatal(err)
			}
			score, bestMoves := NegamaxSync(model, gameCopy, depth - 1, alpha, beta, -player)
			moveResults <- &MoveResult{
				Move: move,
				Score: -score,
				Moves: append([]*chess.Move{move}, bestMoves...),
			}
		}(move)
	}
	wg.Wait()
	close(moveResults)

	bestScore := alpha
	var bestMoveResult *MoveResult
	for moveResult := range moveResults {
		score := moveResult.Score
		if bestMoveResult == nil || score >= bestScore {
			bestMoveResult = moveResult
			bestScore = score
		}
	}

	gameCopy, err := common.CopyGame(game)
	if err != nil {
		log.Fatal(err)
	}
	for _, move := range bestMoveResult.Moves {
		gameCopy.Move(move)
	}
	log.Printf(
		"Moves: %v; Score: %f\nLAST BOARD:%v",
		bestMoveResult.Moves,
		bestMoveResult.Score,
		gameCopy.Position().Board().Draw(),
	)

	return bestMoveResult
}

func NegamaxSync(
	model *ChessBotModel,
	game *chess.Game,
	depth uint8,
	alpha,
	beta,
	player float32,
) (float32, []*chess.Move) {
	if isGameOver(game) {
		return endScore(game, player), []*chess.Move{}
	} else if depth <= 0 {
		boardInput := ParseBoard(game.Position().Board())
		score, err := model.Evaluate([][]float32{boardInput})
		if err != nil {
			log.Fatal(err)
		}
		//log.Printf("%v", score[0])
		return score[0], []*chess.Move{}
	}

	bestScore := alpha
	bestMoves := []*chess.Move{}
	for _, move := range getMoves(model, game) {
		//log.Printf("\t%d %v", depth, move)
		gameCopy, err := common.CopyGame(game)
		if err != nil {
			log.Fatal(err)
		}
		err = gameCopy.Move(move)
		if err != nil {
			log.Fatal(err)
		}

		negamaxScore, cBestMoves := NegamaxSync(model, gameCopy, depth - 1, -beta, -alpha, -player)

		score := negamaxScore

		if len(bestMoves) == 0 || score >= bestScore {
			bestScore = score
			bestMoves = append([]*chess.Move{move}, cBestMoves...)
		}

		if score > alpha {
			alpha = score
		}
		if alpha >= beta {
			break
		}
	}

	if depth == 5 {
		gameCopy, err := common.CopyGame(game)
		if err != nil {
			log.Fatal(err)
		}
		for _, move := range bestMoves {
			gameCopy.Move(move)
		}

		//log.Printf(
		//	"Moves: %v; Score: %f\nLAST BOARD:%v",
		//	bestMoves,
		//	bestScore,
		//	gameCopy.Position().Board().Draw(),
		//)
	}

	return bestScore, bestMoves
}

func isGameOver(game *chess.Game) bool {
	return game.Outcome() != chess.NoOutcome
}

func endScore(game *chess.Game, player float32) float32 {
	outcome := game.Outcome()
	if outcome == chess.WhiteWon && player == 1 {
		return MAX_SCORE
	} else if outcome == chess.WhiteWon && player == -1 {
		return -MAX_SCORE
	} else if outcome == chess.BlackWon && player == 1 {
		return -MAX_SCORE
	} else if outcome == chess.BlackWon && player == -1 {
		return MAX_SCORE
	} else if outcome == chess.Draw {
		return -MAX_SCORE + 1000
	} else {
		panic("Unexpected outcome")
	}
}

func getMoves(model *ChessBotModel, game *chess.Game) []*chess.Move {
	moves := game.ValidMoves()
	boards := [][]float32{}
	for _, move := range moves {
		gameCopy, err := common.CopyGame(game)
		if err != nil {
			log.Fatal(err)
		}
		err = gameCopy.Move(move)
		if err != nil {
			log.Fatal(err)
		}
		boardInput := ParseBoard(gameCopy.Position().Board())
		boards = append(boards, boardInput)
	}
	scores, err := model.Evaluate(boards)
	if err != nil {
		log.Fatal(err)
	}
	results := []MoveResult{}
	for i, move := range moves {
		results = append(results, MoveResult{
			Move: move,
			Score: scores[i],
		})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	movesSorted := []*chess.Move{}
	for _, mr := range results {
		movesSorted = append(movesSorted, mr.Move)
	}
	return movesSorted[:3]
}

func getMovesCheap(model *ChessBotModel, game *chess.Game) []*chess.Move {
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
