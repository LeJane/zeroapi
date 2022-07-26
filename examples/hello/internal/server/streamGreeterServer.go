// Code generated by goctl. DO NOT EDIT!
// Source: hello.proto

package server

import (
	"context"
	"github.com/showurl/zeroapi/examples/hello/internal/logic"
	"github.com/showurl/zeroapi/examples/hello/internal/svc"
	"github.com/showurl/zeroapi/examples/hello/pb"
)

type StreamGreeterServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedStreamGreeterServer
}

func NewStreamGreeterServer(svcCtx *svc.ServiceContext) *StreamGreeterServer {
	return &StreamGreeterServer{
		svcCtx: svcCtx,
	}
}

func (s *StreamGreeterServer) Greet(ctx context.Context, in *pb.StreamReq) (*pb.StreamResp, error) {
	l := logic.NewGreetLogic(ctx, s.svcCtx)
	return l.Greet(in)
}