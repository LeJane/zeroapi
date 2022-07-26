package server

import (
	"github.com/showurl/zeroapi"
	"github.com/showurl/zeroapi/examples/hello/internal/middleware"
	"github.com/showurl/zeroapi/examples/hello/pb"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
)

func (s *StreamGreeterServer) routeV1(group *zeroapi.GatewayEngine) {
	group.GET("/greet", s.Greet)
}

func (s *StreamGreeterServer) Gateway() {
	engine := zeroapi.Engine(rest.RestConf{
		ServiceConf: s.svcCtx.Config.ServiceConf,
		Host:        "0.0.0.0",
		Port:        s.svcCtx.Config.Gateway.Port,
		Timeout:     s.svcCtx.Config.Gateway.CallRpcTimeoutSeconds * 1000,
	}, s.svcCtx.Config.Gateway, [][]byte{pb.ProtoSetCommon, pb.ProtoSetHello})
	s.routeV1(engine.Group("/hello/v1"))
	svr := engine.Server(zeroapi.WithHeaderProcessor(func(header http.Header) []string {
		return []string{
			"User-Agent:" + header.Get("User-Agent"),
			"X-Forwarded-For:" + header.Get("X-Forwarded-For"),
			"X-Real-IP:" + header.Get("X-Real-IP"),
			"token:" + header.Get("token"),
		}
	}))
	svr.Use(middleware.PrintLog)
	defer svr.Stop()
	svr.Start()
}