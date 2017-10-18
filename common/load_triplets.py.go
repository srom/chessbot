package common

import (
	"github.com/golang/protobuf/proto"
	"encoding/binary"
	"bufio"
	"os"
	"log"
)

func YieldTriplets() {
	// WIP
	f, err := os.Open("/Users/srom/Downloads/1508288318.pb")
	if err != nil {
		log.Fatalln("Error opening file:", err)
	}
	reader := bufio.NewReader(f)
	sizeBytes, err := reader.Peek(4)
	if err != nil {
		log.Printf("Error reading size bytes: %v", err)
	}

	size := binary.LittleEndian.Uint32(sizeBytes)
	log.Printf("Size: %d", size)
	reader.Discard(4)
	message := make([]byte, size)
	_, err = reader.Read(message)
	if err != nil {
		log.Printf("Error reading message bytes: %v", err)
	}
	triplet := &ChessBotTriplet{}
	err = proto.Unmarshal(message, triplet)
	if err != nil {
		log.Printf("Error unmarshaling message: %v", err)
	}
	log.Printf("Random: %v (%d)", triplet.Random.Pieces, len(triplet.Random.Pieces))
}
