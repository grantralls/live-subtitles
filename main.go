package main

import (
	"log"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/app"

	"github.com/grantralls/live-transcription/audio"
	"github.com/grantralls/live-transcription/aws"
)

func run() {
	audioDataChan, err := audio.StartRecordingDefaultInput()

	if err != nil {
		log.Fatalf("Error when starting audio: %v", err)
	}

	sender, stream := aws.StartStream()
	results := stream.Events()
	defer stream.Close()

outer:
	for {
		select {
		case rawAudioData, ok := <-audioDataChan:
			if !ok {
				break outer
			}
			err := sender(rawAudioData)
			if err != nil {
				log.Printf("Error when sending audio data to aws: %v", err)
				break outer
			}
		case transcriptionResult := <-results:
			transcript := aws.GetTranscript(transcriptionResult)
			if transcript != nil {
				log.Println(*transcript)
			}
		}
	}
}

type Vehicle interface {
	Vroom()
}

func main() {

	// data := []byte("Hello there!")

	// run()
	gst.Init(nil)
	mainLoop := glib.NewMainLoop(glib.MainContextDefault(), false)

	appSrc, err := gst.NewElement("appsrc")
	checkError("Error when creating appSrc: %v", err)

	src := app.SrcFromElement(appSrc)
	if src == nil {
		log.Fatalf("appSrc is not an appsrc element")
	}

	// videosrc, err := gst.NewElement("v4l2src")
	// checkError("Error when creating video test source: %v", err)
	//
	// videoConvert, err := gst.NewElement("videoconvert")
	// checkError("Error when creating videoconvert: %v", err)
	//
	// capsFilter, err := gst.NewElement("capsfilter")
	// checkError("Error when creating capsFilter: %v", err)
	// capsFilter.SetArg("caps", "video/x-raw,width=1920,height=1080")
	//
	// textOverlay, err := gst.NewElement("textoverlay")
	// checkError("Error when creating textOverlay: %v", err)
	// textOverlay.SetArg("text", "Hello there!!")
	//
	// autoVideoSink, err := gst.NewElement("autovideosink")
	// checkError("Error when creating video sink: %v", err)

	pipeline, err := gst.NewPipeline("pipeline")
	checkError("Error when creating a new pipeline: %v", err)

	ok := pipeline.GetBus().AddWatch(messageHandler)
	if !ok {
		log.Fatalf("failed to add watch to pipeline")
	}

	fakeSink, _ := gst.NewElement("fakesink")
	err = pipeline.AddMany(appSrc, fakeSink)
	checkError("Error when adding many: %v", err)

	err = gst.ElementLinkMany(appSrc, fakeSink)
	checkError("Error when linking many: %v", err)

	err = pipeline.SetState(gst.StatePlaying)
	checkError("Error when setting state: %v", err)
	mainLoop.Run()
}

func messageHandler(msg *gst.Message) bool {
	if msg != nil {
		log.Printf("%+v\n", msg)
	}

	return true
}

func checkError(msg string, err error) {
	if err != nil {
		log.Fatalf(msg, err)
	}
}
