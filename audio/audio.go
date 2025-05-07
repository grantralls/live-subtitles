package audio

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gordonklaus/portaudio"
)

var (
	FAILED_TO_INITIALIZE = errors.New("Failed To Initalize")
)

func StartRecordingDefaultInput() (<-chan []byte, error) {
	if err := portaudio.Initialize(); err != nil {
		return nil, fmt.Errorf("Failed to Initialize: (%T) %w", err, err)
	}

	channel := make(chan []byte)
	// chunk_size_in_bytes = chunk_duration_in_millisecond / 1000 * audio_sample_rate * 2
	stream, err := portaudio.OpenDefaultStream(1, 1, 16000, 1600, func(in []int16, out []int16, timeInfo portaudio.StreamCallbackTimeInfo, flags portaudio.StreamCallbackFlags) {
		buf := make([]byte, len(in)*2)
		binary.Encode(buf, binary.LittleEndian, in)
		channel <- buf
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to Open Stream: %w", err)
	}

	go func() {
		stream.Start()
		fmt.Println("Press ENTER to quit...")
		os.Stdin.Read([]byte{0})
		stream.Stop()
		portaudio.Terminate()
		close(channel)
	}()

	return channel, nil
}
