package handler

import (
	"github.com/gin-gonic/gin"
	apiCoreV1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	schedulerApi "k8s.io/kube-scheduler/extender/v1"
	"net/http"
)

func (h *Handler) Predicate(ctx *gin.Context) {

	args := &schedulerApi.ExtenderArgs{}
	err := ctx.ShouldBindJSON(args)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	filterResult := &schedulerApi.ExtenderFilterResult{Error: ""}
	ctx.JSON(http.StatusOK, filterResult)

	if label := args.Pod.Labels[h.C.PodSelectorLabel]; len(label) == 0 {
		filterResult.Nodes = nil
		klog.V(5).Infof("predicate request: Pod: %+v, NodeList: %+v", args.Pod, filterResult.Nodes)
		return
	}

	filterResult.Nodes = &apiCoreV1.NodeList{}

	for _, node := range args.Nodes.Items {
		if label := node.Labels["x"]; len(label) == 0 {
			filterResult.Nodes.Items = append(filterResult.Nodes.Items, node)
		}
	}
	klog.V(5).Infof("predicate request: Pod: %+v, NodeList: %+v", args.Pod, args.Nodes)
}
