package preprocess

import (
	"compress/bzip2"
	"log"
	"net/http"

	"gopkg.in/freeeve/pgn.v1"
)

func YieldBoardFeaturesAndResult(urls <-chan string, done <-chan struct{}) <-chan *BoardFeaturesAndResult {
	features := make(chan *BoardFeaturesAndResult, CHAN_BUFFER)
	go func() {
		defer close(features)
		for url := range urls {
			log.Printf("Reading data from %v", url)
			response, err := http.Get(url)
			if err != nil {
				log.Printf("Error reading url %v: %v", url, err)
				continue
			}

			reader := bzip2.NewReader(response.Body)

			ps := pgn.NewPGNScanner(reader)
			for ps.Next() {
				game, err := ps.Scan()
				if err != nil {
					continue
				}

				result, err := parseResult(game.Tags["Result"])
				if err != nil {
					log.Printf("Error parsing result: %v\n", err)
					continue
				}

				board := pgn.NewBoard()
				for i, move := range game.Moves {
					if i == 0 {
						features <- parseBoard(board, result)
					}

					board.MakeMove(move)

					select {
					case <-done:
						return
					default:
						features <- parseBoard(board, result)
					}
				}
			}
		}
	}()
	return features
}
