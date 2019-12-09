package utils

import (
	"github.com/imkira/go-libav/avcodec"
	"github.com/imkira/go-libav/avfilter"
	"github.com/imkira/go-libav/avformat"
	"github.com/imkira/go-libav/avutil"
)

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

func newContext() (*context, error) {
	ctx := &context{}
	if err := ctx.alloc(); err != nil {
		ctx.free()
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

func (ctx *context) free() {
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
