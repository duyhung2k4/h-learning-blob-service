package sendblobhandle

import (
	constant "app/internal/constants"
	"app/internal/entity"
	httpresponse "app/pkg/http_response"
	logapp "app/pkg/log"
	"app/pkg/uuidapp"

	"github.com/gin-gonic/gin"
)

func (h *sendblobHandle) InitStream(ctx *gin.Context) {
	uuidStream, err := uuidapp.Create()
	if err != nil {
		logapp.Logger("create-stream", err.Error(), constant.ERROR_LOG)
		httpresponse.InternalServerError(ctx, err)
		return
	}

	profileId := ctx.GetUint("profile_id")

	// info process
	model := entity.ProcessStream{
		ProfileId: profileId,
		Uuid:      uuidStream,
		Status:    entity.PROCESS_PENDING,
	}

	result, err := h.query.Create(model)
	if err != nil {
		logapp.Logger("create-process-stream", err.Error(), constant.ERROR_LOG)
		httpresponse.InternalServerError(ctx, err)
		return
	}

	httpresponse.Success(ctx, result)
}
