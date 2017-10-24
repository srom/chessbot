package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/notnil/chess"
	"github.com/srom/chessbot/estimator/play"
)

const DEPTH = 3

func main() {
	modelPath, err := filepath.Abs("../../model/estimator.pb")
	if err != nil {
		log.Fatalf("Error reading path: %v", err)
	}
	model, err := play.LoadModel(modelPath)
	if err != nil {
		log.Fatalf("Error parsing model: %v", err)
	}
	defer model.Close()

	game := chess.NewGame()
	for _, moveStr := range []string{
		"d4", "d5",
		"Nf3", "e6",
		"c4", "Nf6",
		"Nc3", "Be7",
		"e3", "O-O",
		"Bd3", "h6",
	} {
		game.MoveStr(moveStr)
	}


	fmt.Print(game.Position().Board().Draw())
	fmt.Print("\n")

	player := float32(-1)
	moveUnit := 0
	for game.Outcome() == chess.NoOutcome {
		moveUnit += 1
		start := time.Now()
		player = -1 * player

		if player == -1.0 {
			humanMove(game)
		} else {
			fmt.Println("Chessbot is thinking...")
			moveResult := play.Negamax(model, game, DEPTH, -1*play.MAX_SCORE, play.MAX_SCORE, player)
			err = game.Move(moveResult.Move)
		}

		if err != nil {
			log.Fatal(err)
		}

		if moveUnit%1 == 0 {
			log.Printf("Round %d", moveUnit/2)
			fmt.Print(game.Position().Board().Draw())
			log.Printf("Elapsed: %v", time.Since(start))
			fmt.Print("\n")
		}
	}

	log.Printf("Outcome: %v", game.Outcome())
	fmt.Print(game.Position().Board().Draw())
	fmt.Print("\n")
}

func humanMove(game *chess.Game) {
	for {
		moveStr, err := promptForMove()
		if err == nil {
			err = game.MoveStr(moveStr)
			if err == nil {
				return
			}
		}
	}
}

func promptForMove() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Move: ")
	move, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return move[:len(move)-1], nil
}
