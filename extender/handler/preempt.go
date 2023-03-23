package handler

import (
	"github.com/gin-gonic/gin"
	schedulerApi "k8s.io/kube-scheduler/extender/v1"
	"net/http"
)

func (h *Handler) Preempt(ctx *gin.Context) {

	args := &schedulerApi.ExtenderPreemptionArgs{}
	err := ctx.ShouldBindJSON(args)
	if err != nil {
		return
	}
	preemptionResult := &schedulerApi.ExtenderPreemptionResult{
		NodeNameToMetaVictims: args.NodeNameToMetaVictims,
	}

	ctx.JSON(http.StatusOK, preemptionResult)
}
