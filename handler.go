package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/proto" //lint:ignore SA1019 we have to import this because it appears in exported API
	"github.com/showurl/zeroapi/xhttp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

type ResponseHandler func(in proto.Message) (code int, msg string, data interface{})
type iFailedReason interface {
	GetFailedReason() string
}

type handlerOption struct {
	responseHandler ResponseHandler
}

type OptionFunc func(*handlerOption)

func WithResponseHandler(responseHandler ResponseHandler) OptionFunc {
	return func(o *handlerOption) {
		if responseHandler != nil {
			o.responseHandler = responseHandler
		}
	}
}

type Handler struct {
	*grpcurl.DefaultEventHandler
	responseHandler ResponseHandler
}

func NewHandler(w http.ResponseWriter, source grpcurl.DescriptorSource, optionFs ...OptionFunc) *Handler {
	h := &Handler{
		DefaultEventHandler: &grpcurl.DefaultEventHandler{
			Out:       w,
			Formatter: grpcurl.NewJSONFormatter(true, grpcurl.AnyResolverFromDescriptorSource(source)),
		},
	}
	option := &handlerOption{responseHandler: h.defaultResponseHandler}
	for _, optionF := range optionFs {
		optionF(option)
	}
	h.responseHandler = option.responseHandler
	return h
}

func (h *Handler) OnReceiveResponse(resp proto.Message) {
	h.NumResponses++
	code, msg, data := h.responseHandler(resp)
	res := xhttp.XResponse{
		Code:       int32(code),
		Msg:        msg,
		Data:       data,
		ServerTime: time.Now().UnixMilli(),
	}
	buf, _ := json.Marshal(res)
	fmt.Fprintln(h.Out, string(buf))
}

func (h *Handler) OnReceiveTrailers(stat *status.Status, md metadata.MD) {
	h.Status = stat
	if stat.Code() != codes.OK {
		fmt.Fprintln(h.Out, defaultErrBuf)
	}
}

func (h *Handler) defaultResponseHandler(in proto.Message) (code int, msg string, data interface{}) {
	if in == nil {
		return 0, "", nil
	} else {
		if failed, ok := in.(iFailedReason); ok {
			if failed.GetFailedReason() != "" {
				return -1, failed.GetFailedReason(), nil
			}
		}
		return 0, "", in
	}
}
