package main

import (
	"log"

	"github.com/grantralls/live-transcription/audio"
)

func main() {
	channel, err := audio.StartRecordingDefaultInput()

	if err != nil {
		log.Fatalf("Error when starting audio: ", err.Error())
	}

	log.Println("waiting for data")

outer:
	for {
		select {
		case data, ok := <-channel:
			if !ok {
				break outer
			}
			log.Println(data)
		}
	}
}
