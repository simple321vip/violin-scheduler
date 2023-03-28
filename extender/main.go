package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"violin-home.cn/violin-extender/common"
	"violin-home.cn/violin-extender/config"
	"violin-home.cn/violin-extender/handler"
)

func main() {

	c := config.InitConfig()
	var kubeConfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeConfig = flag.String("kubeConfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeConfig file")
	} else {
		kubeConfig = flag.String("kubeConfig", "", "absolute path to the kubeConfig file")
	}
	flag.Parse()
	cfg, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		panic(err)
	}
	clientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	h := &handler.Handler{
		Client: clientSet,
		C:      c,
	}
	r.Group("/scheduler")
	{
		r.POST("/predicate", h.Predicate)
		r.POST("/preempt", h.Preempt)
		r.POST("/prioritise", h.Prioritise)
		r.POST("/bind", h.Bind)
	}
	common.Run(r, c.SC.Name, c.SC.Addr)
}
