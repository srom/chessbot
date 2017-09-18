package main

import (
	"compress/bzip2"
	"log"
	"net/http"

	"gopkg.in/freeeve/pgn.v1"
)

var SOURCES = []string{
	"http://www.ficsgames.org/dl/ficsgamesdb_201708_chess_nomovetimes_1494519.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_201707_chess_nomovetimes_1494520.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_201706_chess_nomovetimes_1494521.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_201705_chess_nomovetimes_1494522.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_201704_chess_nomovetimes_1494523.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_201703_chess_nomovetimes_1494524.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_201702_chess_nomovetimes_1494525.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_201701_chess_nomovetimes_1494526.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2016_chess_nomovetimes_1494493.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2015_chess_nomovetimes_1494494.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2014_chess_nomovetimes_1494497.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2013_chess_nomovetimes_1494498.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2012_chess_nomovetimes_1494499.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2011_chess_nomovetimes_1494501.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2010_chess_nomovetimes_1494502.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2009_chess_nomovetimes_1494503.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2008_chess_nomovetimes_1494504.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2007_chess_nomovetimes_1494505.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2006_chess_nomovetimes_1494506.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2005_chess_nomovetimes_1494507.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2004_chess_nomovetimes_1494508.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2003_chess_nomovetimes_1494509.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2002_chess_nomovetimes_1494510.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2001_chess_nomovetimes_1494511.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2000_chess_nomovetimes_1494513.pgn.bz2",
}

func YieldSourceURLs(done <-chan struct{}) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for _, url := range SOURCES {
			select {
			case <-done:
				return
			default:
				out <- url
			}
		}
	}()
	return out
}

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
