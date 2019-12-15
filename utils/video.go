package utils

import (
	"fmt"
	"github.com/imkira/go-libav/avcodec"
	"github.com/imkira/go-libav/avformat"
	"github.com/imkira/go-libav/avutil"
	"log"
	m "simple-video-info/model"
)

func OpenInput(ctx *context, fileName string) {
	var err error
	ctx.decFmt, err = avformat.NewContextForInput()
	if err != nil {
		log.Fatalf("Failed to open input context: %v\n", err)
	}
	options := avutil.NewDictionary()
	defer options.Free()
	if err := options.Set("scan_all_pmts", "1"); err != nil {
		log.Fatalf("Failed to set input options: %v\n", err)
	}
	if err := ctx.decFmt.OpenInput(fileName, nil, options); err != nil {
		log.Fatalf("Failed to open input file: %v\n", err)
	}
	if err := ctx.decFmt.FindStreamInfo(nil); err != nil {
		log.Fatalf("Failed to find stream info: %v\n", err)
	}
	ctx.decFmt.Dump(0, fileName, false)
	openVideoStream(ctx)
	openAudioStream(ctx)
}

func openVideoStream(ctx *context) {
	var err error
	if ctx.decVideoStream = videoStream(ctx.decFmt); ctx.decVideoStream == nil {
		log.Fatalf("Could not find a video stream. Aborting...\n")
	}
	codecCtx := *ctx.decVideoStream.CodecContext()
	codec := avcodec.FindDecoderByID(codecCtx.CodecID())
	if codec == nil {
		log.Fatalf("Could not find decoder: %v\n", codecCtx.CodecID())
	}

	if ctx.decVideoCodec, err = avcodec.NewContextWithCodec(codec); err != nil {
		log.Fatalf("Failed to create codec context: %v\n", err)
	}
	//TODO: go build -tags ffmpeg33
	if err := codecCtx.CopyTo(ctx.decVideoCodec); err != nil {
		log.Fatalf("Failed to copy codec context: %v\n", err)
	}
	if err := ctx.decVideoCodec.SetInt64Option("refcounted_frames", 1); err != nil {
		log.Fatalf("Failed to copy codec context: %v\n", err)
	}
	if err := ctx.decVideoCodec.OpenWithCodec(codec, nil); err != nil {
		log.Fatalf("Failed to open codec: %v\n", err)
	}

	GetVideoInfo(ctx)
}

func openAudioStream(ctx *context) {
	if ctx.decAudioStream = audioStream(ctx.decFmt); ctx.decAudioStream == nil {
		log.Println("Could not find a audio stream. Aborting...")
		return
	}
	audioCodecCtx := *ctx.decAudioStream.CodecContext()
	audioCodec := avcodec.FindDecoderByID(audioCodecCtx.CodecID())
	if audioCodec == nil {
		log.Printf("Could not find decoder: %v\n", audioCodecCtx.CodecID())
	}
	var err error
	if ctx.decAudioCodec, err = avcodec.NewContextWithCodec(audioCodec); err != nil {
		log.Printf("Failed to create codec context: %v\n", err)
	}

}

func GetVideoInfo(ctx *context) m.Info {
	info := &m.Info{}
	for _, stream := range ctx.decFmt.Streams() {

		switch stream.CodecContext().CodecType() {
		case avutil.MediaTypeVideo:
			video := &m.Video{
				Name:     ctx.decVideoCodec.Codec().Name(),
				Width:    stream.CodecContext().Width(),
				Height:   stream.CodecContext().Height(),
				BitRate:  stream.CodecContext().BitRate(),
				Duration: fmt.Sprint(stream.Duration()),
			}
			info.Video = *video
		case avutil.MediaTypeAudio:
			audio := &m.Audio{
				Name:     ctx.decAudioCodec.Codec().Name(),
				BitRate:  stream.CodecContext().BitRate(),
				Duration: fmt.Sprint(stream.Duration()),
			}
			info.Audio = *audio
		}

	}
	return *info
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

func audioStream(fmtCtx *avformat.Context) *avformat.Stream {
	for _, stream := range fmtCtx.Streams() {
		switch stream.CodecContext().CodecType() {
		case avutil.MediaTypeAudio:
			return stream
		}
	}
	return nil
}

type context struct {
	decFmt         *avformat.Context
	decVideoStream *avformat.Stream
	decAudioStream *avformat.Stream
	decVideoCodec  *avcodec.Context
	decAudioCodec  *avcodec.Context
	decPkt         *avcodec.Packet
	decFrame       *avutil.Frame
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
	return nil
}

func (ctx *context) Free() {
	if ctx.decPkt != nil {
		ctx.decPkt.Free()
		ctx.decPkt = nil
	}
	if ctx.decFrame != nil {
		ctx.decFrame.Free()
		ctx.decFrame = nil
	}
	if ctx.decVideoCodec != nil {
		ctx.decVideoCodec.Free()
		ctx.decVideoCodec = nil
	}
	if ctx.decAudioCodec != nil {
		ctx.decAudioCodec.Free()
		ctx.decAudioCodec = nil
	}
	if ctx.decFmt != nil {
		ctx.decFmt.CloseInput()
		ctx.decFmt.Free()
		ctx.decFmt = nil
	}
}
