package sendblobhandle

import (
	"app/generated/proto/servicegrpc"
	constant "app/internal/constants"
	httpresponse "app/pkg/http_response"
	logapp "app/pkg/log"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (h *sendblobHandle) InitStream(ctx *gin.Context) {
	uuid := ctx.Query("uuid")
	ipQuantity360p := ctx.Query("quantity_360p")
	ipMergeBlob := ctx.Query("ip_merge_blob")

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		httpresponse.InternalServerError(ctx, err)
		logapp.Logger("up-connection", err.Error(), constant.ERROR_LOG)
		return
	}
	defer conn.Close()

	// quantity-grpc
	log.Println("Quantity: ", ipQuantity360p)
	log.Println("Merge blob: ", ipMergeBlob)
	connQuantityGrpc, err := grpc.NewClient(ipQuantity360p, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		httpresponse.InternalServerError(ctx, err)
		logapp.Logger("connection-quantity-grpc", err.Error(), constant.ERROR_LOG)
		return
	}
	grpcClientQuantity := servicegrpc.NewQuantityServiceClient(connQuantityGrpc)

	ctxGrpc := metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
		"ip_merge_blob": ipMergeBlob,
		"uuid":          uuid,
	}))
	stream, err := grpcClientQuantity.SendBlobQuantity(ctxGrpc)

	if err != nil {
		httpresponse.InternalServerError(ctx, err)
		logapp.Logger("stream-quantity", err.Error(), constant.ERROR_LOG)
		return
	}

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Lỗi khi nhận tin nhắn từ WebSocket: %v", err)
			break
		}

		err = stream.Send(&servicegrpc.SendBlobQuantityRequest{
			Uuid: uuid,
			Blob: data,
		})

		if err != nil {
			log.Println("Error send: ", err)
		}
	}
}
