package server

import (
	"log"

	"github.com/LeJane/zeroapi"
	"github.com/LeJane/zeroapi/examples/hello/internal/middleware"
	"github.com/LeJane/zeroapi/examples/hello/pb"
	"github.com/zeromicro/go-zero/rest"
)

func (s *StreamGreeterServer) routeV1(group *zeroapi.GatewayEngine) {
	group.GET("/greet", s.Greet)
}

func (s *StreamGreeterServer) Gateway() {
	engine := zeroapi.Engine(s.svcCtx.Config.Gateway, pb.ProtoSetCommon, pb.ProtoSetHello)
	s.routeV1(engine.Group("/hello/v1"))
	svr := engine.Server(rest.WithCors())
	svr.Use(middleware.PrintLog)
	defer svr.Stop()
	log.Println("gateway is started at 0.0.0.0:", s.svcCtx.Config.Gateway.Port)
	svr.Start()
}
