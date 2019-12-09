package utils

import (
	"fmt"
	"github.com/imkira/go-libav/avcodec"
	"github.com/imkira/go-libav/avfilter"
	"github.com/imkira/go-libav/avformat"
	"github.com/imkira/go-libav/avutil"
	"log"
)

func OpenInput(ctx *context, fileName string) {
	var err error

	// open format (container) context
	ctx.decFmt, err = avformat.NewContextForInput()
	if err != nil {
		log.Fatalf("Failed to open input context: %v\n", err)
	}

	// set some options for opening file
	options := avutil.NewDictionary()
	defer options.Free()
	if err := options.Set("scan_all_pmts", "1"); err != nil {
		log.Fatalf("Failed to set input options: %v\n", err)
	}
	fmt.Println("fileName", fileName)
	// open file for decoding
	if err := ctx.decFmt.OpenInput(fileName, nil, options); err != nil {
		log.Fatalf("Failed to open input file: %v\n", err)
	}

	// initialize context with stream information
	if err := ctx.decFmt.FindStreamInfo(nil); err != nil {
		log.Fatalf("Failed to find stream info: %v\n", err)
	}

	// dump streams to standard output
	ctx.decFmt.Dump(0, fileName, false)

	// prepare first video stream for decoding
	openInputVideoStream(ctx)
}

func openInputVideoStream(ctx *context) {
	var err error

	// find first video stream
	if ctx.decStream = videoStream(ctx.decFmt); ctx.decStream == nil {
		log.Fatalf("Could not find a video stream. Aborting...\n")
	}

	codecCtx := *ctx.decStream.CodecContext()

	codec := avcodec.FindDecoderByID(codecCtx.CodecID())
	if codec == nil {
		log.Fatalf("Could not find decoder: %v\n", codecCtx.CodecID())
	}
	if ctx.decCodec, err = avcodec.NewContextWithCodec(codec); err != nil {
		log.Fatalf("Failed to create codec context: %v\n", err)
	}
	//TODO: go build -tags ffmpeg33
	if err := codecCtx.CopyTo(ctx.decCodec); err != nil {
		log.Fatalf("Failed to copy codec context: %v\n", err)
	}
	if err := ctx.decCodec.SetInt64Option("refcounted_frames", 1); err != nil {
		log.Fatalf("Failed to copy codec context: %v\n", err)
	}
	if err := ctx.decCodec.OpenWithCodec(codec, nil); err != nil {
		log.Fatalf("Failed to open codec: %v\n", err)
	}

	ctx.srcFilter = addFilter(ctx, "buffer", "in")
	if err = ctx.srcFilter.SetImageSizeOption("video_size", ctx.decCodec.Width(), ctx.decCodec.Height()); err != nil {
		log.Fatalf("Failed to set filter option: %v\n", err)
	}
	if err = ctx.srcFilter.SetPixelFormatOption("pix_fmt", ctx.decCodec.PixelFormat()); err != nil {
		log.Fatalf("Failed to set filter option: %v\n", err)
	}
	if err = ctx.srcFilter.SetRationalOption("time_base", ctx.decCodec.TimeBase()); err != nil {
		log.Fatalf("Failed to set filter option: %v\n", err)
	}
	if err = ctx.srcFilter.Init(); err != nil {
		log.Fatalf("Failed to initialize buffer filter: %v\n", err)
	}

}

func videoStream(fmtCtx *avformat.Context) *avformat.Stream {
	for _, stream := range fmtCtx.Streams() {
		switch stream.CodecContext().CodecType() {
		case avutil.MediaTypeVideo:
			return stream
		}
	}
	return nil
}

func addFilter(ctx *context, name, id string) *avfilter.Context {
	filter := avfilter.FindFilterByName(name)
	if filter == nil {
		log.Fatalf("Could not find %s/%s filter\n", name, id)
	}
	fctx, err := ctx.filterGraph.AddFilter(filter, id)
	if err != nil {
		log.Fatalf("Failed to add %s/%s filter: %v\n", name, id, err)
	}
	return fctx
}

type context struct {
	decFmt    *avformat.Context
	decStream *avformat.Stream
	decCodec  *avcodec.Context
	decPkt    *avcodec.Packet
	decFrame  *avutil.Frame
	srcFilter *avfilter.Context

	encFmt     *avformat.Context
	encStream  *avformat.Stream
	encCodec   *avcodec.Context
	encIO      *avformat.IOContext
	encPkt     *avcodec.Packet
	encFrame   *avutil.Frame
	sinkFilter *avfilter.Context

	filterGraph *avfilter.Graph
}

func NewContext() (*context, error) {
	ctx := &context{}
	if err := ctx.alloc(); err != nil {
		ctx.Free()
		return nil, err
	}
	return ctx, nil
}

func (ctx *context) alloc() error {
	var err error
	if ctx.decPkt, err = avcodec.NewPacket(); err != nil {
		return err
	}
	if ctx.decFrame, err = avutil.NewFrame(); err != nil {
		return err
	}
	if ctx.encPkt, err = avcodec.NewPacket(); err != nil {
		return err
	}
	if ctx.encFrame, err = avutil.NewFrame(); err != nil {
		return err
	}
	if ctx.filterGraph, err = avfilter.NewGraph(); err != nil {
		return err
	}
	return nil
}

func (ctx *context) Free() {
	if ctx.encIO != nil {
		_ = ctx.encIO.Close()
		ctx.encIO = nil
	}
	if ctx.encFmt != nil {
		ctx.encFmt.Free()
		ctx.encFmt = nil
	}
	if ctx.filterGraph != nil {
		ctx.filterGraph.Free()
		ctx.filterGraph = nil
	}
	if ctx.encPkt != nil {
		ctx.encPkt.Free()
		ctx.encPkt = nil
	}
	if ctx.encFrame != nil {
		ctx.encFrame.Free()
		ctx.encFrame = nil
	}
	if ctx.decPkt != nil {
		ctx.decPkt.Free()
		ctx.decPkt = nil
	}
	if ctx.decFrame != nil {
		ctx.decFrame.Free()
		ctx.decFrame = nil
	}
	if ctx.decCodec != nil {
		ctx.decCodec.Free()
		ctx.decCodec = nil
	}
	if ctx.decFmt != nil {
		ctx.decFmt.CloseInput()
		ctx.decFmt.Free()
		ctx.decFmt = nil
	}
}
