package handler

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"violin-home.cn/violin-extender/config"
)

type Handler struct {
	Client *kubernetes.Clientset
	C      *config.Config
}

type Zone struct {
	NodeMap map[string]v1.Node
}
