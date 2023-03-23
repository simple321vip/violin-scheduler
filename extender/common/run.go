package common

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Run(eg *gin.Engine, srvName string, addr string) {

	srv := &http.Server{Addr: addr, Handler: eg}

	go func() {
		log.Printf("%s running in %s \n", srvName, addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("%s", err)
		}
		LG.Info("VIOLIN-NOTICE SERVER STARTED SUCCESSFUL")
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	LG.Info("VIOLIN-NOTICE SERVER STOP SUCCESSFUL")

}
