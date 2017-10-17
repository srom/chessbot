package fetch

import (
	"compress/bzip2"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/srom/chessbot/common"
	"gopkg.in/freeeve/pgn.v1"
	"github.com/notnil/chess"
)

const DIMENSION = 768 // 8 * 8 squares * 12 piece types
const CHAN_BUFFER = 1e4

func YieldTriplets(urls <-chan string) <-chan *ChessBotTriplet {
	features := make(chan *ChessBotTriplet, CHAN_BUFFER)
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
					log.Printf("Error reading game: %v", err)
					continue
				}
				parseGame(game, features)
			}
		}
	}()
	return features
}

func parseGame(game *pgn.Game, features chan *ChessBotTriplet) {
	board := pgn.NewBoard()
	var parent *BoardBits
	for i, move := range game.Moves {
		if i == 0 {
			parent = getBoardBits(board)
		}

		pgnBoard := pgn.FENFromBoard(board).String()
		randomNextBoard, err := getRandomNextBoard(pgnBoard)
		if err != nil {
			log.Printf("Error getting random board: %v", err)
			return
		}

		random := getBoardBits(randomNextBoard)

		board.MakeMove(move)

		observed := getBoardBits(board)

		features <- getTriplet(parent, observed, random)
	}
}

func getBoardBits(board *pgn.Board) *BoardBits {
	pieces := make([]uint32, 0, DIMENSION)
	for _, piece := range common.PIECES {
		for _, position := range common.POSITIONS {
			if board.GetPiece(position) == piece {
				pieces = append(pieces, 1)
			} else {
				pieces = append(pieces, 0)
			}
		}
	}
	return &BoardBits{
		Pieces: pieces,
	}
}

func getTriplet(parent, observed, random *BoardBits) *ChessBotTriplet {
	return &ChessBotTriplet{
		Parent: parent,
		Observed: observed,
		Random: random,
	}
}

func getRandomNextBoard(boardFen string) (*pgn.Board, error) {
	fen, err := chess.FEN(boardFen)
	if err != nil {
		return nil, err
	}
	chessGame := chess.NewGame(fen)
	moves := chessGame.ValidMoves()
	move := randomMove(moves)
	chessGame.Move(move)
	return pgn.NewBoardFEN(chessGame.FEN())
}

func randomMove(moves []*chess.Move) *chess.Move {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(moves))
	return moves[index]
}
