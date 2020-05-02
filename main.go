package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/valyala/fasthttp"
)

var (
	log = getLogger("web-proxy-server")
)

func main() {
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			switch string(ctx.Path()) {
			case "/health":
				healthCheck(ctx)
			default:
				HandleRequestAndRedirect(ctx)
			}
		},
		ReadBufferSize: 8192,
	}

	log.Infof("Server will run on: %d\n", config.Port)

	go s.ListenAndServe(fmt.Sprintf("%s%d", ":", config.Port))

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Debug("start shutting down")
	if err := s.Shutdown(); err != nil {
		log.Error(err)
	}
	log.Debug("finish shutting down")

}
