package fetcher

import (
	"sync"
)
import (
	"github.com/xgo11/spider/core"
)

func NewFetcher(scheduler2FetcherQ, fetcher2ProcessorQ, statusQ core.Queue, hooks ...core.FetcherHook) core.IFetcher {
	hf := &httpFetcher{
		schedule2FetcherQ:  scheduler2FetcherQ,
		fetcher2ProcessorQ: fetcher2ProcessorQ,
		statusQ:            statusQ,
		wg:                 &sync.WaitGroup{},
		pause:              false,
		isRunning:          false,
	}
	hf.registerHooks(hooks...)
	return hf
}
