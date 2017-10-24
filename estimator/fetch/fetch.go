package fetch

import (
	"bufio"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/golang/protobuf/proto"
	"github.com/srom/chessbot/common"
)

const BATCH_SIZE = 1e5
const BUCKET_NAME = "chessbot"
const KEY_FORMAT = "triplets_all/%v.pb.gz"

func FetchData(awsSession *session.Session, done <-chan struct{}, featureChannels ...<-chan *common.ChessBotTriplet) {
	start := time.Now()
	loopStart := time.Now()
	iteration := 0
	var wgOutChan sync.WaitGroup

	out := make(chan *common.ChessBotTriplet, CHAN_BUFFER)
	defer close(out)

	output := func(c <-chan *common.ChessBotTriplet) {
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

	triplets := []*common.ChessBotTriplet{}

	for triplet := range out {
		iteration += 1
		triplets = append(triplets, triplet)

		if len(triplets) == BATCH_SIZE {
			flushToS3(awsSession, triplets)
			log.Printf("%v: Elapsed %v; Batch %v", iteration, time.Since(start), time.Since(loopStart))
			loopStart = time.Now()
			triplets = []*common.ChessBotTriplet{}
		}
	}
	if len(triplets) >= BATCH_SIZE {
		flushToS3(awsSession, triplets)
	}
}

func flushToS3(sess *session.Session, triplets []*common.ChessBotTriplet) {
	r, w := io.Pipe()

	uploader := s3manager.NewUploader(sess)
	keyName := fmt.Sprintf(KEY_FORMAT, time.Now().Unix())

	go func(triplets []*common.ChessBotTriplet) {
		gz := gzip.NewWriter(w)
		tripletWriter := bufio.NewWriter(gz)

		defer func() {
			gz.Close()
			w.Close()
		}()

		for _, triplet := range triplets {
			data, err := proto.Marshal(triplet)
			delimiter := make([]byte, 4)
			binary.LittleEndian.PutUint32(delimiter, uint32(len(data)))
			if err != nil {
				fmt.Printf("Error marshaling triplets: %v", err)
				return
			}
			tripletWriter.Write(append(delimiter, data...))
		}

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
