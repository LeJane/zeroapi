package config

import (
	"github.com/LeJane/zeroapi"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Gateway zeroapi.Config
}
