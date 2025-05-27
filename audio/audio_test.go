package audio

import (
	"bytes"
	"testing"
)

func TestCallbackWithChannelOutput(t *testing.T) {
	resultStream := make(chan []byte)
	receiveAudioData := callbackWithChannelOutput(resultStream)

	tests := []struct {
		name   string
		input  []int16
		output []byte
	}{
		{
			"Ones",
			[]int16{1, 1, 1},
			[]byte{00000001, 00000000, 00000001, 00000000, 00000001, 00000000},
		},
	}

	for _, tt := range tests {
		go receiveAudioData(tt.input, nil)
		result := <-resultStream

		if !bytes.Equal(result, tt.output) {
			t.Fatalf("Test: %v failed. Expected=%v, got=%v", tt.name, tt.output, result)
		}
	}
}
