package main

import (
	"fmt"
	"net/http"

	"github.com/walk1ng/gin-photo-gallery-storage/conf"
	"github.com/walk1ng/gin-photo-gallery-storage/constant"

	"github.com/walk1ng/gin-photo-gallery-storage/routers"
)

func main() {
	// get global router from routers
	router := routers.Router

	// setup a http server
	server := http.Server{
		Addr:           fmt.Sprintf(":%s", conf.ServerCfg.Get(constant.ServerPort)),
		Handler:        router,
		MaxHeaderBytes: 1 << 20,
	}

	// run and listen
	server.ListenAndServe()
}
