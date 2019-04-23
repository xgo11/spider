package main

import (
	"log"
	"net/http"

	"github.com/xgo11/spider/components/fetcher"
)

func main() {

	ft := fetcher.NewFetcher(nil, nil, nil)

	handle := ft.HttpServe()
	svr := http.Server{
		Addr:    ":21000",
		Handler: handle,
	}

	if err := svr.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
