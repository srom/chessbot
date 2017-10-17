package fetch

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/golang/protobuf/proto"
)

const BATCH_SIZE = 1e3
const BUCKET_NAME = "chessbot"
const KEY_FORMAT = "triplets/%v.pb.gz"

func FetchData(awsSession *session.Session, done <-chan struct{}, featureChannels ...<-chan *ChessBotTriplet) {
	start := time.Now()
	loopStart := time.Now()
	iteration := 0
	var wgOutChan sync.WaitGroup

	out := make(chan *ChessBotTriplet, CHAN_BUFFER)
	defer close(out)

	output := func(c <-chan *ChessBotTriplet) {
		defer wgOutChan.Done()
		for feature := range c {
			select {
			case <-done:
				return
			default:
				out <- feature
			}
		}
	}

	wgOutChan.Add(len(featureChannels))
	for _, c := range featureChannels {
		go output(c)
	}
	go func() {
		wgOutChan.Wait()
		close(out)
	}()

	triplets := &ChessBotTriplets{}

	for triplet := range out {
		iteration += 1
		triplets.Triplets = append(triplets.Triplets, triplet)

		if len(triplets.Triplets) == BATCH_SIZE {
			flushToS3(awsSession, triplets)
			log.Printf("%v: Elapsed %v; Batch %v", iteration, time.Since(start), time.Since(loopStart))
			loopStart = time.Now()
			triplets = &ChessBotTriplets{}
		}
	}
	if len(triplets.Triplets) > 1e4 {
		flushToS3(awsSession, triplets)
	}
}

func flushToS3(sess *session.Session, triplets *ChessBotTriplets) {
	r, w := io.Pipe()

	uploader := s3manager.NewUploader(sess)
	keyName := fmt.Sprintf(KEY_FORMAT, time.Now().Unix())

	go func(triplets *ChessBotTriplets) {
		gz := gzip.NewWriter(w)

		defer func() {
			gz.Close()
			w.Close()
		}()

		data, err := proto.Marshal(triplets)
		if err != nil {
			fmt.Printf("Error marshaling triplets: %v", err)
			return
		}

		gz.Write(data)

	}(triplets)

	log.Println("Uploading...")
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(keyName),
		Body:   r,
	})
	if err != nil {
		log.Printf("Error uploading to S3: %v", err)
	}
	log.Println("Upload completed")
}