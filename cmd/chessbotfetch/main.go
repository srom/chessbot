package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/srom/chessbot/common"
	"github.com/srom/chessbot/estimator/fetch"
)

const NUM_PARSERS = 4

func main() {
	done := make(chan struct{})
	defer close(done)

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("eu-west-1"),
		Credentials: credentials.NewSharedCredentials("", "default"),
	}))

	featureChannels := []<-chan *fetch.ChessBotTriplet{}

	urls := common.YieldSourceURLs(done)
	for i := 0; i < NUM_PARSERS; i++ {
		featureChannels = append(
			featureChannels,
			fetch.YieldTriplets(urls),
		)
	}

	fetch.FetchData(sess, done, featureChannels...)
}
