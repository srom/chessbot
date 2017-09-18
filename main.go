package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

const NUM_PARSERS = 4
const CHAN_BUFFER = 1e4

func main() {
	done := make(chan struct{})
	defer close(done)

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	}))

	featureChannels := []<-chan *BoardFeaturesAndResult{}

	urls := YieldSourceURLs(done)
	for i := 0; i < NUM_PARSERS; i++ {
		featureChannels = append(
			featureChannels,
			YieldBoardFeaturesAndResult(urls, done),
		)
	}

	ExportFeaturesToS3(sess, done, featureChannels...)
}
