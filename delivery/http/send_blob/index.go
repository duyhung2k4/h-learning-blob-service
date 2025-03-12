package sendblobhandle

import (
	"app/internal/connection"
	constant "app/internal/constants"
	"app/internal/entity"
	middlewareapp "app/internal/middleware"
	routerconfig "app/internal/router_config"
	query "app/pkg/query/basic"

	"github.com/gin-gonic/gin"
)

type sendblobHandle struct {
	connect connection.Connection
	query   query.QueryService[entity.ProcessStream]
}

type SendblobHandle interface {
	StreamEncoding(ctx *gin.Context)
	InitStream(ctx *gin.Context)
}

func NewHandle() SendblobHandle {
	return &sendblobHandle{
		connect: connection.GetConnect(),
		query:   query.Register[entity.ProcessStream](),
	}
}

func Register(r *gin.Engine) {
	handle := NewHandle()

	routerconfig.AddRouter(r, routerconfig.RouterConfig{
		Method:   constant.POST_HTTP,
		Endpoint: "blob-stream/init-stream",
		Middleware: []gin.HandlerFunc{
			middlewareapp.ValidateToken,
			middlewareapp.GetProfileId,
		},
		Handle: handle.InitStream,
	})

	routerconfig.AddRouter(r, routerconfig.RouterConfig{
		Method:     constant.GET_HTTP,
		Endpoint:   "blob-stream/stream-encoding",
		Middleware: []gin.HandlerFunc{},
		Handle:     handle.StreamEncoding,
	})
}
