package audio

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/gordonklaus/portaudio"
)

const (
	NUM_OF_CHANNELS = 1
	SAMPLE_RATE     = 16000
	// chunk_size_in_bytes = chunk_duration_in_millisecond (100ms) / 1000 * audio_sample_rate
	OneHundrdMilliS   = 50
	FRAMES_PER_BUFFER = OneHundrdMilliS / 1000 * SAMPLE_RATE
)

func StartRecordingDefaultInput() (<-chan []byte, error) {
	// Initialize PortAudio
	if err := portaudio.Initialize(); err != nil {
		return nil, fmt.Errorf("Failed to Initialize: (%T) %w", err, err)
	}

	destination := make(chan []byte)

	// The callback gets ran on a separate thread
	stream, err := portaudio.OpenDefaultStream(NUM_OF_CHANNELS, NUM_OF_CHANNELS, SAMPLE_RATE, FRAMES_PER_BUFFER, callbackWithChannelOutput(destination))

	if err != nil {
		return nil, fmt.Errorf("Failed to Open Stream: %w", err)
	}

	stream.Start()
	fmt.Println("Press ENTER to quit...")

	// Keep collecting microphone data on a separate thread until the user pressed "Enter"
	// This is so we can return the channel even though we don't want to stop the stream yet.
	go func() {
		os.Stdin.Read([]byte{0})
		stream.Stop()
		portaudio.Terminate()
		close(destination)
	}()

	return destination, nil
}

// The callback used by portaudio to capture microphone data.
// Receives a byte slice channel to forward microphone data to.
func callbackWithChannelOutput(destination chan<- []byte) func([]int16, []int16) {
	return func(in []int16, out []int16) {
		// log.Println(in)
		buf := make([]byte, len(in)*2)
		numberOfBytesWritten, err := binary.Encode(buf, binary.LittleEndian, in)

		if numberOfBytesWritten != len(in)*2 {
			panic("Unexpected number of bytes")
			return
		}

		if err != nil {
			panic("Failed to Encode microphone data into binary")
		}

		destination <- buf
	}
}
