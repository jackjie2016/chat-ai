package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ouqiang/goproxy"
)

func main() {
	proxy := goproxy.New()
	server := &http.Server{
		Addr:         "0.0.0.0:7890",
		Handler:      proxy,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}
	fmt.Println("代理地址：0.0.0.0:7890")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}
