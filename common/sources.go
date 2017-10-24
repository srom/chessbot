package common

// http://www.ficsgames.org/download.html
// These URLS expire quickly so they need to be updated often.
var SOURCES = []string{
	"http://www.ficsgames.org/dl/ficsgamesdb_2016_chess_nomovetimes_1502841.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2015_chess_nomovetimes_1502842.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2014_chess_nomovetimes_1502844.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2013_chess_nomovetimes_1502845.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2012_chess_nomovetimes_1502846.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2011_chess_nomovetimes_1502847.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2010_chess_nomovetimes_1502848.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2009_chess_nomovetimes_1502849.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2008_chess_nomovetimes_1502850.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2007_chess_nomovetimes_1502851.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2006_chess_nomovetimes_1502852.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2005_chess_nomovetimes_1502853.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2004_chess_nomovetimes_1502854.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2003_chess_nomovetimes_1502855.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2002_chess_nomovetimes_1502856.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2001_chess_nomovetimes_1502857.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2000_chess_nomovetimes_1502858.pgn.bz2",
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
