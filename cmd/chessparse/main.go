package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/srom/chess"
)

const NUM_PARSERS = 4

func main() {
	done := make(chan struct{})
	defer close(done)

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	}))

	featureChannels := []<-chan *chess.BoardFeaturesAndResult{}

	urls := chess.YieldSourceURLs(done)
	for i := 0; i < NUM_PARSERS; i++ {
		featureChannels = append(
			featureChannels,
			chess.YieldBoardFeaturesAndResult(urls, done),
		)
	}

	chess.ExportFeaturesToS3(sess, done, featureChannels...)
}
