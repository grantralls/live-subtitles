package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/transcribestreaming"
	"github.com/aws/aws-sdk-go-v2/service/transcribestreaming/types"
)

func WithMockedReaderAndWriter(client *transcribestreaming.StartStreamTranscriptionEventStream) {
	client.Reader = &MockedTranscriptResultStreamReader{}
	client.Writer = &MockedAudioStreamWriter{}
}

type MockedTranscriptResultStreamReader struct{}

func (r *MockedTranscriptResultStreamReader) Events() <-chan types.TranscriptResultStream {
	return nil
}
func (r *MockedTranscriptResultStreamReader) Close() error {
	return nil
}
func (r *MockedTranscriptResultStreamReader) Err() error {
	return nil
}

type MockedAudioStreamWriter struct{}

func (w *MockedAudioStreamWriter) Send(context.Context, types.AudioStream) error {
	return nil
}
func (w *MockedAudioStreamWriter) Close() error {
	return nil
}
func (w *MockedAudioStreamWriter) Err() error {
	return nil
}
