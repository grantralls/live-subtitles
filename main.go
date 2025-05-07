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

outer:
	for {
		select {
		case _, ok := <-channel:
			if !ok {
				break outer
			}
		}
	}
}
