package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	schedulerApi "k8s.io/kube-scheduler/extender/v1"
	"net/http"
)

func (h *Handler) Bind(ctx *gin.Context) {

	args := &schedulerApi.ExtenderBindingArgs{}
	err := ctx.ShouldBindJSON(args)
	if err != nil {
		return
	}
	result := &schedulerApi.ExtenderBindingResult{}

	podName, ns := args.PodName, args.PodNamespace
	binding := &corev1.Binding{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      podName,
			UID:       args.PodUID,
		},
		Target: corev1.ObjectReference{
			Kind: "node",
			Name: args.Node,
		},
	}

	err = h.Client.CoreV1().Pods("").Bind(context.TODO(), binding, metav1.CreateOptions{})
	if err != nil {
		result.Error = err.Error()
		return
	}

	_, err = h.Client.CoreV1().Pods("").Get(context.TODO(), podName, metav1.GetOptions{})

	ctx.JSON(http.StatusOK, result)

}
