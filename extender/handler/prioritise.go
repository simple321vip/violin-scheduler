package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	apiCoreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	schedulerApi "k8s.io/kube-scheduler/extender/v1"
	"net/http"
	"violin-home.cn/violin-extender/common"
)

func (h *Handler) Prioritise(ctx *gin.Context) {

	args := &schedulerApi.ExtenderArgs{}
	err := ctx.ShouldBindJSON(args)

	if err != nil {
		common.LG.Error(err.Error())
	}
	var nodeNameList []string
	for _, node := range args.Nodes.Items {
		nodeNameList = append(nodeNameList, node.Name)
	}

	var priorityList schedulerApi.HostPriorityList
	priorityList = make([]schedulerApi.HostPriority, len(args.Nodes.Items))
	podList, err := h.Client.CoreV1().Pods("").List(context.TODO(), metaV1.ListOptions{
		LabelSelector: h.C.PodSelectorLabelKey,
	})
	zoneMap := h.GetZoneMap(args.Nodes.Items)
	maxScore := len(podList.Items)

	common.LG.Info("")

	for i, node := range args.Nodes.Items {

		// by node
		podCount := 0
		for _, pod := range podList.Items {
			if pod.Spec.NodeName == node.Name {
				podCount++
			}
		}
		scoreByNode := (maxScore - podCount) << 48

		// by zone
		podCount = 0
		zone := zoneMap[node.Labels[h.C.ZoneSelectorLabel]]
		for _, pod := range podList.Items {
			if _, ok := zone.NodeMap[pod.Spec.NodeName]; ok {
				podCount++
			}
		}
		scoreByZone := (maxScore - podCount) << 32

		// by podCount / nodeCount of zone
		nodeCountOfZone := len(zone.NodeMap)
		per := (podCount/nodeCountOfZone)*1000 + 1
		scoreByPer := per << 16

		// sum
		score := int64(scoreByNode + scoreByZone + scoreByPer)
		priorityList[i] = schedulerApi.HostPriority{Host: node.Name, Score: score}
	}
	ctx.JSON(http.StatusOK, priorityList)
}

func (h *Handler) GetZoneMap(items []apiCoreV1.Node) map[string]Zone {
	zoneMap := make(map[string]Zone)
	for _, node := range items {
		zone, ok := zoneMap[node.Labels[h.C.ZoneSelectorLabel]]
		if !ok {
			zone = Zone{NodeMap: map[string]apiCoreV1.Node{}}
		}
		zone.NodeMap[node.Name] = node
		zoneMap[node.Labels[h.C.ZoneSelectorLabel]] = zone
	}
	return zoneMap
}
