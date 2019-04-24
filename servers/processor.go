package main

import (
	"log"
	"net/http"

	"github.com/xgo11/spider/components/processor"
)

func main() {
	pt := processor.NewProcessor(nil, nil, nil, nil)
	handle := pt.HttpServe()
	svr := http.Server{
		Addr:    ":22000",
		Handler: handle,
	}

	if err := svr.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
