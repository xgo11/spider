format:
	@for s in $(shell git status -s |grep -e '.go$$'|grep -v '^D ' |awk '{print $$NF}'); do \
		go fmt $$s;\
	done

v:
	GOPROXY=https://goproxy.io go mod vendor


clean:
	@git clean -f -d -X

fetcher: v
	@go build -mod=vendor -ldflags "-s -w"  -v -o bin/fetcherServer servers/fetcher.go
