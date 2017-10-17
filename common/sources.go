package common

// http://www.ficsgames.org/download.html
// These URLS expire quickly so they need to be updated often.
var SOURCES = []string{
	"http://www.ficsgames.org/dl/ficsgamesdb_2016_chess2000_nomovetimes_1501117.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2015_chess2000_nomovetimes_1501118.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2014_chess2000_nomovetimes_1501119.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2013_chess2000_nomovetimes_1501120.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2012_chess2000_nomovetimes_1501121.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2011_chess2000_nomovetimes_1501122.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2010_chess2000_nomovetimes_1501125.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2009_chess2000_nomovetimes_1501126.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2008_chess2000_nomovetimes_1501127.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2007_chess2000_nomovetimes_1501128.pgn.bz2",
	"http://www.ficsgames.org/dl/ficsgamesdb_2006_chess2000_nomovetimes_1501123.pgn.bz2",
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
