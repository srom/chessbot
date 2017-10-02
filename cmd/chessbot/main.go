package main

import (
	"fmt"
	"log"
	"time"

	"github.com/notnil/chess"
	"github.com/srom/chessbot/play"
	"path/filepath"
)

const DEPTH = 5

func main() {
	modelPath, err := filepath.Abs("../../model/chessbot.pb")
	if err != nil {
		log.Fatalf("Error reading path: %v", err)
	}

	model, err := play.LoadModel(modelPath)
	if err != nil {
		log.Fatalf("Error parsing model: %v", err)
	}
	defer model.Close()

	game := chess.NewGame()

	player := int64(-1)
	for game.Outcome() == chess.NoOutcome {
		start := time.Now()
		moveNodes, _ := play.Negamax(model, game, DEPTH, -1 * play.MAX_SCORE, play.MAX_SCORE, -1 * player)
		randomMoveNode := play.PickRandomMove(moveNodes)
		game.Move(randomMoveNode.Move)

		log.Printf("Num moves: %d", len(moveNodes))
		fmt.Print(game.Position().Board().Draw())
		log.Printf("Elapsed: %v", time.Since(start))
		fmt.Print("\n")
	}

	log.Printf("Outcome: %v", game.Outcome())
}
