package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const BATCH_SIZE = 1e6
const BUCKET_NAME = "chessbot"
const KEY_FORMAT = "inputs/%v.csv.gz"

type Inputs struct {
	Inputs []BoardFeaturesAndResult `json:"inputs"`
}

func ExportFeaturesToS3(sess *session.Session, done <-chan struct{}, featureChannels ...<-chan *BoardFeaturesAndResult) {
	start := time.Now()
	loopStart := time.Now()
	iteration := 0
	var wgOutChan sync.WaitGroup

	out := make(chan *BoardFeaturesAndResult, CHAN_BUFFER)
	defer close(out)

	output := func(c <-chan *BoardFeaturesAndResult) {
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

	batch := []BoardFeaturesAndResult{}
	for feature := range out {
		iteration += 1
		batch = append(batch, *feature)

		if len(batch) == BATCH_SIZE {
			flushToS3(sess, batch)
			log.Printf("%v: Elapsed %v; Batch %v", iteration, time.Since(start), time.Since(loopStart))
			loopStart = time.Now()
			batch = []BoardFeaturesAndResult{}
			break
		}
	}
	if len(batch) > 1e4 {
		flushToS3(sess, batch)
	}
}

func flushToS3(sess *session.Session, batch []BoardFeaturesAndResult) {
	r, w := io.Pipe()

	inputs := Inputs{Inputs: batch}

	uploader := s3manager.NewUploader(sess)
	keyName := fmt.Sprintf(KEY_FORMAT, time.Now().Unix())

	go func(inputs *Inputs) {
		gz := gzip.NewWriter(w)
		defer func() {
			gz.Close()
			w.Close()
		}()

		err := json.NewEncoder(gz).Encode(inputs)
		if err != nil {
			log.Printf("Error marshalling inputs to JSON: %v", err)
		}
	}(&inputs)

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
