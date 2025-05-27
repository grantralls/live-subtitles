package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/transcribestreaming"
	"github.com/aws/aws-sdk-go-v2/service/transcribestreaming/types"
	"github.com/grantralls/live-transcription/audio"
)

// Sends a slice of audio bytes for AWS to transcribe
type Send func([]byte) error

// Starts a stream with a given config
// Returns a function "Send" used to send bytes to aws, and an aws stream to read results
func StartStream() (Send, *transcribestreaming.StartStreamTranscriptionEventStream) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Error when loading SDK config, %v", err)
	}

	client := transcribestreaming.NewFromConfig(cfg)
	sampleRate := int32(audio.SAMPLE_RATE)
	params := transcribestreaming.StartStreamTranscriptionInput{
		MediaEncoding:        types.MediaEncodingPcm,
		MediaSampleRateHertz: &sampleRate,
		LanguageCode:         types.LanguageCodeEnUs,
	}

	output, err := client.StartStreamTranscription(context.TODO(), &params)
	if err != nil {
		log.Fatalf("Error when starting stream transcription: %v", err)
	}

	return send(output.GetStream()), output.GetStream()
}

// A helper function to get the transcript out of AWS's types
func GetTranscript(transcriptionResult types.TranscriptResultStream) *string {
	transcriptResult, ok := transcriptionResult.(*types.TranscriptResultStreamMemberTranscriptEvent)

	if !ok {
		return nil
	}

	transcript := transcriptResult.Value.Transcript
	if len(transcript.Results) > 0 && len(transcript.Results[0].Alternatives) > 0 {
		return transcript.Results[0].Alternatives[0].Transcript
	}

	return nil
}

// A helper function to wrap the audio bytes in AWS's required types
func send(stream *transcribestreaming.StartStreamTranscriptionEventStream) func([]byte) error {
	return func(data []byte) error {
		wrappedData := types.AudioStreamMemberAudioEvent{
			Value: types.AudioEvent{
				AudioChunk: data,
			},
		}
		return stream.Send(context.TODO(), &wrappedData)
	}
}
