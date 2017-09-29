package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/srom/chessbot/preprocess"
)

const NUM_PARSERS = 4

func main() {
	done := make(chan struct{})
	defer close(done)

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	}))

	featureChannels := []<-chan *preprocess.BoardFeaturesAndResult{}

	urls := preprocess.YieldSourceURLs(done)
	for i := 0; i < NUM_PARSERS; i++ {
		featureChannels = append(
			featureChannels,
			preprocess.YieldBoardFeaturesAndResult(urls, done),
		)
	}

	preprocess.ExportFeaturesToS3(sess, done, featureChannels...)
}
