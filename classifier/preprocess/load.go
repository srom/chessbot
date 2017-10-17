package preprocess

import (
	"compress/bzip2"
	"log"
	"net/http"

	"gopkg.in/freeeve/pgn.v1"
)

// http://www.ficsgames.org/download.html
var SOURCES = []string{
	"http://www.ficsgames.org/dl/ficsgamesdb_2016_chess_nomovetimes_1495633.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2015_chess_nomovetimes_1495634.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2014_chess_nomovetimes_1495635.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2013_chess_nomovetimes_1495484.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2012_chess_nomovetimes_1495631.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2011_chess_nomovetimes_1495632.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2010_chess_nomovetimes_1495636.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2009_chess_nomovetimes_1495637.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2008_chess_nomovetimes_1495638.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2007_chess_nomovetimes_1495639.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2006_chess_nomovetimes_1495640.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2005_chess_nomovetimes_1495641.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2004_chess_nomovetimes_1495643.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2003_chess_nomovetimes_1495644.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2002_chess_nomovetimes_1495645.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2001_chess_nomovetimes_1495646.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2000_chess_nomovetimes_1495647.pgn.bz2",
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
