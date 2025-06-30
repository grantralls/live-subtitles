package gpipeline

import (
	"github.com/go-gst/go-gst/pkg/gst"
	"github.com/go-gst/go-gst/pkg/gstapp"
)

type Pipeline struct {
	Pipeline     gst.Pipeline
	Src          gstapp.AppSrc
	Textrender   gst.Element
	Converter1   gst.Element
	Converter2   gst.Element
	Sink         gst.Element
	Compositor   gst.Element
	Queue        gst.Element
	Timeoverlay  gst.Element
	Capsfilter1  gst.Element
	Capsfilter2  gst.Element
	Uridecodebin gst.Element
	Fakesink     gst.Element
	Videotestsrc gst.Element
}

func New() *Pipeline {
	pipeline := Pipeline{}
	pipeline.init()
	pipeline.add()
	pipeline.setProperties()
	pipeline.link()
	return &pipeline
}

// Create the pipeline elements
func (p *Pipeline) init() {
	p.Pipeline = gst.NewPipeline("").(gst.Pipeline)

	p.Src = gst.ElementFactoryMake("appsrc", "").(gstapp.AppSrc)
	p.Textrender = gst.ElementFactoryMake("textrender", "")
	p.Converter1 = gst.ElementFactoryMake("videoconvert", "")
	p.Converter2 = gst.ElementFactoryMake("videoconvert", "")
	p.Sink = gst.ElementFactoryMake("autovideosink", "")
	p.Compositor = gst.ElementFactoryMake("compositor", "")
	p.Queue = gst.ElementFactoryMake("queue", "")
	p.Timeoverlay = gst.ElementFactoryMake("timeoverlay", "")
	p.Videotestsrc = gst.ElementFactoryMake("videotestsrc", "")
	p.Capsfilter1 = gst.ElementFactoryMake("capsfilter", "")
	p.Capsfilter2 = gst.ElementFactoryMake("capsfilter", "")
	p.Uridecodebin = gst.ElementFactoryMake("uridecodebin", "")
	p.Fakesink = gst.ElementFactoryMake("fakesink", "")
}

// Add the elements to the pipeline bin
func (p *Pipeline) add() {
	p.Pipeline.AddMany(p.Src, p.Textrender, p.Converter1, p.Compositor, p.Converter2, p.Sink, p.Capsfilter1, p.Capsfilter2)
}

func (p *Pipeline) setProperties() {
	videoCaps := gst.CapsFromString("video/x-raw,width=1020,height=436")
	p.Capsfilter1.SetObjectProperty("caps", videoCaps)
	p.Capsfilter2.SetObjectProperty("caps", videoCaps)

	caps := gst.CapsFromString("text/x-raw, format=(string)pango-markup")
	p.Src.SetObjectProperty("caps", caps)
	p.Src.SetObjectProperty("format", gst.FormatTime)
	p.Src.SetObjectProperty("is-live", true)
}

func (p *Pipeline) requestPads() {
}

func (p *Pipeline) link() {
	gst.LinkMany(p.Src, p.Textrender, p.Converter1)
	gst.LinkMany(p.Compositor, p.Converter2, p.Sink)
	p.Converter1.Link(p.Compositor)
}
