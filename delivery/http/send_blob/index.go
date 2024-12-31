package sendblobhandle

import (
	constant "app/internal/constants"
	routerconfig "app/internal/router_config"

	"github.com/gin-gonic/gin"
)

type sendblobHandle struct{}

type SendblobHandle interface {
	InitStream(ctx *gin.Context)
}

func NewHandle() SendblobHandle {
	return &sendblobHandle{}
}

func Register(r *gin.Engine) {
	handle := NewHandle()

	routerconfig.AddRouter(r, routerconfig.RouterConfig{
		Method:     constant.GET_HTTP,
		Endpoint:   "blob-stream/init-stream",
		Middleware: []gin.HandlerFunc{},
		Handle:     handle.InitStream,
	})
}
