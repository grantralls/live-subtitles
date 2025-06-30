// This example shows how to use the appsrc element.
//
// Also see: https://gstreamer.freedesktop.org/documentation/tutorials/basic/short-cutting-the-pipeline.html?gi-language=c
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-gst/go-gst/pkg/gst"
	"github.com/go-gst/go-gst/pkg/gstapp"
	"github.com/grantralls/live-transcription/aws"
	"github.com/grantralls/live-transcription/gpipeline"
)

const width = 320
const height = 240

func createPipeline(dataSrc <-chan []byte) (gst.Pipeline, error) {
	println("Creating pipeline")
	gst.Init()

	p := gpipeline.New()

	// Initialize a frame counter
	var i int

	// Since our appsrc element operates in pull mode (it asks us to provide data),
	// we add a handler for the need-data callback and provide new data from there.
	// In our case, we told gstreamer that we do 2 frames per second. While the
	// buffers of all elements of the pipeline are still empty, this will be called
	// a couple of times until all of them are filled. After this initial period,
	// this handler will be called (on average) twice per second.

	var data []byte
	p.Src.ConnectNeedData(func(self gstapp.AppSrc, _ uint) {
		select {
		case textData := <-dataSrc:
			data = textData
		default:
			if data == nil {
				data = []byte("<span font=\"100\">Hello there :)</span>")
			}
		}

		// Create a buffer that can hold exactly one video RGBA frame.
		buffer := gst.NewBufferAllocate(nil, uint(len(data)), nil)

		// For each frame we produce, we set the timestamp when it should be displayed
		// The autovideosink will use this information to display the frame at the right time.
		buffer.SetPTS(p.Src.GetClock().GetTime() - p.Pipeline.GetBaseTime())
		buffer.SetDuration(gst.ClockTime(time.Millisecond))

		// At this point, buffer is only a reference to an existing memory region somewhere.
		// When we want to access its content, we have to map it while requesting the required
		// mode of access (read, read/write).
		// See: https://gstreamer.freedesktop.org/documentation/plugin-development/advanced/allocation.html
		mapped, ok := buffer.Map(gst.MapWrite)
		if !ok {
			panic("Failed to map buffer")
		}
		_, err := mapped.Write(data)
		if err != nil {
			println("Failed to write to buffer:", err)
			panic("Failed to write to buffer")
		}

		mapped.Unmap()

		// Push the buffer onto the pipeline.
		self.PushBuffer(buffer)

		i++
	})

	return p.Pipeline, nil
}

func mainLoop(pipeline gst.Pipeline) error {
	// Start the pipeline

	pipeline.SetState(gst.StatePlaying)

	for msg := range pipeline.GetBus().Messages(context.Background()) {
		switch msg.Type() {
		case gst.MessageEos:
			return nil
		case gst.MessageError:
			debug, gerr := msg.ParseError()
			if debug != "" {
				fmt.Println(gerr.Error(), debug)
			}
			return gerr
		default:
			fmt.Println(msg)
		}

		pipeline.DebugBinToDotFileWithTs(gst.DebugGraphShowVerbose, "pipeline")
	}

	return fmt.Errorf("unexpected end of messages without EOS")
}

func main() {
	run()
}

func run() {
	// audioDataChan, err := audio.StartRecordingDefaultInput()

	// if err != nil {
	// 	log.Fatalf("Error when starting audio: %v", err)
	// }

	_, stream := aws.StartStream()
	results := stream.Events()
	defer stream.Close()
	rawText := make(chan []byte)

	go func() {
		pipeline, err := createPipeline(rawText)
		if err != nil {
			fmt.Println("Error creating pipeline:", err)
			return
		}

		err = mainLoop(pipeline)

		if err != nil {
			fmt.Println("Error running pipeline:", err)
		}
	}()

	// outer:
	for {
		select {
		// case rawAudioData, ok := <-audioDataChan:
		// 	if !ok {
		// 		break outer
		// 	}
		// 	err := sender(rawAudioData)
		// 	if err != nil {
		// 		log.Printf("Error when sending audio data to aws: %v", err)
		// 		break outer
		// 	}
		case transcriptionResult := <-results:
			transcript := aws.GetTranscript(transcriptionResult)
			if transcript != nil {
				rawText <- []byte("<span font=\"50\">" + *transcript + "</span>")
			}
		}
	}
}
