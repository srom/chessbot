package fetch

import (
	"compress/bzip2"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/notnil/chess"
	"github.com/srom/chessbot/common"
	"gopkg.in/freeeve/pgn.v1"
)

const DIMENSION = 768 // 8 * 8 squares * 12 piece types
const CHAN_BUFFER = 1e4

func YieldTriplets(urls <-chan string) <-chan *common.ChessBotTriplet {
	rand.Seed(time.Now().UnixNano())
	features := make(chan *common.ChessBotTriplet, CHAN_BUFFER)
	go func() {
		defer close(features)
		for url := range urls {
			log.Printf("Reading data from %v", url)
			response, err := http.Get(url)
			if err != nil {
				log.Printf("Error reading url %v: %v", url, err)
				continue
			}
			scanGames(response, features)
		}
	}()
	return features
}

func scanGames(response *http.Response, features chan *common.ChessBotTriplet) {
	defer func() {
		if r := recover(); r != nil {
		    log.Printf("Recovered: %v", r)
		}
	}()

	reader := bzip2.NewReader(response.Body)

	ps := pgn.NewPGNScanner(reader)
	for ps.Next() {
		scanGame(ps, features)
	}
}

func scanGame(ps *pgn.PGNScanner, features chan *common.ChessBotTriplet) {
	game, err := ps.Scan()
	if err != nil {
		log.Printf("Error reading game: %v", err)
		return
	}
	parseGame(game, features)
}

func parseGame(game *pgn.Game, features chan *common.ChessBotTriplet) {
	board := pgn.NewBoard()
	var parent *common.BoardBits
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

func getBoardBits(board *pgn.Board) *common.BoardBits {
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
	return &common.BoardBits{
		Pieces: pieces,
	}
}

func getTriplet(parent, observed, random *common.BoardBits) *common.ChessBotTriplet {
	return &common.ChessBotTriplet{
		Parent:   parent,
		Observed: observed,
		Random:   random,
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
	index := rand.Intn(len(moves))
	return moves[index]
}
