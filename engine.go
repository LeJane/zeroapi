package zeroapi

import (
	"net/http"
	"path"
	"time"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type (
	GatewayEngine struct {
		Config      Config
		prefix      string
		router      *router
		ProtoSets   [][]byte
		middlewares []rest.Middleware
	}
	router struct {
		routers []RouteMapping
	}
)

func (r *router) add(s string, url string, rpcPath string, fs ...OptionFunc) {
	r.routers = append(r.routers, RouteMapping{
		Method:   s,
		Path:     url,
		RpcPath:  rpcPath,
		optionFs: fs,
	})
}

func Engine(conf Config, protoSets ...[]byte) *GatewayEngine {

	return &GatewayEngine{
		Config:    conf,
		router:    &router{},
		ProtoSets: protoSets,
	}
}

func (e *GatewayEngine) Use(middlewares ...rest.Middleware) {
	e.middlewares = append(e.middlewares, middlewares...)
}

func (e *GatewayEngine) Server(serverOptions ...rest.RunOption) *Server {
	var opts []Option
	opts = append(opts, WithHeaderProcessor(func(header http.Header) []string {
		return []string{
			"User-Agent:" + header.Get("User-Agent"),
			"X-Forwarded-For:" + header.Get("X-Forwarded-For"),
			"X-Real-IP:" + header.Get("X-Real-IP"),
			"app-version:" + header.Get("app-version"),
			"platform:" + header.Get("platform"),
			"user_id:" + header.Get("user_id"),
			"token:" + header.Get("token"),
		}
	}))
	svr := &Server{
		Server:    rest.MustNewServer(e.Config.RestConf, serverOptions...),
		upstreams: e.Upstreams(),
		timeout:   time.Duration(e.Config.CallRpcTimeoutSeconds) * time.Second,
	}
	for _, middleware := range e.middlewares {
		svr.Use(middleware)
	}
	for _, opt := range opts {
		opt(svr)
	}

	svr.Start()
	return svr
}

func (e *GatewayEngine) formatPrefix(prefix string) string {
	return path.Join(e.prefix, prefix)
}

func (e *GatewayEngine) Group(prefix string) (n *GatewayEngine) {
	n = &GatewayEngine{
		Config:    e.Config,
		prefix:    e.formatPrefix(prefix),
		router:    e.router,
		ProtoSets: e.ProtoSets,
	}
	return
}

func (e *GatewayEngine) GET(url string, handler interface{}, optionFs ...OptionFunc) {
	rpcPath := path.Join(e.Config.PbGroup, funcName(handler))
	e.router.add("get", e.formatPrefix(url), rpcPath, optionFs...)
}
func (e *GatewayEngine) POST(url string, handler interface{}, optionFs ...OptionFunc) {
	rpcPath := path.Join(e.Config.PbGroup, funcName(handler))
	e.router.add("post", e.formatPrefix(url), rpcPath, optionFs...)
}
func (e *GatewayEngine) PUT(url string, handler interface{}, optionFs ...OptionFunc) {
	rpcPath := path.Join(e.Config.PbGroup, funcName(handler))
	e.router.add("put", e.formatPrefix(url), rpcPath, optionFs...)
}
func (e *GatewayEngine) DELETE(url string, handler interface{}, optionFs ...OptionFunc) {
	rpcPath := path.Join(e.Config.PbGroup, funcName(handler))
	e.router.add("delete", e.formatPrefix(url), rpcPath, optionFs...)
}
func (e *GatewayEngine) PATCH(url string, handler interface{}, optionFs ...OptionFunc) {
	rpcPath := path.Join(e.Config.PbGroup, funcName(handler))
	e.router.add("patch", e.formatPrefix(url), rpcPath, optionFs...)
}

func (e *GatewayEngine) Upstreams() []upstream {

	return []upstream{{
		Grpc: zrpc.RpcClientConf{
			Endpoints: []string{e.Config.RpcListenOn},
			NonBlock:  true,
		},
		ProtoSets: e.ProtoSets,
		Mappings:  e.router.routers,
	}}
}
